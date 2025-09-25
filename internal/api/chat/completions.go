// internal/api/chat/completions.go
package chat

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"llmapisrv/internal/service"
	"llmapisrv/pkg/logger"
	"llmapisrv/pkg/queue"
	"llmapisrv/pkg/util"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	newAPIService *service.NewAPIService
	logService    *service.LogService
	queue         *queue.RedisQueue
}

func NewChatHandler(
	newAPIService *service.NewAPIService,
	logService *service.LogService,
	queue *queue.RedisQueue,
) *ChatHandler {
	return &ChatHandler{
		newAPIService: newAPIService,
		logService:    logService,
		queue:         queue,
	}
}

// ChatCompletions 处理聊天完成请求
func (h *ChatHandler) ChatCompletions(c *gin.Context) {
	// 获取API Key
	apiKey := c.GetHeader("Authorization")
	if apiKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "API key is required"})
		return
	}

	// 移除Bearer前缀
	apiKey = util.ExtractToken(apiKey)

	// 读取请求体
	var requestBody map[string]interface{}
	if v, exists := c.Get("req"); exists {
		requestBody = v.(map[string]interface{})
	}
	logger.Infof("ChatCompletions reqBody: %v", util.ToJSONString(requestBody))

	// 检查是否为流式响应
	isStream, ok := requestBody["stream"].(bool)
	if !ok {
		isStream = false
	}

	// 转发请求
	startTime := time.Now()
	resp, err := h.newAPIService.ChatCompletion(apiKey, requestBody)
	if err != nil {
		logger.Infof("ChatCompletion got err: %v", err.Error())
		util.ServerError(c, err)
		return
	}
	defer resp.Body.Close()

	// 根据是否为流式响应选择不同的处理方式
	if isStream {
		// 处理流式响应
		h.handleStreamResponse(c, resp, apiKey, requestBody)
	} else {
		// 处理非流式响应
		// 读取响应
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var responseData map[string]interface{}
		if err := json.Unmarshal(body, &responseData); err == nil {
			// 提取使用情况
			if usage, ok := responseData["usage"].(map[string]interface{}); ok {
				// 发送到队列，异步记录日志
				logData := map[string]interface{}{
					"api_key":  strings.Replace(apiKey, "sk-", "", -1),
					"model":    requestBody["model"],
					"usage":    usage,
					"duration": time.Since(startTime).Milliseconds(),
				}
				h.queue.Push("log:chat", logData)
			}
		}

		// 返回原始响应
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
	}
}

// 流式响应处理
func (h *ChatHandler) handleStreamResponse(c *gin.Context, resp *http.Response, apiKey string, requestBody map[string]interface{}) {
	// 设置响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// 创建缓冲区
	var lastChunk []byte
	var chunkList [][]byte

	// 读取响应流
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		chunk := scanner.Bytes()

		// 发送给客户端
		c.Writer.Write(chunk)
		c.Writer.Write([]byte("\n"))
		c.Writer.Flush()

		// 保存最后一块数据
		// lastChunk = chunk

		chunkList = append(chunkList, chunk)
	}

	var chunkData map[string]interface{}
	chunkLength := len(chunkList)
	for i := chunkLength - 1; i >= int(chunkLength/2); i-- {
		if len(chunkList[i]) < 6 {
			continue
		}
		fmt.Println(string(chunkList[i][6:])) // 跳过开头的 data:
		err := json.Unmarshal(chunkList[i][6:], &chunkData)
		if err != nil {
			continue
		}
		if _, ok := chunkData["usage"]; ok {
			lastChunk = chunkList[i][6:]
			break
		}

	}

	// 处理最后一块数据，提取使用情况
	if len(lastChunk) > 0 {
		if usage, ok := chunkData["usage"].(map[string]interface{}); ok {
			// 发送到队列，异步记录日志
			logData := map[string]interface{}{
				"api_key": strings.Replace(apiKey, "sk-", "", -1),
				"model":   requestBody["model"],
				"usage":   usage,
			}
			h.queue.Push("log:chat", logData)
		}
	}
}

// internal/api/log_handler.go
package api

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"llmapisrv/internal/service"
	"llmapisrv/pkg/util"
)

type LogHandler struct {
	logService *service.LogService
}

func NewLogHandler(logService *service.LogService) *LogHandler {
	return &LogHandler{
		logService: logService,
	}
}

// GetLogs 获取日志
func (h *LogHandler) GetLogs(c *gin.Context) {
	userID := c.GetUint("user_id")

	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 获取日志
	logs, total, err := h.logService.GetLogsByUserID(userID, page, pageSize)
	if err != nil {
		util.Fail(c, util.FailCode, err.Error())
		return
	}

	/**
	logs数组结构
	{
		"id": 52,
		"user_id": 1,
		"remote_log_id": 4015,
		"created_at": 1746433683,
		"type": 2,
		"content": "模型倍率 15.00，补全倍率 4.00，分组倍率 1.00",
		"model_name": "xa-claude37",
		"quota": 10785,
		"prompt_tokens": 19,
		"completion_tokens": 175,
		"use_time": 11,
		"is_stream": false,
		"channel": 0,
		"channel_name": "",
		"token_id": 2,
		"group": "default",
		"other": "{\"admin_info\":{\"use_channel\":[\"3\"]},\"cache_ratio\":1,\"cache_tokens\":0,\"completion_ratio\":4,\"frt\":-1000,\"group_ratio\":1,\"is_model_mapped\":true,\"model_price\":-1,\"model_ratio\":15,\"upstream_model_name\":\"claude-3-7-sonnet-20250219\"}",
		"upstream_model_name": "claude-3-7-sonnet-20250219"
	}
	*/
	// 构建响应
	response := gin.H{
		"data": logs,
		"meta": gin.H{
			"current_page": page,
			"page_size":    pageSize,
			"total":        total,
			"total_pages":  (total + int64(pageSize) - 1) / int64(pageSize),
		},
	}

	util.Success(c, response)
}

// CleanupOldLogs 清理旧日志（仅清理本地数据库）
func (h *LogHandler) CleanupOldLogs(c *gin.Context) {
	// 清理旧日志
	err := h.logService.CleanupOldLogs()
	if err != nil {
		util.Fail(c, util.FailCode, err.Error())
		return
	}

	util.Success(c, "CleanupOldLogs success")
}

package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"llmapisrv/pkg/logger"
)

// Logger 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)                 // Use io.ReadAll
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody)) // Use io.NopCloser
		}

		// 使用自定义ResponseWriter记录响应
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 计算延迟
		latency := time.Since(start)

		// 获取客户端信息
		clientInfo, exists := c.Get("client_info")
		if exists {
			info := clientInfo.(*ClientInfo)
			// 使用客户端信息
			logger.InfoWithCtx(ctx, "Request",
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", info.RealIP),
				zap.String("user-agent", info.UserAgent),
				zap.Int("status", c.Writer.Status()),
				zap.Duration("latency", latency),
				zap.String("browser", info.Browser),
				zap.String("forwarded-for", info.ForwardedFor),
				zap.String("os", info.OS),
				zap.String("device-type", info.DeviceType),
				zap.String("referer", info.Referer),
				zap.String("Authorization", info.Authorization),
				zap.String("request", string(requestBody[:min(len(requestBody), 2048)])),      // 只输出前2048个字符
				zap.String("response", blw.body.String()[:min(len(blw.body.String()), 2048)]), // 只输出前2048个字符
			)
		} else {
			// 记录日志
			logger.InfoWithCtx(ctx, "Request",
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.Int("status", c.Writer.Status()),
				zap.Duration("latency", latency),
				zap.String("request", string(requestBody[:min(len(requestBody), 2048)])),      // 只输出前2048个字符
				zap.String("response", blw.body.String()[:min(len(blw.body.String()), 2048)]), // 只输出前2048个字符
			)
		}
	}
}

// bodyLogWriter 自定义ResponseWriter，用于记录响应体
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 重写Write方法，同时写入原ResponseWriter和缓冲区
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

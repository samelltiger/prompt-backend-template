// internal/middleware/metrics.go
package middleware

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// 请求计数器
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// 请求持续时间
	httpRequestDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{0.1, 0.3, 0.5, 0.7, 1, 5, 10, 20, 30, 60},
		},
		[]string{"method", "endpoint", "status"},
	)

	// 模型调用计数器
	modelCallsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "model_calls_total",
			Help: "Total number of model API calls",
		},
		[]string{"model", "status"},
	)

	// Token使用计数器
	tokenUsageTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "token_usage_total",
			Help: "Total number of tokens used",
		},
		[]string{"model", "type"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDurationSeconds)
	prometheus.MustRegister(modelCallsTotal)
	prometheus.MustRegister(tokenUsageTotal)
}

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 获取请求路径和方法
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method

		// 获取模型名称（如果存在）
		var model string
		if c.Request.Method == "POST" && (strings.Contains(path, "/completions") || strings.Contains(path, "/chat/completions")) {
			var requestBody map[string]interface{}
			requestBodyStr, _ := io.ReadAll(c.Request.Body)                // Use io.ReadAll
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBodyStr)) // Use io.NopCloser
			if c.ShouldBindJSON(&requestBody) == nil {
				if modelName, ok := requestBody["model"].(string); ok {
					model = modelName
				}
			}
			c.Set("req", requestBody)
		}

		// 处理请求
		c.Next()

		// 记录指标
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		// HTTP请求指标
		httpRequestsTotal.WithLabelValues(method, path, status).Inc()
		httpRequestDurationSeconds.WithLabelValues(method, path, status).Observe(duration)

		// 如果是模型调用，记录模型调用指标
		if model != "" {
			modelCallsTotal.WithLabelValues(model, status).Inc()

			// 如果存在token使用信息，记录token使用指标
			if tokenUsage, exists := c.Get("token_usage"); exists {
				if usage, ok := tokenUsage.(map[string]interface{}); ok {
					if promptTokens, ok := usage["prompt_tokens"].(float64); ok {
						tokenUsageTotal.WithLabelValues(model, "prompt").Add(promptTokens)
					}
					if completionTokens, ok := usage["completion_tokens"].(float64); ok {
						tokenUsageTotal.WithLabelValues(model, "completion").Add(completionTokens)
					}
				}
			}
		}
	}
}

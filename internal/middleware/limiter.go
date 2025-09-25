// internal/middleware/limiter.go
package middleware

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	"llmapisrv/config"
	"llmapisrv/pkg/cache"
	"llmapisrv/pkg/util"
)

// RateLimiterMiddleware 限流中间件
func RateLimiterMiddleware(cache *cache.RedisCache, limit int, window int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户标识
		userID, exists := c.Get("user_id")
		if !exists {
			util.Fail(c, util.UnauthorizedCode, "Unauthorized")
			c.Abort()
			return
		}

		// 构建缓存键
		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = c.Request.URL.Path
		}

		key := "rate_limit:" + endpoint + ":" + strconv.Itoa(int(userID.(uint)))

		// 获取当前计数
		countStr, err := cache.Get(key)
		if err != nil && err != redis.Nil {
			// 发生错误，但允许请求通过
			c.Next()
			return
		}

		var count int
		if err != redis.Nil {
			count, _ = strconv.Atoi(countStr)
		}

		// 检查是否超过限制
		if count >= limit {
			limitErr := fmt.Sprintf("error: %v, limit: %v, window: %v, unit: seconds", "Rate limit exceeded", limit, window)
			util.Fail(c, util.LimitErrorCode, limitErr)
			c.Abort()
			return
		}

		// 增加计数
		if err == redis.Nil {
			// 首次请求，设置初始值和过期时间
			cache.Set(key, "1", window)
		} else {
			// 增加计数，但保持原有过期时间
			ttl, _ := cache.TTL(key)
			cache.Set(key, strconv.Itoa(count+1), int(ttl.Seconds()))
		}

		c.Next()
	}
}

// BillingRateLimiter 账单查询限流
func BillingRateLimiter(cache *cache.RedisCache, cfg *config.Config) gin.HandlerFunc {
	return RateLimiterMiddleware(cache, cfg.RateLimit.BillingQueryLimit, 60) // 每分钟限制
}

// LogRateLimiter 日志查询限流
func LogRateLimiter(cache *cache.RedisCache, cfg *config.Config) gin.HandlerFunc {
	return RateLimiterMiddleware(cache, cfg.RateLimit.LogQueryLimit, 60) // 每分钟限制
}

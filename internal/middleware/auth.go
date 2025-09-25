// internal/middleware/auth.go
package middleware

import (
	"strconv"
	"time"

	"llmapisrv/internal/service"
	"llmapisrv/pkg/cache"
	"llmapisrv/pkg/logger"
	"llmapisrv/pkg/util"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(userService *service.UserService, cache *cache.RedisCache) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过对健康检查的认证
		if c.Request.URL.Path == "/api/about" {
			c.Next()
			return
		}

		clientInfoInterface, exists := c.Get("client_info")
		if !exists {
			util.Fail(c, util.UnauthorizedCode, "client_info get error, in AuthMiddleware")
			c.Abort()
			return
		}

		clientInfo := clientInfoInterface.(*ClientInfo)

		// 获取API Key
		authHeader := clientInfo.Authorization
		if authHeader == "" {
			util.Fail(c, util.UnauthorizedCode, "API key is required")
			c.Abort()
			return
		}

		// 提取token
		token := clientInfo.AuthNoSk

		// 检查缓存
		cacheKey := "auth:" + token
		if v, err := cache.Get(cacheKey); err == nil {
			// 缓存命中，继续请求
			userId, err := strconv.ParseUint(v, 10, 64)
			logger.Infof("in AuthMiddleware: value: %v, user id: %v", v, userId)
			if err == nil {
				c.Set("user_id", uint(userId))
				c.Set("api_key", token)
				c.Next()
				return
			}
		}

		// 缓存未命中，查询数据库
		user, err := userService.GetUserByAPIKey(token)
		if err != nil {
			util.Fail(c, util.UnauthorizedCode, "Invalid API key")
			c.Abort()
			return
		}

		// 检查用户状态
		if user.Status != 1 {
			util.Fail(c, util.UnauthorizedCode, "API key is disabled")
			c.Abort()
			return
		}

		// 检查余额
		// if user.RemainQuota <= 0 {
		// 	util.Fail(c, util.UnauthorizedCode, "Insufficient quota")
		// 	c.Abort()
		// 	return
		// }

		// 检查过期时间
		if user.ExpiredTime > 0 && user.ExpiredTime < time.Now().Unix() {
			util.Fail(c, util.UnauthorizedCode, "API key has expired")
			c.Abort()
			return
		}

		// 将用户信息添加到上下文
		c.Set("user_id", user.ID)
		c.Set("api_key", token)

		// 更新缓存，设置两小时过期
		cache.Set(cacheKey, strconv.FormatInt(int64(user.ID), 10), 2*3600)

		c.Next()
	}
}

// internal/middleware/admin_auth.go
package middleware

import (
	"github.com/gin-gonic/gin"

	"llmapisrv/config"
	"llmapisrv/pkg/util"
)

// AdminAuthMiddleware 管理员认证中间件
func AdminAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取API Key
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			util.Fail(c, util.UnauthorizedCode, "API key is required")
			c.Abort()
			return
		}

		// 提取token
		token := util.ExtractToken(authHeader)

		// 检查是否是管理员key
		if token != cfg.NewAPI.AdminKey {
			util.Fail(c, util.UnauthorizedCode, "Requires administrator privileges")
			c.Abort()
			return
		}

		c.Next()
	}
}

// pkg/util/token.go
package util

import (
	"strings"
)

// ExtractToken 从Authorization头提取token
func ExtractToken(authHeader string) string {
	// 移除Bearer前缀
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return authHeader
}

func RemoveStartSk(apiKey string) string {
	return strings.Replace(apiKey, "sk-", "", 1)
}

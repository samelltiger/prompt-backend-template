package middleware

import (
	"llmapisrv/pkg/util"
	"strings"

	"github.com/gin-gonic/gin"
)

// ClientInfo 客户端信息结构体
type ClientInfo struct {
	IP            string            // 客户端IP
	UserAgent     string            // 用户代理
	Referer       string            // 来源URL
	RequestID     string            // 请求ID（如果有）
	DeviceType    string            // 设备类型
	OS            string            // 操作系统
	Browser       string            // 浏览器
	IsMobile      bool              // 是否为移动设备
	Headers       map[string]string // 所有请求头
	ForwardedFor  string            // X-Forwarded-For
	RealIP        string            // X-Real-IP
	ContentType   string            // 内容类型
	Accept        string            // Accept头
	AcceptLang    string            // Accept-Language
	Authorization string            // 授权头
	AuthNoSk      string            // 去掉sk-的key值，数据库里查询用这个数据
}

// ClientInfoMiddleware 客户端信息中间件
func ClientInfoMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 收集客户端信息
		clientInfo := &ClientInfo{
			IP:            c.ClientIP(),
			UserAgent:     c.Request.UserAgent(),
			Referer:       c.Request.Referer(),
			ForwardedFor:  c.GetHeader("X-Forwarded-For"),
			RealIP:        c.GetHeader("X-Real-IP"),
			ContentType:   c.GetHeader("Content-Type"),
			Accept:        c.GetHeader("Accept"),
			AcceptLang:    c.GetHeader("Accept-Language"),
			Authorization: c.GetHeader("Authorization"),
			Headers:       make(map[string]string),
		}

		token := util.ExtractToken(clientInfo.Authorization)
		clientInfo.AuthNoSk = strings.Replace(token, "sk-", "", -1)

		// 收集所有请求头
		for k, v := range c.Request.Header {
			if len(v) > 0 {
				clientInfo.Headers[k] = v[0]
			}
		}

		// 解析设备信息
		clientInfo.DeviceType, clientInfo.OS, clientInfo.Browser = parseUserAgent(clientInfo.UserAgent)
		clientInfo.IsMobile = strings.Contains(strings.ToLower(clientInfo.UserAgent), "mobile")

		// 将客户端信息存储到上下文中
		c.Set("client_info", clientInfo)

		c.Next()
	}
}

// parseUserAgent 简单的UserAgent解析
func parseUserAgent(ua string) (deviceType, os, browser string) {
	ua = strings.ToLower(ua)

	// 设备类型检测
	switch {
	case strings.Contains(ua, "mobile"):
		deviceType = "mobile"
	case strings.Contains(ua, "tablet"):
		deviceType = "tablet"
	default:
		deviceType = "desktop"
	}

	// 操作系统检测
	switch {
	case strings.Contains(ua, "windows"):
		os = "windows"
	case strings.Contains(ua, "mac os"):
		os = "macos"
	case strings.Contains(ua, "linux"):
		os = "linux"
	case strings.Contains(ua, "android"):
		os = "android"
	case strings.Contains(ua, "iphone"):
		os = "ios"
	case strings.Contains(ua, "ipad"):
		os = "ios"
	default:
		os = "unknown"
	}

	// 浏览器检测
	switch {
	case strings.Contains(ua, "chrome"):
		browser = "chrome"
	case strings.Contains(ua, "firefox"):
		browser = "firefox"
	case strings.Contains(ua, "safari"):
		browser = "safari"
	case strings.Contains(ua, "edge"):
		browser = "edge"
	case strings.Contains(ua, "opera"):
		browser = "opera"
	case strings.Contains(ua, "msie") || strings.Contains(ua, "trident"):
		browser = "ie"
	default:
		browser = "unknown"
	}

	return deviceType, os, browser
}

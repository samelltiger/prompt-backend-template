// internal/api/proxy.go
package api

import (
	"llmapisrv/pkg/oss"
	"llmapisrv/pkg/proxy"
	"llmapisrv/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProxyHandler struct {
	proxyService *proxy.ProxyService
}

func NewProxyHandler(ossClient *oss.OSSClient) *ProxyHandler {
	return &ProxyHandler{
		proxyService: proxy.NewProxyService(ossClient),
	}
}

// ImageProxy 图片代理接口
func (h *ProxyHandler) ImageProxy(c *gin.Context) {
	// 获取参数
	path := c.Query("path")
	thumbStr := c.DefaultQuery("thumb", "false")

	if path == "" {
		util.Fail(c, util.FailCode, "path parameter is required")
		return
	}

	isThumb := thumbStr == "true"

	// 下载并缓存图片
	imageData, err := h.proxyService.DownloadAndCacheImage(c.Request.Context(), path, isThumb)
	if err != nil {
		util.Fail(c, util.FailCode, "Failed to get image: "+err.Error())
		return
	}

	// 设置响应头
	contentType := h.getContentType(path)
	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=86400") // 缓存1天
	
	// 返回图片数据
	c.Data(http.StatusOK, contentType, imageData)
}

// getContentType 根据文件扩展名获取Content-Type
func (h *ProxyHandler) getContentType(filename string) string {
	switch {
	case len(filename) >= 4 && filename[len(filename)-4:] == ".jpg" || filename[len(filename)-5:] == ".jpeg":
		return "image/jpeg"
	case len(filename) >= 4 && filename[len(filename)-4:] == ".png":
		return "image/png"
	case len(filename) >= 4 && filename[len(filename)-4:] == ".gif":
		return "image/gif"
	case len(filename) >= 5 && filename[len(filename)-5:] == ".webp":
		return "image/webp"
	case len(filename) >= 4 && filename[len(filename)-4:] == ".svg":
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}
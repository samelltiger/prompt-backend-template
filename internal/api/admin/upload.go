// internal/api/admin/upload.go
package admin

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"llmapisrv/config"
	"llmapisrv/pkg/oss"
	"llmapisrv/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadImageRequest struct {
	ImageData string `json:"image_data" binding:"required"` // base64编码的图片数据
	ImageType string `json:"image_type"`                    // 图片类型，可选，如 "png", "jpg", "jpeg"
}

type UploadImageResponse struct {
	OSSShortLink string `json:"oss_short_link"` // OSS短链
	PublicURL    string `json:"public_url"`     // 2小时过期的公开访问URL
	FileName     string `json:"file_name"`      // 上传的文件名
}

type UploadHandler struct {
	ossClient *oss.OSSClient
}

func NewUploadHandler(ossClient *oss.OSSClient) *UploadHandler {
	return &UploadHandler{
		ossClient: ossClient,
	}
}

// UploadImage 上传base64编码的图片到OSS
func (h *UploadHandler) UploadImage(c *gin.Context) {
	var req UploadImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.ParamError(c, err.Error())
		return
	}

	// 解析base64数据
	// 检查是否包含数据URL前缀（如 data:image/png;base64,）
	imageData := req.ImageData
	if strings.Contains(imageData, ",") {
		parts := strings.SplitN(imageData, ",", 2)
		if len(parts) == 2 {
			imageData = parts[1]
		}
	}

	// 解码base64数据
	decodedData, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		util.ParamError(c, "Invalid base64 image data")
		return
	}

	// 验证图片数据大小（限制10MB）
	if len(decodedData) > 10*1024*1024 {
		util.ParamError(c, "Image size exceeds 10MB limit")
		return
	}

	// 检测图片类型
	imageType := h.detectImageType(decodedData)
	if imageType == "" {
		if req.ImageType != "" {
			imageType = req.ImageType
		} else {
			imageType = "png" // 默认为png
		}
	}

	// 生成唯一文件名
	fileName := h.generateFileName(imageType)

	// 创建读取器
	reader := bytes.NewReader(decodedData)

	// 上传到OSS（设置为私有，支持临时URL访问）
	ctx := context.Background()
	err = h.ossClient.UploadFile(ctx, fileName, reader, oss.WithPrivate(true))
	if err != nil {
		util.Fail(c, util.FailCode, "Failed to upload image to OSS")
		return
	}

	// 生成2小时过期的访问URL
	publicURL, err := h.ossClient.GetFileURL(ctx, fileName, 2*time.Hour)
	if err != nil {
		util.Fail(c, util.FailCode, "Failed to generate access URL")
		return
	}

	// 生成OSS短链（基于配置）
	ossShortLink := h.generateOSSShortLink(fileName)

	util.Success(c, UploadImageResponse{
		OSSShortLink: ossShortLink,
		PublicURL:    publicURL,
		FileName:     fileName,
	})
}

// detectImageType 检测图片类型
func (h *UploadHandler) detectImageType(data []byte) string {
	if len(data) < 8 {
		return ""
	}

	// PNG: 89 50 4E 47 0D 0A 1A 0A
	if len(data) >= 8 && data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return "png"
	}

	// JPEG: FF D8 FF
	if len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return "jpg"
	}

	// GIF: 47 49 46 38
	if len(data) >= 4 && data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x38 {
		return "gif"
	}

	// WebP: RIFF...WEBP
	if len(data) >= 12 && string(data[0:4]) == "RIFF" && string(data[8:12]) == "WEBP" {
		return "webp"
	}

	return ""
}

// generateFileName 生成唯一文件名
func (h *UploadHandler) generateFileName(imageType string) string {
	// 生成UUID作为文件名
	id := uuid.New().String()

	// 添加时间戳确保唯一性
	timestamp := time.Now().Format("20060102_150405")

	// 构建文件名：images/年月日/timestamp_uuid.扩展名
	dateDir := time.Now().Format("20060102")
	fileName := fmt.Sprintf("images/%s/%s_%s.%s", dateDir, timestamp, id, imageType)

	return fileName
}

// generateOSSShortLink 生成OSS短链
func (h *UploadHandler) generateOSSShortLink(fileName string) string {
	// 基于配置生成OSS短链
	cfg := &config.AppConfig.OSS.Aliyun

	scheme := "https"
	if !cfg.SSL {
		scheme = "http"
	}

	if cfg.IsCNAME {
		return fmt.Sprintf("%s://%s/%s", scheme, cfg.Endpoint, fileName)
	}

	return fmt.Sprintf("%s://%s.%s/%s", scheme, cfg.Bucket, cfg.Endpoint, fileName)
}

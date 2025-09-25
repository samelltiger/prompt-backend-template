// pkg/proxy/proxy.go
package proxy

import (
	"context"
	"fmt"
	"io"
	"llmapisrv/config"
	"llmapisrv/pkg/oss"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ProxyService struct {
	ossClient *oss.OSSClient
}

func NewProxyService(ossClient *oss.OSSClient) *ProxyService {
	return &ProxyService{
		ossClient: ossClient,
	}
}

// SelectProxyServer 根据权重选择代理服务器
func (p *ProxyService) SelectProxyServer() string {
	servers := config.AppConfig.OSS.OSSProxySrvs
	if len(servers) == 0 {
		return ""
	}

	// 计算总权重
	totalWeight := 0
	for _, server := range servers {
		totalWeight += server.Weight
	}

	if totalWeight == 0 {
		return servers[0].IP // 如果所有权重都是0，返回第一个
	}

	// 生成随机数
	rand.Seed(time.Now().UnixNano())
	randomWeight := rand.Intn(totalWeight)

	// 根据权重选择服务器
	currentWeight := 0
	for _, server := range servers {
		currentWeight += server.Weight
		if randomWeight < currentWeight {
			return server.IP
		}
	}

	return servers[0].IP // 默认返回第一个
}

// GenerateProxyURL 生成代理URL
func (p *ProxyService) GenerateProxyURL(ossPath string, isThumb bool) string {
	proxyServer := p.SelectProxyServer()
	if proxyServer == "" {
		return ""
	}

	// 构建代理URL
	thumbParam := "false"
	if isThumb {
		thumbParam = "true"
	}

	// 确保服务器地址包含协议
	if !strings.HasPrefix(proxyServer, "http://") && !strings.HasPrefix(proxyServer, "https://") {
		proxyServer = "http://" + proxyServer
	}

	return fmt.Sprintf("%s/api/image-proxy?path=%s&thumb=%s", proxyServer, ossPath, thumbParam)
}

// DownloadAndCacheImage 下载并缓存图片
func (p *ProxyService) DownloadAndCacheImage(ctx context.Context, ossPath string, isThumb bool) ([]byte, error) {
	// 构建本地缓存路径
	cacheDir := config.AppConfig.OSS.OssCacheDir
	if cacheDir == "" {
		cacheDir = "/images"
	}

	var cachePath string
	if isThumb {
		// 为缩略图添加特殊标识
		ext := filepath.Ext(ossPath)
		nameWithoutExt := strings.TrimSuffix(ossPath, ext)
		cachePath = filepath.Join(cacheDir, nameWithoutExt+"_thumb"+ext)
	} else {
		cachePath = filepath.Join(cacheDir, ossPath)
	}

	// 检查缓存是否存在
	if data, err := os.ReadFile(cachePath); err == nil {
		return data, nil
	}

	// 缓存不存在，从OSS下载
	var signedURL string
	var err error

	if isThumb {
		signedURL, err = p.ossClient.GetThumbnailURL(ctx, ossPath, 2*time.Hour)
	} else {
		signedURL, err = p.ossClient.GetFileURL(ctx, ossPath, 2*time.Hour)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get OSS signed URL: %w", err)
	}

	// 下载图片
	resp, err := http.Get(signedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download image, status: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	// 缓存到本地
	if err := p.cacheImage(cachePath, data); err != nil {
		// 缓存失败不影响返回结果，只记录错误
		fmt.Printf("Failed to cache image %s: %v\n", cachePath, err)
	}

	return data, nil
}

// cacheImage 将图片数据缓存到本地
func (p *ProxyService) cacheImage(cachePath string, data []byte) error {
	// 创建目录
	dir := filepath.Dir(cachePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// 写入文件
	return os.WriteFile(cachePath, data, 0644)
}
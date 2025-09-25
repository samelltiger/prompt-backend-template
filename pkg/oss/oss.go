// pkg/oss/oss.go
package oss

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"llmapisrv/config"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// OSSClient OSS客户端
type OSSClient struct {
	config         *config.AliyunOSS
	uploadClient   *oss.Client
	downloadClient *oss.Client
	bucket         *oss.Bucket
}

// NewOSSClient 创建OSS客户端
func NewOSSClient(cfg *config.AliyunOSS) (*OSSClient, error) {
	// 创建上传客户端
	uploadClient, err := oss.New(cfg.Endpoint, cfg.AccessID, cfg.AccessKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create upload client: %w", err)
	}

	// 创建下载客户端（使用相同的配置）
	downloadClient, err := oss.New(cfg.Endpoint, cfg.AccessID, cfg.AccessKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create download client: %w", err)
	}

	// 获取存储空间
	bucket, err := uploadClient.Bucket(cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to get bucket: %w", err)
	}

	return &OSSClient{
		config:         cfg,
		uploadClient:   uploadClient,
		downloadClient: downloadClient,
		bucket:         bucket,
	}, nil
}

// UploadFile 上传文件
func (c *OSSClient) UploadFile(ctx context.Context, fileName string, reader io.Reader, options ...UploadOption) error {
	// 默认配置
	opts := &uploadOptions{
		private:       true,
		retryTimes:    1,
		retryInterval: 10 * time.Second,
	}

	// 应用选项
	for _, opt := range options {
		opt(opts)
	}

	// 设置请求头
	headers := []oss.Option{
		oss.ObjectACL(oss.ACLPrivate),
	}

	if !opts.private {
		headers[0] = oss.ObjectACL(oss.ACLPublicRead)
	}

	if opts.downloadFilename != "" {
		// 设置Content-Disposition头，支持中文文件名
		encodedFilename := url.QueryEscape(opts.downloadFilename)
		headers = append(headers, oss.ContentDisposition(
			fmt.Sprintf("attachment;filename=%s;filename*=%s", encodedFilename, encodedFilename)))
	}

	// 重试上传
	var lastErr error
	for i := 0; i < opts.retryTimes; i++ {
		if i > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(opts.retryInterval):
			}
		}

		err := c.bucket.PutObject(fileName, reader, headers...)
		if err == nil {
			return nil
		}

		lastErr = err
		if i < opts.retryTimes-1 {
			// 检查错误是否可重试
			if !isRetryableError(err) {
				return err
			}
		}
	}

	return fmt.Errorf("upload failed after %d retries: %w", opts.retryTimes, lastErr)
}

// GetFileURL 获取文件临时下载地址
func (c *OSSClient) GetFileURL(ctx context.Context, fileName string, expires time.Duration) (string, error) {
	if fileName == "" {
		return "", fmt.Errorf("file name cannot be empty")
	}

	// 检查文件权限
	acl, err := c.bucket.GetObjectACL(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to get object ACL: %w", err)
	}

	// 如果是公开读，直接返回公开URL
	if acl.ACL == string(oss.ACLPublicRead) {
		return c.getPublicDownloadURL(fileName), nil
	}

	// 私有文件，生成签名URL
	signedURL, err := c.bucket.SignURL(fileName, oss.HTTPGet, int64(expires.Seconds()))
	if err != nil {
		return "", fmt.Errorf("failed to sign URL: %w", err)
	}

	return signedURL, nil
}

// GetThumbnailURL 获取缩略图URL，包含OSS图片处理参数
func (c *OSSClient) GetThumbnailURL(ctx context.Context, fileName string, expires time.Duration) (string, error) {
	if fileName == "" {
		return "", fmt.Errorf("file name cannot be empty")
	}

	// 检查文件权限
	acl, err := c.bucket.GetObjectACL(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to get object ACL: %w", err)
	}

	// 如果是公开读，生成带处理参数的公开URL
	if acl.ACL == string(oss.ACLPublicRead) {
		baseURL := c.getPublicDownloadURL(fileName)
		return baseURL + "?x-oss-process=image%2Fresize%2Cm_fill%2Cw_330%2Ch_240", nil
	}

	// 私有文件，生成包含处理参数的签名URL
	// 使用阿里云OSS SDK的正确方式传递查询参数
	options := []oss.Option{
		oss.Process("image/resize,m_fill,w_330,h_240"),
	}

	signedURL, err := c.bucket.SignURL(fileName, oss.HTTPGet, int64(expires.Seconds()), options...)
	if err != nil {
		return "", fmt.Errorf("failed to sign URL: %w", err)
	}

	return signedURL, nil
}

// DownloadFile 下载文件
func (c *OSSClient) DownloadFile(ctx context.Context, fileName string) (io.ReadCloser, error) {
	if fileName == "" {
		return nil, fmt.Errorf("file name cannot be empty")
	}

	body, err := c.bucket.GetObject(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	return body, nil
}

// DeleteFile 删除文件
func (c *OSSClient) DeleteFile(ctx context.Context, fileName string) error {
	if fileName == "" {
		return fmt.Errorf("file name cannot be empty")
	}

	err := c.bucket.DeleteObject(fileName)
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

// getPublicDownloadURL 获取公开下载URL
func (c *OSSClient) getPublicDownloadURL(fileName string) string {
	if fileName == "" {
		return ""
	}

	// 构建公开访问URL
	scheme := "https"
	if !c.config.SSL {
		scheme = "http"
	}

	if c.config.IsCNAME {
		return fmt.Sprintf("%s://%s/%s", scheme, c.config.Endpoint, fileName)
	}

	return fmt.Sprintf("%s://%s.%s/%s", scheme, c.config.Bucket, c.config.Endpoint, fileName)
}

// isRetryableError 判断错误是否可重试
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// 网络错误可重试
	if _, ok := err.(*url.Error); ok {
		return true
	}

	// HTTP状态码判断
	if ossErr, ok := err.(oss.ServiceError); ok {
		switch ossErr.StatusCode {
		case http.StatusRequestTimeout,
			http.StatusTooManyRequests,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout:
			return true
		}
	}

	return false
}

// UploadOption 上传选项
type UploadOption func(*uploadOptions)

// uploadOptions 上传配置
type uploadOptions struct {
	private          bool
	downloadFilename string
	retryTimes       int
	retryInterval    time.Duration
}

// WithPrivate 设置私有访问
func WithPrivate(private bool) UploadOption {
	return func(o *uploadOptions) {
		o.private = private
	}
}

// WithDownloadFilename 设置下载文件名
func WithDownloadFilename(filename string) UploadOption {
	return func(o *uploadOptions) {
		o.downloadFilename = filename
	}
}

// WithRetry 设置重试配置
func WithRetry(times int, interval time.Duration) UploadOption {
	return func(o *uploadOptions) {
		o.retryTimes = times
		o.retryInterval = interval
	}
}

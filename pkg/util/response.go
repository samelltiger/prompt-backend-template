package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	// SuccessCode 成功状态码
	SuccessCode = 200
	// FailCode 失败状态码
	FailCode = 400
	// UnauthorizedCode 未授权状态码
	UnauthorizedCode = 401
	// ServerErrorCode 服务器错误状态码
	ServerErrorCode = 500
	// 请求次数限制错误
	LimitErrorCode = 1001
)

// Response 标准响应结构
type Response struct {
	// 状态码
	// Example: 200
	Code int `json:"code"`
	// 错误信息（成功时为空）
	// Example: "参数错误"
	Message string `json:"message"`
	// 业务数据
	Data interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// ParamError 参数错误
func ParamError(c *gin.Context, message string) {
	if message == "" {
		message = "参数错误"
	}
	c.JSON(http.StatusOK, Response{
		Code:    400,
		Message: message,
	})
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    401,
		Message: "未授权或授权已过期",
	})
}

// Unauthorized 未授权
func UnauthorizedWithMsg(c *gin.Context, msg string) {

	if msg == "" {
		Unauthorized(c)
	} else {
		c.JSON(http.StatusUnauthorized, Response{
			Code:    401,
			Message: msg,
		})
	}
}

// ServerError 服务器错误
func ServerError(c *gin.Context, err error) {
	msg := "服务器内部错误"
	if err != nil {
		msg = err.Error()
	}
	c.JSON(http.StatusInternalServerError, Response{
		Code:    500,
		Message: msg,
	})
}

// PageResult 分页结果
type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// PageSuccess 分页成功响应
func PageSuccess(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data: PageResult{
			List:     list,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

// 资源没找到
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "资源不存在"
	}
	c.JSON(http.StatusOK, Response{
		Code:    400,
		Message: message,
	})
}

// SuccessWithPagination sends a success response with pagination information
func SuccessWithPagination(c *gin.Context, data interface{}, total int64, page, size int) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    data,
		"pagination": gin.H{
			"total":       total,
			"page":        page,
			"size":        size,
			"total_pages": (total + int64(size) - 1) / int64(size),
		},
	})
}

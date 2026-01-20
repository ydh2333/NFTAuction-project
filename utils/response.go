package utils

import (
	"github.com/gin-gonic/gin"
)

// SuccessResponse 成功响应结构体
type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse 错误响应结构体
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SendSuccess 发送成功响应
func SendSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(200, SuccessResponse{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

// SendError 发送错误响应
func SendError(c *gin.Context, code int, message string) {
	c.JSON(code, ErrorResponse{
		Code:    code,
		Message: message,
	})
	c.Abort() // 终止后续处理
}

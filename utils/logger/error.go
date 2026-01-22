package logger

import (
	"fmt"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

// 自定义错误结构体
type AppError struct {
	Msg  string
	Err  error
	Code int
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

// 包装错误
func WrapError(err error, msg string, args ...interface{}) *AppError {
	return &AppError{
		Msg:  sprintf(msg, args...),
		Err:  err,
		Code: 500,
	}
}

// 创建错误
func NewErrorf(msg string, args ...interface{}) *AppError {
	return &AppError{
		Msg:  sprintf(msg, args...),
		Err:  nil,
		Code: 500,
	}
}

// 格式化字符串（简化版fmt.Sprintf）
func sprintf(msg string, args ...interface{}) string {
	if len(args) == 0 {
		return msg
	}
	return fmt.Sprintf(msg, args...)
}

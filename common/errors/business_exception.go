package errors

import (
	"fmt"
)

// BusinessError 定义一个结构体来表示业务请求异常
type BusinessError struct {
	// 异常状态码
	code int
}

// NewBusinessError 创建一个新的 BusinessError 实例
func NewBusinessError(code int) *BusinessError {
	return &BusinessError{
		code: code,
	}
}

// Code 获取异常状态码
func (e *BusinessError) Code() int {
	return e.code
}

// Error 实现 error 接口
func (e *BusinessError) Error() string {
	return fmt.Sprintf("Business request exception with code: %d", e.code)
}

package common

import (
	"fmt"
)

// BusinessRequestException 定义一个结构体来表示业务请求异常
type BusinessRequestException struct {
	// 异常状态码
	code int
}

// NewBusinessRequestException 创建一个新的 BusinessRequestException 实例
func NewBusinessRequestException(code int) *BusinessRequestException {
	return &BusinessRequestException{
		code: code,
	}
}

// Code 获取异常状态码
func (e *BusinessRequestException) Code() int {
	return e.code
}

// Error 实现 error 接口
func (e *BusinessRequestException) Error() string {
	return fmt.Sprintf("Business request exception with code: %d", e.code)
}

package client

import "time"

// RequestCallback 定义了请求回调接口
type RequestCallback interface {
	OnSuccess(callBack any)
	OnError(error)
}

// RequestResponseFuture 定义了请求响应状态
type RequestResponseFuture struct {
	start           int
	Cause           error
	Response        any
	RequestCallback RequestCallback
	waitResponse    chan any
	waitCause       chan error
}

func (f *RequestResponseFuture) isTimeout() bool {
	now := time.Now().Second()
	return now-f.start > 5
}

package client

import "time"

// RequestCallback 定义了请求回调接口
type RequestCallback interface {
	OnSuccess(callBack any)
	OnError(error)
}

// RequestResponseFuture 定义了请求响应未来状态
type RequestResponseFuture struct {
	start           int
	Cause           error
	Response        any
	RequestCallback RequestCallback
}

func (f *RequestResponseFuture) isTimeout() bool {
	now := time.Now().Second()
	return now-f.start > 5
}

// ExecuteRequestCallback 执行请求回调
func (r *RequestResponseFuture) ExecuteRequestCallback() {
	if r.RequestCallback != nil {
		if r.Cause != nil {
			r.RequestCallback.OnError(r.Cause)
		} else {
			r.RequestCallback.OnSuccess(r.Response)
		}
	}
}

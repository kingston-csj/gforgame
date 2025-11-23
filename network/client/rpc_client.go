package client

import (
	"io/github/gforgame/network"
	"sync/atomic"
	"time"
)

var (
	// 客户端消息流水号，用于实现消息回调
	counter int32 = 10
)

// 客户端发送消息，并注册回调函数
func Callback(session *network.Session, request any, callback RequestCallback) {
	atomic.AddInt32(&counter, 1)
	future := &RequestResponseFuture{RequestCallback: callback, start: time.Now().Second()}
	CallBackManager.Register(int(counter), future)
	session.Send(request, int(counter))
}

func Request(session *network.Session, request any) (any, error) {
	atomic.AddInt32(&counter, 1)
	session.Send(request, int(counter))
	future := &RequestResponseFuture{start: time.Now().Second()}
	future.waitCause = make(chan error)
	future.waitResponse = make(chan any)
	CallBackManager.Register(int(counter), future)
	// 调用成功，获得消息返回值; 失败，获得错误（如超时）
	// 这里的代码相对java来得优雅
	select {
	case r := <-future.waitResponse:
		return r, nil
	case e := <-future.waitCause:
		return nil, e
	}

}

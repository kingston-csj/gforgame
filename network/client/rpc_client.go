package client

import (
	"io/github/gforgame/network"
	"sync/atomic"
	"time"
)

var (
	counter int32 = 10
)

// 客户端发送消息，并注册回调函数
func SendCallback(session *network.Session, request any, callback RequestCallback) {
	atomic.AddInt32(&counter, 1)
	future := &RequestResponseFuture{RequestCallback: callback, start: time.Now().Second()}
	CallBackManager.Register(int(counter), future)
	session.Send(request, int(counter))
}

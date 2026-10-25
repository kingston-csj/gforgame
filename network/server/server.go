package server

import "github.com/forfun/gforgame/network/protocol"

// BaseServer 为 TCP 和 WebSocket 服务器提供独立的消息工厂。
type BaseServer struct {
	*protocol.MessageFactory
}

func NewBaseServer(routerFactory *protocol.MessageFactory) BaseServer {
	if routerFactory == nil {
		routerFactory = protocol.NewMessageFactory()
	}
	return BaseServer{MessageFactory: routerFactory}
}

func (s *BaseServer) GetMessageFactory() *protocol.MessageFactory { return s.MessageFactory }

type Server interface {
	GetMessageFactory() *protocol.MessageFactory
	// Addr 监听地址
	Addr() string
	// Start 启动服务器
	Start() error
	// Stop 关闭服务器
	Stop()
	// RunningChan 接收外部关服信号
	RunningChan() <-chan bool
	// NotifyStop 主动通知服务器进入关停流程
	NotifyStop()
}

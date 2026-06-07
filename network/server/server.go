package server

type Server interface {
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

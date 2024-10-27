package network

type Server interface {
	// Addr 监听地址
	Addr() string
	// Start 启动服务器
	Start() error
	// Stop 关闭服务器
	Stop() error
}

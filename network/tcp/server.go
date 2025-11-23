package tcp

import (
	"errors"
	"fmt"
	"io/github/gforgame/logger"
	"io/github/gforgame/network"
	"net"
)

type TcpServer struct {
	Options
	Name    string // 服务器名称
	Running chan bool
}

func NewServer(opts ...Option) *TcpServer {
	opt := Options{}
	for _, option := range opts {
		option(&opt)
	}

	s := &TcpServer{
		Options: opt,
		Running: make(chan bool),
	}

	return s
}

func (s *TcpServer) Start() error {
	if s.ServiceAddr == "" {
		return errors.New("service address cannot be empty")
	}

	modules := s.modules
	for _, c := range modules {
		c.Init()
		err := s.Router.RegisterMessageHandlers(c)
		if err != nil {
			panic(err)
		}
	}

	go func() {
		s.startListen()
	}()

	return nil
}

func (s *TcpServer) Addr() string {
	return s.ServiceAddr
}

// Enable current server accept connection
func (s *TcpServer) startListen() {
	listener, err := net.Listen("tcp", s.ServiceAddr)
	if err != nil {
		logger.Error(err)
	}

	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error(fmt.Errorf("new tcp conn failed %v", err))
			continue
		}
		go onClientConnected(s, conn)
	}
}

// 处理客户端连接，包括socket,websocket
func onClientConnected(node *TcpServer, conn net.Conn) {
	defer func() {
		// 处理客户端网络断开
		s := network.GetSession(conn)
		node.IoDispatch.OnSessionClosed(s)
		network.UnregisterSession(conn)
		err := conn.Close()
		if err != nil {
			logger.Error(fmt.Errorf("close tcp conn failed %v", err))
		}
	}()

	ioSession := network.NewSession(conn, node.MessageCodec)
	network.RegisterSession(conn, ioSession)

	// session created hook
	node.IoDispatch.OnSessionCreated(ioSession)

	// 异步读写数据
	go ioSession.Read()
	go ioSession.Write()

	// read loop
	//  轮询，保证异步任务和客户端消息的执行是线程安全的
	select {
	case task := <-ioSession.AsynTasks:
		task()
	case ioFrame := <-ioSession.DataReceived:
		node.IoDispatch.OnMessageReceived(ioSession, ioFrame)
	case <-ioSession.Die:
		// 关闭session，执行defer函数
	}
}

func (n *TcpServer) Stop() {
	for _, c := range n.modules {
		c.Shutdown()
	}
}

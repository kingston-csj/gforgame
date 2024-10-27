package tcp

import (
	"errors"
	"fmt"
	"io/github/gforgame/logger"
	"io/github/gforgame/network"
	"net"
)

type tcpServer struct {
	Options
	Name    string // 服务器名称
	Running chan bool
}

func NewServer(opts ...Option) *tcpServer {
	opt := Options{}
	for _, option := range opts {
		option(&opt)
	}

	s := &tcpServer{
		Options: opt,
		Running: make(chan bool),
	}

	return s
}

func (s *tcpServer) Start() error {
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

func (s *tcpServer) Addr() string {
	return s.ServiceAddr
}

// Enable current server accept connection
func (s *tcpServer) startListen() {
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
func onClientConnected(node *tcpServer, conn net.Conn) {
	defer func() {
		// 处理客户端网络断开
		s := network.GetSession(conn)
		node.IoDispatch.OnSessionCreated(s)
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
	for ioFrame := range ioSession.DataReceived {
		node.IoDispatch.OnMessageReceived(ioSession, ioFrame)
	}
}

func (n *tcpServer) Stop() {
	for _, c := range n.modules {
		c.Shutdown()
	}
}

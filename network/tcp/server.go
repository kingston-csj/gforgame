package tcp

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"

	"github.com/forfun/gforgame/network"
)

type TcpServer struct {
	Options
	Name     string // 服务器名称
	Running  chan bool
	listener net.Listener
	stopOnce sync.Once
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
			return err
		}
	}

	listener, err := net.Listen("tcp", s.ServiceAddr)
	if err != nil {
		return err
	}
	s.listener = listener

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
	if s.listener == nil {
		slog.Error("tcp listener is nil")
		return
	}

	defer func() {
		_ = s.listener.Close()
	}()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			// listener 被关闭时会返回错误，此时应退出 accept 循环
			if errors.Is(err, net.ErrClosed) {
				return
			}
			slog.Error(fmt.Sprintf("new tcp conn failed %v", err))
			return
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
		_ = conn.Close()
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
	for {
		select {
		case task := <-ioSession.AsynTasks:
			task()
		case ioFrame := <-ioSession.DataReceived:
			node.IoDispatch.OnMessageReceived(ioSession, ioFrame)
		case <-ioSession.Die:
			// 关闭session，执行defer函数
			return
		}
	}
}

func (n *TcpServer) Stop() {
	n.stopOnce.Do(func() {
		if n.listener != nil {
			_ = n.listener.Close()
		}
		network.CloseAllSessions()
		for _, c := range n.modules {
			c.Shutdown()
		}
	})
}

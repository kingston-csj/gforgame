package tcp

import (
	"errors"
	"fmt"
	"net"
	"runtime"
	"sync"

	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/network"
	serverpkg "github.com/forfun/gforgame/network/server"
)

type TcpServer struct {
	Options
	Name     string // 服务器名称
	Running  chan bool
	listener net.Listener
	stopOnce sync.Once
}

var _ serverpkg.Server = (*TcpServer)(nil)

func NewServer(opts ...Option) *TcpServer {
	opt := Options{BaseServerOptions: serverpkg.BaseServerOptions{DispatchWorkers: 1}}
	for _, option := range opts {
		option(&opt)
	}
	if opt.UseGateway && opt.DispatchWorkers <= 0 {
		opt.DispatchWorkers = int32(runtime.NumCPU())
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

func (s *TcpServer) RunningChan() <-chan bool {
	return s.Running
}

func (s *TcpServer) NotifyStop() {
	s.Running <- true
}

// Enable current server accept connection
func (s *TcpServer) startListen() {
	if s.listener == nil {
		logger.ErrorNoStack("tcp listener is nil")
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
			logger.ErrorNoStack(fmt.Sprintf("new tcp conn failed %v", err))
			return
		}
		go func(conn net.Conn) {
			network.ServeSessionConn(conn, s.MessageCodec, s.IoDispatch, s.DispatchWorkers, s.PayloadMode)
		}(conn)
	}
}

func (n *TcpServer) Stop() {
	n.stopOnce.Do(func() {
		if n.listener != nil {
			_ = n.listener.Close()
		}
		network.CloseAllSessions()
	})
}

package ws

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/network"
	serverpkg "github.com/forfun/gforgame/network/server"

	"github.com/gorilla/websocket"
)

type WsServer struct {
	Options
	Name     string // 服务器名称
	Running  chan bool
	server   *http.Server
	listener net.Listener
	stopOnce sync.Once
}

var _ serverpkg.Server = (*WsServer)(nil)

func NewServer(opts ...Option) *WsServer {
	opt := Options{BaseServerOptions: serverpkg.BaseServerOptions{DispatchWorkers: 1}}
	for _, option := range opts {
		option(&opt)
	}
	if opt.UseGateway && opt.DispatchWorkers <= 0 {
		opt.DispatchWorkers = int32(runtime.NumCPU())
	}

	if opt.wsPath == "" {
		opt.wsPath = "ws"
	}

	s := &WsServer{
		Options: opt,
		Running: make(chan bool),
	}

	return s
}

func (n *WsServer) Start() error {
	if n.ServiceAddr == "" {
		return errors.New("service address cannot be empty")
	}

	if err := n.startListen(); err != nil {
		return err
	}

	return nil
}

func (n *WsServer) Addr() string {
	return n.ServiceAddr
}

func (n *WsServer) RunningChan() <-chan bool {
	return n.Running
}

func (n *WsServer) NotifyStop() {
	n.Running <- true
}

// Enable current server accept connection
func (n *WsServer) startListen() error {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// 允许所有来源
			return true
		},
	}
	path := "ws"
	if len(n.wsPath) > 0 {
		path = n.wsPath
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/"+path, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.ErrorNoStack(fmt.Sprintf("websocket conn failed %v", err))
			return
		}

		c, err := newWSConn(conn)
		if err != nil {
			logger.ErrorNoStack(fmt.Sprintf("new websocket conn failed %v", err))
			return
		}
		go func(conn net.Conn) {
			network.ServeSessionConn(conn, n.MessageCodec, n.IoDispatch, n.DispatchWorkers, n.PayloadMode)
		}(c)
	})

	listener, err := net.Listen("tcp", n.ServiceAddr)
	if err != nil {
		return err
	}
	n.listener = listener
	n.server = &http.Server{
		Addr:    n.ServiceAddr,
		Handler: mux,
	}
	go func() {
		if err := n.server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.ErrorNoStack(fmt.Sprintf("websocket server failed %v", err))
		}
	}()
	return nil
}

func (n *WsServer) Stop() {
	n.stopOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if n.server != nil {
			_ = n.server.Shutdown(ctx)
		}
		if n.listener != nil {
			_ = n.listener.Close()
		}

		network.CloseAllSessions()
	})
}

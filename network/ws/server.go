package ws

import (
	"context"
	"errors"
	"fmt"
	"io/github/gforgame/network"
	"io/github/gforgame/network/protocol"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

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

func NewServer(opts ...Option) *WsServer {
	opt := Options{}
	for _, option := range opts {
		option(&opt)
	}

	if opt.wsPath == "" {
		panic("ws path cannot be empty")
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

	modules := n.modules
	for _, c := range modules {
		c.Init()
		err := n.Router.RegisterMessageHandlers(c)
		if err != nil {
			return err
		}
	}

	if err := n.startListen(); err != nil {
		return err
	}

	return nil
}

func (n *WsServer) Addr() string {
	return n.ServiceAddr
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
			slog.Error(fmt.Sprintf("websocket conn failed %v", err))
			return
		}

		c, err := newWSConn(conn)
		if err != nil {
			slog.Error(fmt.Sprintf("new websocket conn failed %v", err))
			return
		}
		go onClientConnected(n, c)
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
			slog.Error(fmt.Sprintf("websocket server failed %v", err))
		}
	}()
	return nil
}

// 处理客户端连接，包括socket,websocket
func onClientConnected(node *WsServer, conn net.Conn) {
	defer func() {
		slog.Debug(fmt.Sprintf("客户端连接关闭: %s", conn.RemoteAddr().String()))
		// 处理客户端网络断开
		s := network.GetSession(conn)
		node.IoDispatch.OnSessionClosed(s)
		network.UnregisterSession(conn)
		_ = conn.Close()
	}()

	// 先创建默认的Session，协议类型会在第一次收到消息时确定
	var protocolType protocol.ProtocolType
	if _, ok := conn.(*wsConn); ok {
		// WebSocket连接，先使用二进制协议，后续会根据消息类型调整
		protocolType = protocol.ProtocolTypeBinary
		slog.Debug("WebSocket客户端连接，等待确定协议类型")
	} else {
		// TCP连接，默认使用二进制协议
		protocolType = protocol.ProtocolTypeBinary
		slog.Debug("TCP客户端使用二进制协议")
	}

	// 创建Session
	ioSession := network.NewSessionWithProtocol(conn, node.MessageCodec, protocolType)
	network.RegisterSession(conn, ioSession)

	// session created hook
	node.IoDispatch.OnSessionCreated(ioSession)

	// 异步读写数据
	go ioSession.Read()
	go ioSession.Write()

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
		for _, c := range n.modules {
			c.Shutdown()
		}
	})
}

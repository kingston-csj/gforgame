package ws

import (
	"errors"
	"fmt"
	"io/github/gforgame/logger"
	"io/github/gforgame/network"
	"io/github/gforgame/network/protocol"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

type WsServer struct {
	Options
	Name    string // 服务器名称
	Running chan bool
}

func NewServer(opts ...Option) *WsServer {
	opt := Options{}
	for _, option := range opts {
		option(&opt)
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
			panic(err)
		}
	}

	go func() {
		n.startListen()
	}()

	return nil
}

func (n *WsServer) Addr() string {
	return n.ServiceAddr
}

// Enable current server accept connection
func (n *WsServer) startListen() {
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
	http.HandleFunc("/"+path, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error(fmt.Errorf("websocket conn failed %v", err))
			return
		}

		c, err := newWSConn(conn)
		if err != nil {
			logger.Error(fmt.Errorf("new websocket conn failed %v", err))
			return
		}
		go onClientConnected(n, c)
	})
	if err := http.ListenAndServe(n.ServiceAddr, nil); err != nil {
		panic(err)
	}
}

// 处理客户端连接，包括socket,websocket
func onClientConnected(node *WsServer, conn net.Conn) {
	defer func() {
		// 处理客户端网络断开
		// s := network.GetSession(conn)
		// node.IoDispatch.OnSessionCreated(s)
		// network.UnregisterSession(conn)
		// err := conn.Close()
		// if err != nil {
		// 	logger.Error(fmt.Errorf("close ws conn failed %v", err))
		// }
		if r := recover(); r != nil {
			logger.Error(fmt.Errorf("Recovered from panic: %v", r))
			// 可能的话发送错误响应到客户端
			// 记录详细错误日志
		}
	}()

	// 先创建默认的Session，协议类型会在第一次收到消息时确定
	var protocolType protocol.ProtocolType
	if _, ok := conn.(*wsConn); ok {
		// WebSocket连接，先使用二进制协议，后续会根据消息类型调整
		protocolType = protocol.ProtocolTypeBinary
		logger.Debugf("WebSocket客户端连接，等待确定协议类型")
	} else {
		// TCP连接，默认使用二进制协议
		protocolType = protocol.ProtocolTypeBinary
		logger.Debugf("TCP客户端使用二进制协议")
	}

	// 创建Session
	ioSession := network.NewSessionWithProtocol(conn, node.MessageCodec, protocolType)
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

func (n *WsServer) Stop() {
	for _, c := range n.modules {
		c.Shutdown()
	}
}

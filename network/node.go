package network

import (
	"errors"
	"fmt"
	"io/github/gforgame/logger"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

type Node struct {
	Name    string // 服务器名称
	option  Options
	Running chan bool
	Router  MessageRoute
}

func (n *Node) Startup(opts ...Option) error {
	// 设置参数
	opt := Options{}
	for _, option := range opts {
		option(&opt)
	}
	n.option = opt

	if n.option.ServiceAddr == "" {
		return errors.New("service address cannot be empty in master node")
	}

	modules := n.option.modules
	for _, c := range modules {
		c.Init()
		err := n.Router.RegisterMessageHandlers(c)
		if err != nil {
			panic(err)
		}
	}

	go func() {
		if n.option.isWebsocket {
			n.listenWsConn()
		} else {
			n.listenTcpConn()
		}
	}()

	return nil
}

// Enable current server accept connection
func (n *Node) listenTcpConn() {
	listener, err := net.Listen("tcp", n.option.ServiceAddr)
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
		go handleClient(n, conn)
	}
}

// 处理客户端连接，包括socket,websocket
func handleClient(node *Node, conn net.Conn) {
	defer func() {
		// 处理客户端网络断开
		s := GetSession(conn)
		node.option.IoDispatch.OnSessionCreated(s)
		unregisterSession(conn)
		err := conn.Close()
		if err != nil {
			logger.Error(fmt.Errorf("close tcp conn failed %v", err))
		}
	}()

	ioSession := NewSession(conn, node.option.MessageCodec)
	registerSession(conn, ioSession)

	// session created hook
	node.option.IoDispatch.OnSessionCreated(ioSession)

	// 异步读写数据
	go ioSession.Read()
	go ioSession.Write()

	// read loop
	for ioFrame := range ioSession.DataReceived {
		node.option.IoDispatch.OnMessageReceived(ioSession, ioFrame)
	}
}

func (n *Node) listenWsConn() {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// 允许所有来源
			return true
		},
	}
	path := "ws"
	if len(n.option.wsPath) > 0 {
		path = n.option.wsPath
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
		go handleClient(n, c)
	})
	if err := http.ListenAndServe(n.option.ServiceAddr, nil); err != nil {
		panic(err)
	}
}

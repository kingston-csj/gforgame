package network

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"reflect"
)

type Node struct {
	Name   string // 服务器名称
	option Options
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
		log.Fatal(err.Error())
	}

	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}

		go handleClient(n, conn)
	}
}

// 处理客户端连接，包括socket,websocket
func handleClient(node *Node, conn net.Conn) {
	defer conn.Close() // 确保在函数结束时关闭连接

	ioSession := NewSession(&conn, node.option.MessageCodec)
	// 异步向客户端写数据
	go ioSession.Write()

	// read loop
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf(fmt.Sprintf("Read message error: %s, session will be closed immediately", err.Error()))
			return
		}
		if n <= 0 {
			continue
		}
		packets, err := ioSession.ProtocolCodec.Decode(buf[:n])
		if err != nil {
			log.Println(err.Error())
			return
		}
		// process packets decoded
		for _, p := range packets {
			typ, _ := GetMessageType(p.Header.Cmd)
			msg := reflect.New(typ.Elem()).Interface()
			err := node.option.MessageCodec.Decode(p.Data, msg)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			ioFrame := &RequestDataFrame{Header: p.Header, Msg: msg}
			node.option.IoDispatch.OnMessageReceived(ioSession, *ioFrame)
		}
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
			log.Println(fmt.Sprintf("Upgrade failure, URI=%s, Error=%s", r.RequestURI, err.Error()))
			return
		}

		c, err := newWSConn(conn)
		if err != nil {
			log.Println(err)
			return
		}
		go handleClient(n, c)
	})
	if err := http.ListenAndServe(n.option.ServiceAddr, nil); err != nil {
		log.Fatal(err.Error())
	}
}

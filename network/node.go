package network

import (
	"errors"
	"fmt"
	"log"
	"net"
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
	n.listenTcpConn()
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

// 处理客户端连接的函数
func handleClient(node *Node, conn net.Conn) {
	defer conn.Close() // 确保在函数结束时关闭连接

	ioSession := NewSession(&conn, node.option.MessageCodec)

	go ioSession.Write()

	// read loop
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf(fmt.Sprintf("Read message error: %s, session will be closed immediately", err.Error()))
			return
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
			node.option.MessageCodec.Decode(p.Data, msg)

			io_frame := &RequestDataFrame{Header: p.Header, Msg: msg}
			node.option.IoDispatch.OnMessageReceived(ioSession, *io_frame)
		}

	}
}

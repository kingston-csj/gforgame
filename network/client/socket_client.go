package client

import (
	"io/github/gforgame/codec"
	"io/github/gforgame/network"
	"log"
	"net"
)

type TcpSocketClient struct {

	// 服务器ip+端口
	RemoteAddress string
	// 消息编解码
	MsgCodec codec.MessageCodec
	// 消息分发器
	IoDispatcher network.IoDispatch
}

func (c *TcpSocketClient) OpenSession() (*network.Session, error) {
	conn, err := net.Dial("tcp", c.RemoteAddress)
	if err != nil {
		log.Fatal("连接服务器失败:", err)
	}
	// defer conn.Close()
	session := network.NewSession(conn, c.MsgCodec)
	go session.Write()
	go session.Read()
	return session, nil
}

package client

import (
	"net"

	"github.com/forfun/gforgame/codec"
	"github.com/forfun/gforgame/network"
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
		return nil, err
	}
	// defer conn.Close()
	session := network.NewSession(conn, c.MsgCodec)
	go session.Write()
	go session.Read()
	return session, nil
}

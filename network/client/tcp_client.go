package client

import (
	"fmt"
	"net"

	"github.com/forfun/gforgame/codec"
	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/dispatch"
	"github.com/forfun/gforgame/network/protocol"
	"github.com/forfun/gforgame/network/session"
)

type TcpSocketClient struct {

	// 服务器ip+端口
	RemoteAddress string
	// 消息编解码
	MsgCodec codec.MessageCodec
	// 消息分发器
	IoDispatcher dispatch.IoDispatch

	MessageFactory *protocol.MessageFactory
}

func (c *TcpSocketClient) OpenSession() (network.Session, error) {
	if c.MessageFactory == nil {
		return nil, fmt.Errorf("message factory cannot be nil")
	}
	conn, err := net.Dial("tcp", c.RemoteAddress)
	if err != nil {
		return nil, err
	}
	// defer conn.Close()
	s := session.NewSession(conn, c.MsgCodec, c.MessageFactory)
	go s.Write()
	go s.Read()
	return s, nil
}

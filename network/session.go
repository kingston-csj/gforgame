package network

import (
	"fmt"
	"io/github/gforgame/codec"
	"log"
	"net"
)

type Session struct {
	conn *net.Conn

	ProtocolCodec *Protocol

	// TODO 改为接口
	MessageCodec codec.MessageCodec

	attrs map[string]interface{}
	// (当前链接的本地地址)
	localAddr string
	// (当前链接的远程地址)
	remoteAddr string
}

func NewSession(conn *net.Conn, messageCodec codec.MessageCodec) *Session {
	return &Session{conn: conn,
		ProtocolCodec: NewDecoder(),
		MessageCodec:  messageCodec,
		localAddr:     (*conn).LocalAddr().String(),
		remoteAddr:    (*conn).RemoteAddr().String(),
	}
}

func (s *Session) Send(msg any) {
	msg_data, err := s.MessageCodec.Encode(msg)
	if err != nil {
		log.Fatal("连接服务器失败:", err)
	}

	cmd, _ := GetMessageCmd(msg)
	fmt.Println("发送消息:", cmd)
	frame, _ := s.ProtocolCodec.Encode(cmd, msg_data)
	(*s.conn).Write(frame)
}

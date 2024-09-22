package network

import (
	"fmt"
	"io/github/gforgame/codec"
	"log"
	"net"
)

type Session struct {
	conn *net.Conn

	die chan bool

	ProtocolCodec *Protocol

	MessageCodec codec.MessageCodec

	dataToSend chan []byte

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
		dataToSend:    make(chan []byte),
		localAddr:     (*conn).LocalAddr().String(),
		remoteAddr:    (*conn).RemoteAddr().String(),
	}
}

func (s *Session) Send(msg any) error {
	if msg == nil {
		return nil
	}
	msg_data, err := s.MessageCodec.Encode(msg)
	if err != nil {
		return fmt.Errorf("encode message %s cmd failed", msg)
	}

	cmd, e2 := GetMessageCmd(msg)
	if e2 != nil {
		return fmt.Errorf("get message %s cmd failed:%v", msg, e2)
	}
	fmt.Println("发送消息:", cmd)
	frame, _ := s.ProtocolCodec.Encode(cmd, msg_data)
	s.dataToSend <- frame
	return nil
}

func (s *Session) Write() {
	defer close(s.dataToSend)

	for {
		select {
		case data := <-s.dataToSend:
			if _, err := (*s.conn).Write(data); err != nil {
				log.Println(err.Error())
				//s.Close()
			}
		case <-s.die:
			return
		}
	}
}

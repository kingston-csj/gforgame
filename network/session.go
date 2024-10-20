package network

import (
	"fmt"
	"io/github/gforgame/codec"
	"io/github/gforgame/logger"
	"io/github/gforgame/network/protocol"
	"log"
	"net"
	"reflect"
)

type Session struct {
	conn net.Conn

	die chan bool

	ProtocolCodec *protocol.Protocol

	MessageCodec codec.MessageCodec

	dataToSend chan []byte

	DataReceived chan *protocol.RequestDataFrame

	Attrs map[string]interface{}
	// (当前链接的本地地址)
	localAddr string
	// (当前链接的远程地址)
	remoteAddr string
}

func NewSession(conn net.Conn, messageCodec codec.MessageCodec) *Session {
	return &Session{conn: conn,
		ProtocolCodec: protocol.NewDecoder(),
		MessageCodec:  messageCodec,
		dataToSend:    make(chan []byte, 128),
		DataReceived:  make(chan *protocol.RequestDataFrame, 128),
		Attrs:         map[string]interface{}{},
		localAddr:     conn.LocalAddr().String(),
		remoteAddr:    conn.RemoteAddr().String(),
	}
}

func (s *Session) Send(msg any, index int) error {
	if msg == nil {
		return nil
	}
	msgData, err := s.MessageCodec.Encode(msg)
	if err != nil {
		return fmt.Errorf("encode message %s cmd failed", msg)
	}

	cmd, e2 := GetMessageCmd(msg)
	if e2 != nil {
		return fmt.Errorf("get message %s cmd failed:%v", msg, e2)
	}
	fmt.Println("发送消息:", cmd)
	frame, _ := s.ProtocolCodec.Encode(cmd, index, msgData)
	s.dataToSend <- frame
	return nil
}

func (s *Session) Write() {
	defer close(s.dataToSend)

	for {
		select {
		case data := <-s.dataToSend:
			if _, err := s.conn.Write(data); err != nil {
				log.Println(err.Error())
				//s.Close()
			}
		case <-s.die:
			return
		}
	}
}

func (s *Session) Read() {
	buf := make([]byte, 2048)
	for {
		n, err := s.conn.Read(buf)
		if err != nil {
			log.Println(err.Error())
			return
		}
		if n <= 0 {
			continue
		}
		packets, err := s.ProtocolCodec.Decode(buf[:n])
		if err != nil {
			log.Println(fmt.Errorf("decode protocol  failed %v", err))
			return
		}
		// process packets decoded
		for _, p := range packets {
			typ, _ := GetMessageType(p.Header.Cmd)
			msg := reflect.New(typ.Elem()).Interface()
			err := s.MessageCodec.Decode(p.Data, msg)
			if err != nil {
				logger.Error(fmt.Errorf("decode message  failed %v", err))
				continue
			}
			ioFrame := &protocol.RequestDataFrame{Header: p.Header, Msg: msg}
			s.DataReceived <- ioFrame
		}
	}
}

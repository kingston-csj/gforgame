package session

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"strings"

	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/network/protocol"

	"github.com/gorilla/websocket"
)

// WebSocketConn WebSocket连接接口
type WebSocketConn interface {
	net.Conn
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
}

func (s *Session) Write() {
	for {
		select {
		case data := <-s.dataToSend:
			if _, err := s.conn.Write(data); err != nil {
				logger.ErrorNoStack(fmt.Sprintf("session write failed %v", err))
				s.Close()
				return
			}
		case <-s.Die:
			return
		}
	}
}

func (s *Session) Read() {
	defer func() {
		s.Close()
		if r := recover(); r != nil {
			logger.ErrorNoStack(fmt.Sprintf("panic recovered: %v", r))
		}
	}()
	if wsConn, ok := s.conn.(WebSocketConn); ok {
		s.readWebSocketMessages(wsConn)
	} else {
		s.readTCPStream()
	}
}

func (s *Session) readWebSocketMessages(wsConn WebSocketConn) {
	protocolDetermined := false

	for {
		messageType, messageData, err := wsConn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
				return
			}
			errMsg := err.Error()
			if strings.Contains(errMsg, "close 1006") || strings.Contains(errMsg, "unexpected EOF") {
				return
			}
			logger.ErrorNoStack(fmt.Sprintf("websocket read failed %v", err))
			return
		}
		s.MarkReadActivity()

		select {
		case <-s.Die:
			return
		default:
		}

		if !protocolDetermined {
			var newProtocolType protocol.ProtocolType
			if messageType == websocket.TextMessage {
				newProtocolType = protocol.ProtocolTypeJSON
			} else {
				newProtocolType = protocol.ProtocolTypeBinary
			}

			if s.protocolType != newProtocolType {
				factory := &protocol.ProtocolFactory{}
				s.ProtocolCodec = factory.NewProtocolAdapter(newProtocolType)
				s.protocolType = newProtocolType
			}
			protocolDetermined = true
		}

		packets, err := s.ProtocolCodec.Decode(messageData)
		if err != nil {
			logger.ErrorNoStack(fmt.Errorf("decode protocol failed %v", err))
			if errors.Is(err, protocol.ErrInvalidPacketSize) || errors.Is(err, protocol.ErrPacketSizeExceed) {
				return
			}
			continue
		}

		if s.payloadMode == PayloadModeRawBody {
			for _, p := range packets {
				ioFrame := &protocol.RequestDataFrame{Header: p.Header, Msg: p.Data}
				s.DataReceived <- ioFrame
			}
			continue
		}

		for _, p := range packets {
			typ, _ := messageResolver.GetMessageType(p.Header.Cmd)
			if typ == nil {
				logger.ErrorNoStack(fmt.Sprintf("message type not found %v", p.Header.Cmd))
				continue
			}
			msg := reflect.New(typ.Elem()).Interface()
			err := s.MessageCodec.Decode(p.Data, msg)
			if err != nil {
				logger.ErrorNoStack(fmt.Sprintf("decode message failed %v", err))
				continue
			}
			ioFrame := &protocol.RequestDataFrame{Header: p.Header, Msg: msg}
			s.DataReceived <- ioFrame
		}
	}
}

func (s *Session) readTCPStream() {
	buf := make([]byte, 10240)

	for {
		select {
		case <-s.Die:
			return
		default:
		}

		n, err := s.conn.Read(buf)
		if err != nil {
			logger.ErrorNoStack(err.Error())
			return
		}
		if n <= 0 {
			continue
		}
		s.MarkReadActivity()
		packets, err := s.ProtocolCodec.Decode(buf[:n])
		if err != nil {
			logger.ErrorNoStack(fmt.Errorf("decode protocol failed %v", err))
			return
		}
		for _, p := range packets {
			if s.payloadMode == PayloadModeRawBody {
				ioFrame := &protocol.RequestDataFrame{Header: p.Header, Msg: p.Data}
				s.DataReceived <- ioFrame
				continue
			}
			typ, _ := messageResolver.GetMessageType(p.Header.Cmd)
			if typ == nil {
				logger.ErrorNoStack(fmt.Sprintf("message type not found %v", p.Header.Cmd))
				continue
			}
			msg := reflect.New(typ.Elem()).Interface()
			err := s.MessageCodec.Decode(p.Data, msg)
			if err != nil {
				logger.ErrorNoStack(fmt.Sprintf("decode message failed %v", err))
				continue
			}
			ioFrame := &protocol.RequestDataFrame{Header: p.Header, Msg: msg}
			s.DataReceived <- ioFrame
		}
	}
}

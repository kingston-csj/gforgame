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

func (s *BaseSession) Write() {
	defer func() {
		if r := recover(); r != nil {
			logger.ErrorNoStack(fmt.Sprintf("session write panic %v", r))
		}
	}()
	for {
		select {
		case data := <-s.dataToSend:
			_, err := s.conn.Write(data.frame)
			if data.done != nil {
				data.done <- err
			}
			if err != nil {
				logger.ErrorNoStack(fmt.Sprintf("session write failed %v", err))
				s.Close()
				return
			}
		case <-s.Die:
			return
		}
	}
}

func (s *BaseSession) Read() {
	defer func() {
		s.Close()
		if r := recover(); r != nil {
			logger.ErrorNoStack(fmt.Sprintf("panic recovered: %v", r))
		}
	}()

	if wsConn, ok := s.conn.(WebSocketConn); ok {
		s.readWebSocketStream(wsConn)
	} else {
		s.readTcpStream()
	}
}

// 【公共抽取：解析数据包并投递到DataReceived，统一处理Die信号】
func (s *BaseSession) decodeAndDeliverPackets(packets []*protocol.Packet) {
	for _, p := range packets {
		var ioFrame *protocol.RequestDataFrame
		if s.payloadMode == protocol.PayloadModeRawBody {
			ioFrame = &protocol.RequestDataFrame{
				Header: p.Header,
				Msg:    p.Data,
			}
		} else {
			typ, _ := protocol.GetMessageType(p.Header.Cmd)
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
			ioFrame = &protocol.RequestDataFrame{
				Header: p.Header,
				Msg:    msg,
			}
		}

		// 投递前监听Die，防止向关闭的DataReceived写入panic
		select {
		case <-s.Die:
			return
		case s.DataReceived <- ioFrame:
		}
	}
}

func (s *BaseSession) readWebSocketStream(wsConn WebSocketConn) {
	protocolDetermined := false
	for {
		// WebSocket ReadMessage是阻塞调用，无法被Die直接唤醒，依赖底层conn关闭唤醒
		messageType, messageData, err := wsConn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseNoStatusReceived) {
				return
			}
			errMsg := err.Error()
			if strings.Contains(errMsg, "close 1006") || strings.Contains(errMsg, "unexpected EOF") {
				return
			}
			logger.ErrorNoStack(fmt.Sprintf("websocket read failed %v", err))
			return
		}

		// 读到数据后先检查Die
		select {
		case <-s.Die:
			return
		default:
		}

		s.MarkReadActivity()

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
		s.decodeAndDeliverPackets(packets)
	}
}

func (s *BaseSession) readTcpStream() {
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
		s.decodeAndDeliverPackets(packets)
	}
}


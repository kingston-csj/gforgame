package network

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"reflect"
	"strings"
	"sync"

	"github.com/forfun/gforgame/codec"
	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/common/util/jsonutil"
	"github.com/forfun/gforgame/network/protocol"

	"github.com/gorilla/websocket"
)

// WebSocketConn WebSocket连接接口
type WebSocketConn interface {
	net.Conn
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
}

type Session struct {
	conn net.Conn
	// 关闭标记
	Die chan bool
	// 私有协议栈编解码
	ProtocolCodec protocol.ProtocolAdapter
	// 消息编解码
	MessageCodec codec.MessageCodec
	// 准备发送的出队消息(带缓冲)
	dataToSend chan []byte
	// 已经收到的入队消息(带缓冲)
	DataReceived chan *protocol.RequestDataFrame
	// session自定义属性
	Attrs map[string]interface{}
	// 当前链接的本地地址
	localAddr string
	// 当前链接的远程地址
	remoteAddr string
	// 异步任务
	AsynTasks chan func()
	// 协议类型
	protocolType protocol.ProtocolType
	// 关闭只执行一次
	closeOnce sync.Once
}

func NewSession(conn net.Conn, messageCodec codec.MessageCodec) *Session {
	return &Session{conn: conn,
		ProtocolCodec: protocol.NewBinaryProtocolAdapter(),
		MessageCodec:  messageCodec,
		Die:           make(chan bool, 1),
		dataToSend:    make(chan []byte, 128),
		DataReceived:  make(chan *protocol.RequestDataFrame, 128),
		Attrs:         map[string]interface{}{},
		localAddr:     conn.LocalAddr().String(),
		remoteAddr:    conn.RemoteAddr().String(),
		AsynTasks:     make(chan func(), 16),
		protocolType:  protocol.ProtocolTypeBinary,
	}
}

// NewSessionWithProtocol 创建指定协议类型的Session
func NewSessionWithProtocol(conn net.Conn, messageCodec codec.MessageCodec, protocolType protocol.ProtocolType) *Session {
	factory := &protocol.ProtocolFactory{}
	protocolAdapter := factory.NewProtocolAdapter(protocolType)

	return &Session{conn: conn,
		ProtocolCodec: protocolAdapter,
		MessageCodec:  messageCodec,
		Die:           make(chan bool, 1),
		dataToSend:    make(chan []byte, 128),
		DataReceived:  make(chan *protocol.RequestDataFrame, 128),
		Attrs:         map[string]interface{}{},
		localAddr:     conn.LocalAddr().String(),
		remoteAddr:    conn.RemoteAddr().String(),
		AsynTasks:     make(chan func(), 16),
		protocolType:  protocolType,
	}
}

// Send 发送消息
// @param msg 要发送的消息
// @param index 消息索引
// @return 发送过程遇到的异常
func (s *Session) Send(msg any, index int32) error {
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

	msgName, e3 := GetMsgName(cmd)
	if e3 != nil {
		return fmt.Errorf("get message %s name failed:%v", msg, e3)
	}
	jsonStr, err := jsonutil.StructToJSON(msg)
	id, ok := s.GetAttr("id")
	if !ok {
		id = ""
	}
	if err == nil {
		if !strings.Contains(msgName, "Heartbeat") {
			logger.Info(fmt.Sprintf("id:%v 发送消息 cmd:%d, name:%s, 内容:%s", id, cmd, msgName, jsonStr))
		}
	}
	frame, _ := s.ProtocolCodec.Encode(cmd, int32(index), msgData)
	select {
	case <-s.Die:
		return errors.New("session closed")
	case s.dataToSend <- frame:
		return nil
	}
}

func (s *Session) SendWithoutIndex(msg any) error {
	return s.Send(msg, 0)
}

func (s *Session) SetAttr(key string, value any) error {
	s.Attrs[key] = value
	return nil
}

func (s *Session) GetAttr(key string) (any, bool) {
	value, ok := s.Attrs[key]
	return value, ok
}

// SendAndClose 发送消息并关闭连接
// 同步阻塞，消息发送完毕，随即关闭连接
// 注意：执行完毕仅代表数据已写入本地内核缓冲区，并不保证客户端一定会收到
// @param msg 要发送的消息
// @return 发送过程遇到的异常
func (s *Session) SendAndClose(msg any) error {
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
	msgName, e3 := GetMsgName(cmd)
	if e3 != nil {
		return fmt.Errorf("get message %s name failed:%v", msg, e3)
	}
	id, ok := s.GetAttr("id")
	if !ok {
		id = ""
	}
	fmt.Println(fmt.Sprintf("id:%s 发送消息 cmd:%d, name:%s, 内容:%v", id, cmd, msgName, msg))
	frame, _ := s.ProtocolCodec.Encode(cmd, int32(-1), msgData)
	_, err = s.conn.Write(frame)
	if err != nil {
		return err
	}
	// 关闭连接
	err = s.conn.Close()
	if err != nil {
		return err
	}
	return err
}

func (s *Session) Write() {
	for {
		select {
		case data := <-s.dataToSend:
			if _, err := s.conn.Write(data); err != nil {
				slog.Error("session write failed", "error", err)
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
		// 一旦关闭，onClientConnected 就会收到断开信号
		s.Close()
		if r := recover(); r != nil {
			slog.Error(fmt.Sprintf("panic recovered: %v", r))
		}
	}()
	// 检查是否是WebSocket连接
	if wsConn, ok := s.conn.(WebSocketConn); ok {
		// WebSocket连接，按消息处理
		s.readWebSocketMessages(wsConn)
	} else {
		// TCP连接，按字节流处理
		s.readTCPStream()
	}
}

// Close 关闭会话（幂等）
func (s *Session) Close() {
	s.closeOnce.Do(func() {
		close(s.Die)
		_ = s.conn.Close()
	})
}

// readWebSocketMessages 处理WebSocket消息
func (s *Session) readWebSocketMessages(wsConn WebSocketConn) {
	protocolDetermined := false // 标记是否已经确定协议类型

	for {
		// 读取一条完整的WebSocket消息
		messageType, messageData, err := wsConn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
				slog.Debug(fmt.Sprintf("websocket normal close: %v", err))
				return
			}
			// Unity 直接关闭进程时常见 1006/unexpected EOF，按正常断链处理
			errMsg := err.Error()
			if strings.Contains(errMsg, "close 1006") || strings.Contains(errMsg, "unexpected EOF") {
				slog.Debug(fmt.Sprintf("websocket peer disconnected: %v", err))
				return
			}
			slog.Error("websocket read failed", "error", err)
			return
		}

		// 检查是否关闭，结束该goroutine
		select {
		case <-s.Die:
			// 关闭连接
			return
		default:
		}

		// 第一次收到消息时确定协议类型并调整协议适配器
		if !protocolDetermined {
			var newProtocolType protocol.ProtocolType
			if messageType == websocket.TextMessage {
				newProtocolType = protocol.ProtocolTypeJSON
				slog.Debug("WebSocket客户端使用JSON协议")
			} else {
				newProtocolType = protocol.ProtocolTypeBinary
				slog.Debug("WebSocket客户端使用二进制协议")
			}

			// 如果协议类型发生变化，创建新的协议适配器
			if s.protocolType != newProtocolType {
				factory := &protocol.ProtocolFactory{}
				s.ProtocolCodec = factory.NewProtocolAdapter(newProtocolType)
				s.protocolType = newProtocolType
				slog.Debug(fmt.Sprintf("协议适配器已切换: %v", newProtocolType))
			}
			protocolDetermined = true
		}

		// 使用对应的协议解码器处理消息
		packets, err := s.ProtocolCodec.Decode(messageData)
		if err != nil {
			log.Println(fmt.Errorf("decode protocol failed %v", err))
			// WebSocket消息错误时继续处理下一条消息
			continue
		}

		// 处理解码后的数据
		for _, p := range packets {
			typ, _ := GetMessageType(p.Header.Cmd)
			if typ == nil {
				slog.Error(fmt.Sprintf("message type not found %v", p.Header.Cmd))
				continue
			}
			msg := reflect.New(typ.Elem()).Interface()
			err := s.MessageCodec.Decode(p.Data, msg)
			if err != nil {
				slog.Error(fmt.Sprintf("decode message failed %v", err))
				continue
			}
			ioFrame := &protocol.RequestDataFrame{Header: p.Header, Msg: msg}
			s.DataReceived <- ioFrame
		}
	}
}

// readTCPStream 处理TCP字节流
func (s *Session) readTCPStream() {
	buf := make([]byte, 10240)

	for {
		// 检查是否关闭，结束该goroutine
		select {
		case <-s.Die:
			return
		default:
		}

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
			log.Println(fmt.Errorf("decode protocol failed %v", err))
			return
		}
		// process packets decoded
		for _, p := range packets {
			typ, _ := GetMessageType(p.Header.Cmd)
			if typ == nil {
				slog.Error(fmt.Sprintf("message type not found %v", p.Header.Cmd))
				continue
			}
			msg := reflect.New(typ.Elem()).Interface()
			err := s.MessageCodec.Decode(p.Data, msg)
			if err != nil {
				slog.Error(fmt.Sprintf("decode message failed %v", err))
				continue
			}
			ioFrame := &protocol.RequestDataFrame{Header: p.Header, Msg: msg}
			s.DataReceived <- ioFrame
		}
	}
}

func (s *Session) ToString() string {
	id, ok := s.GetAttr("id")
	if !ok {
		id = "anonymous"
	}
	// id + remoteAddr
	return fmt.Sprintf("id:%s, remoteAddr:%s", id, s.conn.RemoteAddr().String())
}

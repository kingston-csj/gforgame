package session

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/forfun/gforgame/codec"
	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/common/util/jsonutil"
	"github.com/forfun/gforgame/network/protocol"
)

// Session 封装单条连接的运行时状态与基础发送能力。
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
	// 消息体处理模式
	payloadMode PayloadMode
	// 关闭只执行一次
	closeOnce sync.Once
	// 最近一次收到客户端数据的时间
	lastRecvUnixNano int64
}

func NewSession(conn net.Conn, messageCodec codec.MessageCodec) *Session {
	nowUnixNano := time.Now().UnixNano()
	return &Session{conn: conn,
		ProtocolCodec:    protocol.NewBinaryProtocolAdapter(),
		MessageCodec:     messageCodec,
		Die:              make(chan bool, 1),
		dataToSend:       make(chan []byte, 128),
		DataReceived:     make(chan *protocol.RequestDataFrame, 128),
		Attrs:            map[string]interface{}{},
		localAddr:        conn.LocalAddr().String(),
		remoteAddr:       conn.RemoteAddr().String(),
		AsynTasks:        make(chan func(), 16),
		protocolType:     protocol.ProtocolTypeBinary,
		payloadMode:      PayloadModeDecode,
		lastRecvUnixNano: nowUnixNano,
	}
}

// NewSessionWithProtocol 创建指定协议类型的Session
func NewSessionWithProtocol(conn net.Conn, messageCodec codec.MessageCodec, protocolType protocol.ProtocolType) *Session {
	factory := &protocol.ProtocolFactory{}
	protocolAdapter := factory.NewProtocolAdapter(protocolType)
	nowUnixNano := time.Now().UnixNano()

	return &Session{
		conn:             conn,
		ProtocolCodec:    protocolAdapter,
		MessageCodec:     messageCodec,
		Die:              make(chan bool, 1),
		dataToSend:       make(chan []byte, 128),
		DataReceived:     make(chan *protocol.RequestDataFrame, 128),
		Attrs:            map[string]interface{}{},
		localAddr:        conn.LocalAddr().String(),
		remoteAddr:       conn.RemoteAddr().String(),
		AsynTasks:        make(chan func(), 16),
		protocolType:     protocolType,
		payloadMode:      PayloadModeDecode,
		lastRecvUnixNano: nowUnixNano,
	}
}

func (s *Session) MarkReadActivity() {
	atomic.StoreInt64(&s.lastRecvUnixNano, time.Now().UnixNano())
}

func (s *Session) LastReadAt() time.Time {
	return time.Unix(0, atomic.LoadInt64(&s.lastRecvUnixNano))
}

func (s *Session) SetPayloadMode(mode PayloadMode) {
	s.payloadMode = mode
}

// Send 发送消息
func (s *Session) Send(msg any, index int32) error {
	if msg == nil {
		return nil
	}
	msgData, err := s.MessageCodec.Encode(msg)
	if err != nil {
		logger.ErrorNoStack(fmt.Errorf("send message failed, fallback to reflect: cmd=%d err=%v", index, err))
		return fmt.Errorf("encode message %s cmd failed", msg)
	}

	cmd, e2 := messageResolver.GetMessageCmd(msg)
	if e2 != nil {
		logger.ErrorNoStack(fmt.Errorf("send message failed, fallback to reflect: cmd=%d err=%v", cmd, e2))
		return fmt.Errorf("get message %s cmd failed:%v", msg, e2)
	}

	msgName, e3 := messageResolver.GetMsgName(cmd)
	if e3 != nil {
		logger.ErrorNoStack(fmt.Errorf("send message failed, fallback to reflect: cmd=%d err=%v", cmd, e3))
		return fmt.Errorf("get message %s name failed:%v", msg, e3)
	}
	jsonStr, err := jsonutil.StructToJSON(msg)
	id, ok := s.GetAttr("id")
	if !ok {
		id = ""
	}
	if err == nil {
		if cmd != -101 && cmd != -300 {
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

// SendRaw 发送原始消息
func (s *Session) SendRaw(frame []byte) error {
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
func (s *Session) SendAndClose(msg any) error {
	if msg == nil {
		return nil
	}
	msgData, err := s.MessageCodec.Encode(msg)
	if err != nil {
		return fmt.Errorf("encode message %s cmd failed", msg)
	}

	cmd, e2 := messageResolver.GetMessageCmd(msg)
	if e2 != nil {
		logger.ErrorNoStack(fmt.Errorf("send message failed, fallback to reflect: cmd=%d err=%v", cmd, e2))
		return fmt.Errorf("get message %s cmd failed:%v", msg, e2)
	}
	msgName, e3 := messageResolver.GetMsgName(cmd)
	if e3 != nil {
		logger.ErrorNoStack(fmt.Errorf("send message failed, fallback to reflect: cmd=%d err=%v", cmd, e3))
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
	err = s.conn.Close()
	if err != nil {
		return err
	}
	return err
}

// Close 关闭会话（幂等）
func (s *Session) Close() {
	s.closeOnce.Do(func() {
		close(s.Die)
		_ = s.conn.Close()
	})
}

func (s *Session) ToString() string {
	id, ok := s.GetAttr("id")
	if !ok {
		id = "anonymous"
	}
	return fmt.Sprintf("id:%s, remoteAddr:%s", id, s.conn.RemoteAddr().String())
}

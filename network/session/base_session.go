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
	"github.com/google/uuid"
)

type queuedWrite struct {
	frame []byte
	done  chan error
}

// BaseSession 封装单条连接的运行时状态与基础发送能力。
type BaseSession struct {
	id ID
	ownerId OwnerID
	conn net.Conn
	// 关闭标记
	Die chan bool
	// 私有协议栈编解码
	ProtocolCodec protocol.ProtocolAdapter
	// 消息编解码
	MessageCodec codec.MessageCodec
	// 准备发送的出队消息(带缓冲)
	dataToSend chan queuedWrite
	// 已经收到的入队消息(带缓冲)
	DataReceived chan *protocol.RequestDataFrame
	// session自定义属性
	Attrs map[string]any
	// 当前链接的本地地址
	localAddr string
	// 当前链接的远程地址
	remoteAddr string
	// 协议类型
	protocolType protocol.ProtocolType
	// 消息体处理模式，仅初始化赋值一次，运行时只读
	payloadMode protocol.PayloadMode

	mu sync.RWMutex // 仅保护 Attrs map

	// 关闭只执行一次
	closeOnce sync.Once
	// 最近一次收到客户端数据的时间
	lastRecvUnixNano int64
}

func NewSession(conn net.Conn, messageCodec codec.MessageCodec) *BaseSession {
	nowUnixNano := time.Now().UnixNano()
	return &BaseSession{
		conn:             conn,
		ProtocolCodec:    protocol.NewBinaryProtocolAdapter(),
		MessageCodec:     messageCodec,
		Die:              make(chan bool, 1),
		dataToSend:       make(chan queuedWrite, 128),
		DataReceived:     make(chan *protocol.RequestDataFrame, 128),
		Attrs:            make(map[string]any),
		localAddr:        conn.LocalAddr().String(),
		remoteAddr:       conn.RemoteAddr().String(),
		protocolType:     protocol.ProtocolTypeBinary,
		payloadMode:      protocol.PayloadModeDecode,
		lastRecvUnixNano: nowUnixNano,
	}
}

// NewSessionWithProtocol 创建指定协议类型的Session
func NewSessionWithProtocol(conn net.Conn, messageCodec codec.MessageCodec, protocolType protocol.ProtocolType) *BaseSession {
	factory := &protocol.ProtocolFactory{}
	protocolAdapter := factory.NewProtocolAdapter(protocolType)
	nowUnixNano := time.Now().UnixNano()
	s := &BaseSession{
		conn:             conn,
		ProtocolCodec:    protocolAdapter,
		MessageCodec:     messageCodec,
		Die:              make(chan bool, 1),
		dataToSend:       make(chan queuedWrite, 128),
		DataReceived:     make(chan *protocol.RequestDataFrame, 128),
		Attrs:            make(map[string]any),
		localAddr:        conn.LocalAddr().String(),
		remoteAddr:       conn.RemoteAddr().String(),
		protocolType:     protocolType,
		payloadMode:      protocol.PayloadModeDecode,
		lastRecvUnixNano: nowUnixNano,
	}
	s.id = uuid.NewString()
	BindSID(s)
	return s
}

func (s *BaseSession) MarkReadActivity() {
	atomic.StoreInt64(&s.lastRecvUnixNano, time.Now().UnixNano())
}

func (s *BaseSession) LastReadAt() time.Time {
	return time.Unix(0, atomic.LoadInt64(&s.lastRecvUnixNano))
}

// SetPayloadMode 仅允许在session初始化阶段调用，运行时禁止调用
func (s *BaseSession) SetPayloadMode(mode protocol.PayloadMode) {
	s.payloadMode = mode
}

// Send 发送消息
func (s *BaseSession) Send(msg any, index int32) error {
	if msg == nil {
		return nil
	}
	msgData, err := s.MessageCodec.Encode(msg)
	if err != nil {
		logger.ErrorNoStack(fmt.Errorf("send message failed, fallback to reflect: cmd=%d err=%v", index, err))
		return fmt.Errorf("encode message %s cmd failed", msg)
	}
	cmd, e2 := protocol.GetMessageCmd(msg)
	if e2 != nil {
		logger.ErrorNoStack(fmt.Errorf("send message failed, fallback to reflect: cmd=%d err=%v", cmd, e2))
		return fmt.Errorf("get message %s cmd failed:%v", msg, e2)
	}
	msgName, e3 := protocol.GetMsgName(cmd)
	if e3 != nil {
		logger.ErrorNoStack(fmt.Errorf("send message failed, fallback to reflect: cmd=%d err=%v", cmd, e3))
		return fmt.Errorf("get message %s name failed:%v", msg, e3)
	}
	jsonStr, err := jsonutil.StructToJSON(msg)
	id := s.GetOwnerId()
	if id == "" {
		id = "anonymous"
	}
	if err == nil {
		if cmd != -151 && cmd != -300 {
			logger.Info(fmt.Sprintf("[%s] 发送消息:  cmd:%d, name:%s, 内容:%s", id, cmd, msgName, jsonStr))
		}
	}
	frame, _ := s.ProtocolCodec.Encode(cmd, int32(index), msgData)
	return s.enqueueFrame(frame, nil)
}

// SendRaw 发送原始消息
func (s *BaseSession) SendRaw(frame []byte) error {
	return s.enqueueFrame(frame, nil)
}

func (s *BaseSession) enqueueFrame(frame []byte, done chan error) error {
	select {
	case <-s.Die:
		return errors.New("session closed")
	case s.dataToSend <- queuedWrite{frame: frame, done: done}:
		return nil
	}
}

func (s *BaseSession) SendWithoutIndex(msg any) error {
	return s.Send(msg, 0)
}

// SetAttr 设置属性，写锁
func (s *BaseSession) SetAttr(key string, value any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Attrs[key] = value
	return nil
}

// GetAttr 获取属性，读锁
func (s *BaseSession) GetAttr(key string) (any, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.Attrs[key]
	return value, ok
}

// SendAndClose 发送消息并关闭连接
// 同步阻塞，等待消息按发送队列写完后，再关闭连接
// 注意：执行完毕仅代表数据已写入本地内核缓冲区，并不保证客户端一定会收到
// @param msg 要发送的消息
// @return 发送过程遇到的异常
func (s *BaseSession) SendAndClose(msg any) error {
	if msg == nil {
		return nil
	}
	msgData, err := s.MessageCodec.Encode(msg)
	if err != nil {
		return fmt.Errorf("encode message %s cmd failed", msg)
	}
	cmd, e2 := protocol.GetMessageCmd(msg)
	if e2 != nil {
		logger.ErrorNoStack(fmt.Errorf("send message failed, fallback to reflect: cmd=%d err=%v", cmd, e2))
		return fmt.Errorf("get message %s cmd failed:%v", msg, e2)
	}
	msgName, e3 := protocol.GetMsgName(cmd)
	if e3 != nil {
		logger.ErrorNoStack(fmt.Errorf("send message failed, fallback to reflect: cmd=%d err=%v", cmd, e3))
		return fmt.Errorf("get message %s name failed:%v", msg, e3)
	}
	logger.Info(fmt.Sprintf("[%s] 发送消息 cmd:%d, name:%s, 内容:%v", s.GetOwnerId(), cmd, msgName, msg))
	frame, _ := s.ProtocolCodec.Encode(cmd, int32(-1), msgData)
	done := make(chan error, 1)
	if err = s.enqueueFrame(frame, done); err != nil {
		return err
	}
	select {
	case err = <-done:
		if err != nil {
			return err
		}
	case <-s.Die:
		select {
		case err = <-done:
			if err != nil {
				return err
			}
		default:
			return errors.New("session closed")
		}
	}
	s.Close()
	return nil
}

// GetId 读锁读取id属性
func (s *BaseSession) GetId() string {
	return s.id
}

// GetOwnerId 读锁读取uid属性
func (s *BaseSession) GetOwnerId() OwnerID {
	return s.ownerId
}

// SetOwnerId 设置uid属性，写锁
func (s *BaseSession) SetOwnerId(id OwnerID) {
	s.ownerId = id
}

// DieChan 会话关闭广播通道（仅用于读取）。
func (s *BaseSession) DieChan() <-chan bool {
	return s.Die
}

// DataReceivedChan 入站消息队列通道（仅用于读取）。
func (s *BaseSession) DataReceivedChan() <-chan *protocol.RequestDataFrame {
	return s.DataReceived
}

// RawConn 返回底层连接。
func (s *BaseSession) RawConn() net.Conn {
	return s.conn
}

// GetProtocolCodec 返回私有协议栈编解码器。
func (s *BaseSession) GetProtocolCodec() protocol.ProtocolAdapter {
	return s.ProtocolCodec
}

// GetMessageCodec 返回消息编解码器。
func (s *BaseSession) GetMessageCodec() codec.MessageCodec {
	return s.MessageCodec
}

// Close 关闭会话（幂等）
func (s *BaseSession) Close() {
	s.closeOnce.Do(func() {
		close(s.Die)
		_ = s.conn.Close()
	})
}

func (s *BaseSession) IsAlive() bool {
	select {
	case <-s.DieChan():
		return false
	default:
		return true
	}
}

func (s *BaseSession) ToString() string {
	id := s.GetOwnerId()
	if id == "" {
		id = "anonymous"
	}
	return fmt.Sprintf("id:%s, remoteAddr:%s", id, s.conn.RemoteAddr().String())
}

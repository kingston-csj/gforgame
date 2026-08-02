package session

import (
	"net"
	"time"

	"github.com/forfun/gforgame/codec"
	"github.com/forfun/gforgame/network/protocol"
)

// Session 抽象一条连接上的会话能力，面向业务层暴露。
type Session interface {
	// Read 启动读循环（应作为 goroutine 启动）。
	Read()
	// Write 启动写循环（应作为 goroutine 启动）。
	Write()

	// MarkReadActivity 标记一次读取活动，用于空闲超时判断。
	MarkReadActivity()
	// LastReadAt 返回最近一次读取时间。
	LastReadAt() time.Time

	// SetPayloadMode 切换消息体解码模式。
	SetPayloadMode(mode PayloadMode)

	// Send 发送一条经过 MessageCodec + ProtocolCodec 编码的消息，带上请求索引。
	Send(msg any, index int32) error
	// SendWithoutIndex 等价于 Send(msg, 0)。
	SendWithoutIndex(msg any) error
	// SendRaw 发送已经打包好的原始帧（跳过上层编解码）。
	SendRaw(frame []byte) error
	// SendAndClose 同步发送一条消息后关闭会话。
	SendAndClose(msg any) error

	// SetAttr 设置自定义属性。
	SetAttr(key string, value any) error
	// GetAttr 读取自定义属性。
	GetAttr(key string) (any, bool)
	// GetId 读取会话绑定的玩家/节点 ID（内部从 Attrs["id"] 读取）。
	GetId() string
	// SetId 设置会话绑定的玩家/节点 ID。
	SetId(id string)

	// Close 关闭会话（幂等）。
	Close()

	// DieChan 会话关闭广播通道。
	DieChan() <-chan bool
	// AsynTasksChan 异步任务队列（可读可写）。
	AsynTasksChan() chan func()
	// DataReceivedChan 入站消息队列。
	DataReceivedChan() <-chan *protocol.RequestDataFrame
	// Conn 返回底层连接。
	Conn() net.Conn

	// GetProtocolCodec 返回私有协议栈编解码器。
	GetProtocolCodec() protocol.ProtocolAdapter
	// GetMessageCodec 返回消息编解码器。
	GetMessageCodec() codec.MessageCodec

	// ToString 打印会话关键信息，用于日志。
	ToString() string
}

// 编译期断言：*BaseSession 必须实现 Session 接口。
var _ Session = (*BaseSession)(nil)

package dispatch

import (
	"fmt"
	"net"
	"time"

	"github.com/forfun/gforgame/codec"
	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/network/protocol"
	"github.com/forfun/gforgame/network/session"
)

const DefaultSessionIdleTimeout = 10 * time.Minute

type IoDispatch interface {

	// OnSessionCreated session创建时触发
	OnSessionCreated(session session.Session)

	// OnMessageReceived 收到消息时触发
	OnMessageReceived(session session.Session, msg *protocol.RequestDataFrame)

	// OnSessionClosed session关闭时触发
	OnSessionClosed(session session.Session)
}

// SessionLoopOpts 会话串行消息循环配置
type SessionLoopOpts struct {
	IdleCheckInterval time.Duration
	IdleTimeout       time.Duration
	OnIdleTimeout     func(session session.Session, idleDuration time.Duration)
}

// ServeSession 启动会话，内置eventLoop
func ServeSession(session session.Session, ioDispatch IoDispatch, loopOpts *SessionLoopOpts) {
	conn := session.RawConn()
	defer func() {
		logger.Info(fmt.Sprintf("客户端连接关闭 %s %s", conn.RemoteAddr().String(), conn.LocalAddr().String()))
		ioDispatch.OnSessionClosed(session)
		_ = conn.Close()
	}()

	ioDispatch.OnSessionCreated(session)
	go session.Read()
	go session.Write()

	// eventLoop 必跑，阻塞在这里直到会话结束
	eventLoop(session, ioDispatch, loopOpts)
}

// CreateServeSession 创建会话并启动派发循环
func CreateServeSession(conn net.Conn, messageCodec codec.MessageCodec, ioDispatch IoDispatch, payloadMode protocol.PayloadMode) {
	ioSession := session.NewSessionWithProtocol(conn, messageCodec, protocol.ProtocolTypeBinary)
	ioSession.SetPayloadMode(payloadMode)

	loopOpts := &SessionLoopOpts{
		IdleCheckInterval: time.Minute,
		IdleTimeout:       DefaultSessionIdleTimeout,
		OnIdleTimeout:     onSessionIdleTimeout,
	}
	ServeSession(ioSession, ioDispatch, loopOpts)
}

// onSessionIdleTimeout 会话空闲超时回调：记录日志并主动关闭连接
func onSessionIdleTimeout(session session.Session, idleDuration time.Duration) {
	conn := session.RawConn()
	logger.Info(fmt.Sprintf(
		"会话空闲超时，主动关闭 remote=%s local=%s idle=%s timeout=%s",
		conn.RemoteAddr().String(),
		conn.LocalAddr().String(),
		idleDuration.Truncate(time.Second),
		DefaultSessionIdleTimeout,
	))
	session.Close()
}

func eventLoop(session session.Session, ioDispatch IoDispatch, opts *SessionLoopOpts) {
	if opts == nil {
		opts = &SessionLoopOpts{
			IdleTimeout: DefaultSessionIdleTimeout,
		}
	}

	var idleCheckTicker *time.Ticker
	if opts.IdleCheckInterval > 0 && opts.IdleTimeout > 0 {
		idleCheckTicker = time.NewTicker(opts.IdleCheckInterval)
		defer idleCheckTicker.Stop()
	}

	die := session.DieChan()
	dataReceived := session.DataReceivedChan()

	for {
		select {
		case ioFrame := <-dataReceived:
			if ioFrame != nil {
				ioDispatch.OnMessageReceived(session, ioFrame)
			}
		case <-idleCheckTicker.C:
			idleDuration := time.Since(session.LastReadAt())
			if idleDuration < opts.IdleTimeout {
				continue
			}
			if opts.OnIdleTimeout != nil {
				opts.OnIdleTimeout(session, idleDuration)
			}
			session.Close()
			return
		case <-die:
			// 收到关闭信号，排空剩余消息，防止写端阻塞
			for {
				select {
				case <-dataReceived:
					// 丢弃未处理消息
				default:
					return
				}
			}
		}
	}
}

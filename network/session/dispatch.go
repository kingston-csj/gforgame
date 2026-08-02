package session

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/forfun/gforgame/codec"
	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/network/protocol"
)

const DefaultDirectSessionIdleTimeout = 10 * time.Minute

type IoDispatch interface {
	OnSessionCreated(session Session)
	OnMessageReceived(session Session, msg *protocol.RequestDataFrame)
	OnSessionClosed(session Session)
}

type SerialSessionLoopOptions struct {
	IdleCheckInterval time.Duration
	IdleTimeout       time.Duration
	OnIdleTimeout     func(session Session, idleDuration time.Duration)
}

type DispatchSessionLoopOptions struct {
	WorkerCount        int32
	SessionKey         string
	ResolveWorkerIndex func(ioFrame *protocol.RequestDataFrame, fallbackSessionIdx int, workerCount int) int
	SerialOptions      *SerialSessionLoopOptions
}

func ServeSession(session Session, ioDispatch IoDispatch, run func(session Session)) {
	conn := session.Conn()
	defer func() {
		logger.Info(fmt.Sprintf("客户端连接关闭 %s %s", conn.RemoteAddr().String(), conn.LocalAddr().String()))
		ioDispatch.OnSessionClosed(session)
		UnregisterSession(conn)
		_ = conn.Close()
	}()

	RegisterSession(conn, session)
	ioDispatch.OnSessionCreated(session)

	go session.Read()
	go session.Write()

	if run != nil {
		run(session)
	}
}

// ServeSessionConn 服务直连模式会话：创建会话并启动派发循环
func ServeSessionConn(conn net.Conn, messageCodec codec.MessageCodec, ioDispatch IoDispatch, dispatchWorkers int32, payloadMode PayloadMode) {
	ioSession := NewSessionWithProtocol(conn, messageCodec, protocol.ProtocolTypeBinary)
	ioSession.SetPayloadMode(payloadMode)

	ServeSession(ioSession, ioDispatch, func(session Session) {
		RunDispatchSessionLoop(session, ioDispatch, newDirectDispatchOptions(conn, dispatchWorkers))
	})
}

// newDirectDispatchOptions 构造直连模式下的派发循环参数
func newDirectDispatchOptions(conn net.Conn, dispatchWorkers int32) *DispatchSessionLoopOptions {
	return &DispatchSessionLoopOptions{
		WorkerCount:        dispatchWorkers,
		SessionKey:         conn.RemoteAddr().String(),
		ResolveWorkerIndex: ResolveWorkerIndex,
		SerialOptions: &SerialSessionLoopOptions{
			IdleCheckInterval: time.Minute,
			IdleTimeout:       DefaultDirectSessionIdleTimeout,
			OnIdleTimeout:     onDirectSessionIdleTimeout,
		},
	}
}

// onDirectSessionIdleTimeout 直连模式会话空闲超时回调：记录日志并主动关闭连接
func onDirectSessionIdleTimeout(session Session, idleDuration time.Duration) {
	conn := session.Conn()
	logger.Info(fmt.Sprintf(
		"直连模式会话空闲超时，主动关闭 remote=%s local=%s idle=%s timeout=%s",
		conn.RemoteAddr().String(),
		conn.LocalAddr().String(),
		idleDuration.Truncate(time.Second),
		DefaultDirectSessionIdleTimeout,
	))
	RemoveSession(session, true)
}

func RunDispatchSessionLoop(session Session, ioDispatch IoDispatch, opts *DispatchSessionLoopOptions) {
	if opts == nil || opts.WorkerCount <= 1 {
		var serialOptions *SerialSessionLoopOptions
		if opts != nil {
			serialOptions = opts.SerialOptions
		}
		RunSerialSessionLoop(session, ioDispatch, serialOptions)
		return
	}

	workerCount := int(opts.WorkerCount)
	resolveWorkerIndex := opts.ResolveWorkerIndex
	if resolveWorkerIndex == nil {
		resolveWorkerIndex = func(ioFrame *protocol.RequestDataFrame, fallbackSessionIdx int, workerCount int) int {
			return fallbackSessionIdx
		}
	}

	workerQueues := make([]chan *protocol.RequestDataFrame, workerCount)
	var dispatchWg sync.WaitGroup
	dispatchWg.Add(workerCount)
	die := session.DieChan()
	for i := 0; i < workerCount; i++ {
		workerQueues[i] = make(chan *protocol.RequestDataFrame, 256)
		go func(workerIdx int) {
			defer dispatchWg.Done()
			for {
				select {
				case ioFrame := <-workerQueues[workerIdx]:
					if ioFrame != nil {
						ioDispatch.OnMessageReceived(session, ioFrame)
					}
				case <-die:
					return
				}
			}
		}(i)
	}
	sessionWorkerIdx := HashSessionWorkerIndex(opts.SessionKey, workerCount)
	asynTasks := session.AsynTasksChan()
	dataReceived := session.DataReceivedChan()
	for {
		select {
		case task := <-asynTasks:
			if task != nil {
				task()
			}
		case ioFrame := <-dataReceived:
			if ioFrame == nil {
				continue
			}
			targetWorkerIdx := resolveWorkerIndex(ioFrame, sessionWorkerIdx, workerCount)
			select {
			case workerQueues[targetWorkerIdx] <- ioFrame:
			case <-die:
				dispatchWg.Wait()
				return
			}
		case <-die:
			dispatchWg.Wait()
			return
		}
	}
}

func RunSerialSessionLoop(session Session, ioDispatch IoDispatch, opts *SerialSessionLoopOptions) {
	var idleCheckTicker *time.Ticker
	if opts != nil && opts.IdleCheckInterval > 0 && opts.IdleTimeout > 0 {
		idleCheckTicker = time.NewTicker(opts.IdleCheckInterval)
		defer idleCheckTicker.Stop()
	}

	die := session.DieChan()
	asynTasks := session.AsynTasksChan()
	dataReceived := session.DataReceivedChan()
	for {
		select {
		case task := <-asynTasks:
			if task != nil {
				task()
			}
		case ioFrame := <-dataReceived:
			if ioFrame != nil {
				ioDispatch.OnMessageReceived(session, ioFrame)
			}
		case <-idleCheckChannel(idleCheckTicker):
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
			return
		}
	}
}

func idleCheckChannel(ticker *time.Ticker) <-chan time.Time {
	if ticker == nil {
		return nil
	}
	return ticker.C
}

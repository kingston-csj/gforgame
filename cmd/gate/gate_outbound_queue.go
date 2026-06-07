package main

import (
	limitedlist "github.com/forfun/gforgame/common/container/list"
	"github.com/forfun/gforgame/common/logger"
)

type backendOutboundMsg struct {
	serverID int32
	transfer any
	index    int32
}

const (
	outboundQueueSize  = 4096
	outboundPendingMax = 4096
)

func startOutboundDispatcher() {
	outboundQueue = make(chan *backendOutboundMsg, outboundQueueSize)
	outboundNotify = make(chan struct{}, 1)
	outboundStop = make(chan struct{})
	outboundWg.Add(1)
	go func() {
		defer outboundWg.Done()
		pending := limitedlist.NewLimitedList[*backendOutboundMsg](outboundPendingMax)
		flush := func() {
			for pending.Len() > 0 {
				msg, ok := pending.Front()
				if !ok || msg == nil {
					_, _ = pending.PopFront()
					continue
				}
				if err := sendTransferToBackend(msg.serverID, msg.transfer, msg.index); err != nil {
					// 后端不可用时保留队列，等待后续重连后继续发送。
					return
				}
				_, _ = pending.PopFront()
			}
		}
		for {
			select {
			case msg := <-outboundQueue:
				if msg != nil {
					// 堆压过大时仅顶掉最旧的一条。
					if pending.Len() >= outboundPendingMax {
						logger.Info("outbound pending overflow, dropped oldest message")
					}
					pending.Push(msg)
				}
				flush()
			case <-outboundNotify:
				flush()
			case <-outboundStop:
				return
			}
		}
	}()
}

func stopOutboundDispatcher() {
	if outboundStop == nil {
		return
	}
	close(outboundStop)
	outboundWg.Wait()
}

func enqueueTransfer(serverID int32, transfer any, index int32) error {
	msg := &backendOutboundMsg{
		serverID: serverID,
		transfer: transfer,
		index:    index,
	}
	// 队列满时阻塞等待，确保消息不丢。
	outboundQueue <- msg
	return nil
}

func notifyOutboundDispatcher() {
	if outboundNotify == nil {
		return
	}
	select {
	case outboundNotify <- struct{}{}:
	default:
	}
}

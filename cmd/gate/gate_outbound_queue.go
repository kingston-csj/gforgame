package main

import "github.com/forfun/gforgame/internal/protos"

type backendOutboundMsg struct {
	serverID int32
	transfer *protos.TransferGateToLogic
	index    int32
}

func startOutboundDispatcher() {
	outboundQueue = make(chan *backendOutboundMsg, 10240)
	outboundNotify = make(chan struct{}, 1)
	outboundStop = make(chan struct{})
	outboundWg.Add(1)
	go func() {
		defer outboundWg.Done()
		pending := make([]*backendOutboundMsg, 0, 10240)
		flush := func() {
			for len(pending) > 0 {
				msg := pending[0]
				if err := sendTransferToBackend(msg.serverID, msg.transfer, msg.index); err != nil {
					// 后端不可用时保留队列，等待后续重连后继续发送。
					return
				}
				pending = pending[1:]
			}
		}
		for {
			select {
			case msg := <-outboundQueue:
				if msg != nil {
					pending = append(pending, msg)
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

func enqueueTransfer(serverID int32, transfer *protos.TransferGateToLogic, index int32) error {
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

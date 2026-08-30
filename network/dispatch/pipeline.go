package dispatch

import (
	"github.com/forfun/gforgame/network/protocol"
	"github.com/forfun/gforgame/network/session"
)

type MessageHandler interface {
	MessageReceived(session session.Session, msg *protocol.RequestDataFrame) bool
}

type BaseIoDispatch struct {
	Pipeline []MessageHandler
}

func (d *BaseIoDispatch) AddHandler(h MessageHandler) {
	d.Pipeline = append(d.Pipeline, h)
}

func (d *BaseIoDispatch) OnSessionCreated(session session.Session) {

}

func (d *BaseIoDispatch) OnMessageReceived(session session.Session, msg *protocol.RequestDataFrame) {
	for _, d := range d.Pipeline {
		// 只要有一个返回false，则终止执行
		if !d.MessageReceived(session, msg) {
			break
		}
	}
}

func (d *BaseIoDispatch) OnSessionClosed(session session.Session) {

}

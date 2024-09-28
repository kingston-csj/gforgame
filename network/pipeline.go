package network

type IoDispatch interface {

	// OnSessionCreated session创建时触发
	OnSessionCreated(session *Session)

	// OnMessageReceived 收到消息时触发
	OnMessageReceived(session *Session, msg *RequestDataFrame)

	// OnSessionClosed session关闭时触发
	OnSessionClosed(session *Session)
}

type MessageHandler interface {
	MessageReceived(session *Session, msg *RequestDataFrame) bool
}

type BaseIoDispatch struct {
	Pipeline []MessageHandler
}

func (d *BaseIoDispatch) AddHandler(h MessageHandler) {
	d.Pipeline = append(d.Pipeline, h)
}

func (d *BaseIoDispatch) OnSessionCreated(session *Session) {

}

func (d *BaseIoDispatch) OnMessageReceived(session *Session, msg *RequestDataFrame) {
	for _, d := range d.Pipeline {
		// 只要有一个返回false，则终止执行
		if d.MessageReceived(session, msg) {
			break
		}
	}
}

func (d *BaseIoDispatch) OnSessionClosed(session *Session) {

}

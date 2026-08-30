package actor

import (
	"time"
)

const maxAskTimeout = 3 * time.Second

type ActorRef struct {
	path    ActorPath
	mailbox *Mailbox
}

func (ref *ActorRef) Path() ActorPath {
	return ref.path
}

func (ref *ActorRef) PathString() string {
	return ref.path.String()
}

// 给ActorRef新增Task方法，用来提交任务
func (ref *ActorRef) Task(fn func()) error {
	return ref.Tell(taskMsg{fn: fn})
}

func (ref *ActorRef) Tell(msg Message) error {
	return ref.mailbox.Tell(msg)
}

func (ref *ActorRef) AskTimeout(msg Message, timeout time.Duration) (Message, error) {
	if timeout <= 0 {
        // 非法，给一个最小兜底，避免无限阻塞
        timeout = 200 * time.Millisecond
    }
	if timeout > maxAskTimeout {
        timeout = maxAskTimeout
    }
    return ref.mailbox.Ask(msg, timeout)
}

func (ref *ActorRef) Ask(msg Message) (Message, error) {
	return ref.mailbox.Ask(msg, 2*time.Second)
}

// Stop 业务可以通过ref关闭自己
func (ref *ActorRef) Stop() {
	ref.mailbox.Close()
}

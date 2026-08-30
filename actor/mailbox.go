package actor

import (
	"errors"
	"sync"
	"time"
)

type askWrap struct {
	payload Message
	resp    chan Message
}

type taskMsg struct {
	fn func()
}

type Mailbox struct {
	msgChan chan Message
	closed  bool
	mu      sync.Mutex
}

func NewMailbox(cap int) *Mailbox {
	return &Mailbox{
		msgChan: make(chan Message, cap),
	}
}

func (m *Mailbox) Tell(msg Message) error {
	m.mu.Lock()
	closed := m.closed
	if closed {
		m.mu.Unlock()
		return errors.New("mailbox closed")
	}
	m.mu.Unlock()

	select {
	case m.msgChan <- msg:
		return nil
	default:
		return errors.New("mailbox full")
	}
}

func (m *Mailbox) Ask(req Message, timeout time.Duration) (Message, error) {
    m.mu.Lock()
    closed := m.closed
    if closed {
        m.mu.Unlock()
        return nil, errors.New("mailbox closed")
    }
    m.mu.Unlock()

    respCh := make(chan Message, 1)
    wrap := askWrap{
        payload: req,
        resp:    respCh,
    }

    timer := time.NewTimer(timeout)
    defer timer.Stop()

    // 阶段1：投递消息进mailbox，共享同一个timer
    select {
    case m.msgChan <- wrap:
    case <-timer.C:
        return nil, errors.New("enqueue timeout: mailbox full")
    }

    // 阶段2：等待actor处理返回，继续复用这个timer，不再新建
    select {
    case r := <-respCh:
        return r, nil
    case <-timer.C:
        return nil, errors.New("ask timeout waiting response")
    }
}

func (m *Mailbox) Recv() <-chan Message {
	return m.msgChan
}

// Close 关闭邮箱；不会立刻终止，消费完队列剩余消息，range循环退出
func (m *Mailbox) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.closed {
		m.closed = true
		close(m.msgChan)
	}
}

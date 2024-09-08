package eventbus

import (
	"fmt"
	"sync"
)

// EventBus 定义了一个事件总线，用于注册、注销和触发事件
type EventBus struct {
	mu       sync.RWMutex
	eventMap map[string][]chan interface{}
}

// NewEventBus 创建一个新的 EventBus 实例
func NewEventBus() *EventBus {
	return &EventBus{
		eventMap: make(map[string][]chan interface{}),
	}
}

// Subscribe 注册一个事件监听器
func (eb *EventBus) Subscribe(event string, ch chan interface{}) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.eventMap[event] = append(eb.eventMap[event], ch)
}

// Unsubscribe 注销一个事件监听器
func (eb *EventBus) Unsubscribe(event string, ch chan interface{}) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	var newChans []chan interface{}
	for _, c := range eb.eventMap[event] {
		if c != ch {
			newChans = append(newChans, c)
		}
	}
	eb.eventMap[event] = newChans
}

// Publish 触发一个事件
func (eb *EventBus) Publish(event string, data interface{}) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	for _, ch := range eb.eventMap[event] {
		go func(ch chan interface{}, data interface{}) {
			ch <- data
		}(ch, data)
	}
}

func main() {
	eventBus := NewEventBus()

	// 创建一个监听器
	listener := make(chan interface{})
	defer close(listener)

	// 注册监听器
	eventBus.Subscribe("testEvent", listener)

	// 触发事件
	go func() {
		eventBus.Publish("testEvent", "Hello, EventBus!")
	}()

	// 监听事件
	for data := range listener {
		fmt.Println("Received event data:", data)
		break
	}
}

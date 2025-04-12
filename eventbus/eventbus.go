package eventbus

import (
	"fmt"
	"sync"
)

// EventBus 结构体用于管理事件的发布和订阅
type EventBus struct {
	handlers map[string][]func(interface{})
	mu       sync.Mutex
}

// NewEventBus 创建一个新的 EventBus 实例
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]func(interface{})),
	}
}

// Subscribe 订阅指定事件，将处理函数添加到对应的事件处理列表中
func (eb *EventBus) Subscribe(event string, handler func(interface{})) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.handlers[event] = append(eb.handlers[event], handler)
}

// Publish 发布事件，调用所有订阅该事件的处理函数
func (eb *EventBus) Publish(event string, data interface{}) {
	eb.mu.Lock()
	handlers, exists := eb.handlers[event]
	eb.mu.Unlock()

	if exists {
		for _, handler := range handlers {
			handler(data)
		}
	}
}

func main() {
	// 创建一个新的 EventBus 实例
	bus := NewEventBus()

	// 定义一个事件处理函数
	handler := func(data interface{}) {
		fmt.Printf("Received event data: %v\n", data)
	}

	// 订阅事件
	bus.Subscribe("testEvent", handler)

	// 发布事件
	bus.Publish("testEvent", "Hello, EventBus!")

	// 等待事件处理完成
	fmt.Scanln()
}

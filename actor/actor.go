package actor

import (
	"fmt"
	"reflect"
)

// Message 消息标记接口
type Message interface{}

// Actor 业务Actor接口
type Actor interface {
	OnStart()
	OnMessage(msg Message) Message
	OnStop()
}

// BaseActor 为Actor提供路由功能
type BaseActor struct {
	router map[reflect.Type]func(Message) Message
}

// NewBaseActor 创建新的BaseActor
func NewBaseActor() *BaseActor {
	return &BaseActor{
		router: make(map[reflect.Type]func(Message) Message),
	}
}

func (b *BaseActor) OnMessage(raw Message) Message {
	typ := reflect.TypeOf(raw)
	if h, ok := b.router[typ]; ok {
		return h(raw)
	}
	return nil
}

// 底层原始注册接口，保留
func (b *BaseActor) Register(msg Message, h func(Message) Message) {
    typ := reflect.TypeOf(msg)
    if _, exists := b.router[typ]; exists {
        panic("duplicate register msg type: " + typ.String())
    }
    b.router[typ] = h
}

// 上层强类型泛型封装，业务唯一入口
func RegisterCmd[T any](base *BaseActor, handler func(*T) Message) {
    var zero *T
    // 调用底层Register，复用重复注册校验等逻辑
    base.Register(zero, func(raw Message) Message {
        return handler(raw.(*T))
    })
}

// ActorPath 仅保存完整路径字符串，不拆分多段字段
type ActorPath struct {
	str string
}

// NewActorPath 辅助构造路径 ns/kind/id
func NewActorPath(namespace, kind, id string) ActorPath {
	return ActorPath{str: fmt.Sprintf("%s/%s/%s", namespace, kind, id)}
}

// String 返回完整路径字符串
func (p ActorPath) String() string {
	return p.str
}

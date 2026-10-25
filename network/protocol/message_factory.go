package protocol

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"
)

// MessageFactory 保存单个服务器的消息注册表，可在启动前注册消息。
// 零值可用；注册与查询支持并发调用。
type MessageFactory struct {
	mu         sync.RWMutex
	id2Msg     map[int32]reflect.Type
	msg2Id     map[reflect.Type]int32
	msgName2Id map[string]int32
	id2MsgName map[int32]string
}

func NewMessageFactory() *MessageFactory { return &MessageFactory{} }

// RegisterMessage 注册消息
// @param cmd 消息命令
// @param msg 消息结构体指针，必须指向结构体
func (f *MessageFactory) RegisterMessage(cmd int32, msg any) {
	f.mu.Lock()
	defer f.mu.Unlock()
	typeOf := reflect.TypeOf(msg)
	if typeOf == nil || typeOf.Kind() != reflect.Ptr {
		panic("msg must be ptr")
	}
	_, ok := f.id2Msg[cmd]
	if ok {
		panic("cmd duplicated: " + strconv.Itoa(int(cmd)))
	}

	structType := typeOf.Elem()

	// 检查底层类型是否为结构体（避免传入非结构体指针，如 *int）
	if structType.Kind() != reflect.Struct {
		panic("msg must point to a struct (指针必须指向结构体)")
	}
	structName := structType.Name()
	if f.id2Msg == nil {
		f.id2Msg = make(map[int32]reflect.Type)
		f.msg2Id = make(map[reflect.Type]int32)
		f.id2MsgName = make(map[int32]string)
		f.msgName2Id = make(map[string]int32)
	}
	f.id2Msg[cmd] = typeOf
	f.msg2Id[typeOf] = cmd
	f.id2MsgName[cmd] = structName
	f.msgName2Id[structName] = cmd
}

// GetMessageCmd 获取消息命令
// @param msg 消息结构体，或指向结构体的指针
// @return 消息命令
func (f *MessageFactory) GetMessageCmd(msg any) (int32, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	typ := reflect.TypeOf(msg)
	if typ == nil {
		return 0, errors.New("GetMessageCmd not found")
	}
	if value, ok := f.msg2Id[typ]; ok {
		return value, nil
	}
	if typ.Kind() == reflect.Struct {
		ptrType := reflect.PointerTo(typ)
		if value, ok := f.msg2Id[ptrType]; ok {
			return value, nil
		}
	}
	return 0, fmt.Errorf("GetMessageCmd not found: %s", typ.String())
}

func (f *MessageFactory) GetMessageCmdFromType(typ reflect.Type) (int32, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	if typ == nil {
		return 0, errors.New("GetMessageCmdFromType not found: nil")
	}
	value, ok := f.msg2Id[typ]
	if ok {
		return value, nil
	} else {
		return 0, fmt.Errorf("GetMessageCmdFromType not found: %s", typ.String())
	}
}

func (f *MessageFactory) GetMessageType(cmd int32) (reflect.Type, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	value, ok := f.id2Msg[cmd]
	if ok {
		return value, nil
	} else {
		return nil, errors.New("not found")
	}
}

func (f *MessageFactory) GetMsgName2IdMapper() map[string]int32 {
	f.mu.RLock()
	defer f.mu.RUnlock()
	result := make(map[string]int32, len(f.msgName2Id))
	for name, cmd := range f.msgName2Id {
		result[name] = cmd
	}
	return result
}

func (f *MessageFactory) GetMsgName(cmd int32) (string, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	value, ok := f.id2MsgName[cmd]
	if ok {
		return value, nil
	} else {
		return "", errors.New("not found")
	}
}

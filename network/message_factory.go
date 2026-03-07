package network

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var (
	id2Msg map[int32]reflect.Type = make(map[int32]reflect.Type)
	msg2Id map[reflect.Type]int32 = make(map[reflect.Type]int32)
	msgName2Id map[string]int32 = make(map[string]int32)
	id2MsgName map[int32]string = make(map[int32]string)
)

// RegisterMessage 注册消息
// @param cmd 消息命令
// @param msg 消息结构体指针，必须指向结构体
func RegisterMessage(cmd int32, msg any) {
	typeOf := reflect.TypeOf(msg)
	if typeOf.Kind() != reflect.Ptr {
		panic("msg must be ptr")
	}
	_, ok := id2Msg[cmd]
	if ok {
		panic("cmd duplicated: " + strconv.Itoa(int(cmd)))
	}

	structType := typeOf.Elem()

	// 检查底层类型是否为结构体（避免传入非结构体指针，如 *int）
	if structType.Kind() != reflect.Struct {
		panic("msg must point to a struct (指针必须指向结构体)")
	}
	structName := structType.Name()
	id2Msg[cmd] = typeOf
	msg2Id[typeOf] = cmd
	id2MsgName[cmd] = structName
	msgName2Id[structName] = cmd
}

// GetMessageCmd 获取消息命令
// @param msg 消息结构体，或指向结构体的指针
// @return 消息命令
func GetMessageCmd(msg any) (int32, error) {
	typ := reflect.TypeOf(msg)
	if typ == nil {
		return 0, errors.New("GetMessageCmd not found")
	}
	if value, ok := msg2Id[typ]; ok {
		return value, nil
	}
	if typ.Kind() == reflect.Struct {
		ptrType := reflect.PtrTo(typ)
		if value, ok := msg2Id[ptrType]; ok {
			return value, nil
		}
	}
	return 0, errors.New("GetMessageCmd not found")
}

func GetMessageCmdFromType(typ reflect.Type) (int32, error) {
	value, ok := msg2Id[typ]
	if ok {
		return value, nil
	} else {
		return 0, errors.New(fmt.Sprintf("GetMessageCmdFromType not found: %s", typ.String()))
	}
}


func GetMessageType(cmd int32) (reflect.Type, error) {
	value, ok := id2Msg[cmd]
	if ok {
		return value, nil
	} else {
		return nil, errors.New("not found")
	}
}

func GetMsgName2IdMapper() map[string]int32 {
	return msgName2Id
}

func GetMsgName(cmd int32) (string, error) {
	value, ok := id2MsgName[cmd]
	if ok {
		return value, nil
	} else {
		return "", errors.New("not found")
	}
}

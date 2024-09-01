package network

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	id2Msg map[int]reflect.Type

	msg2Id map[reflect.Type]int
)

func RegisterMessage(cmd int, msg any) {

	id2Msg[cmd] = reflect.TypeOf(msg)

	msg2Id[reflect.TypeOf(msg)] = cmd
}

func init() {
	id2Msg = make(map[int]reflect.Type)
	msg2Id = make(map[reflect.Type]int)
}

func GetMessageCmd(msg any) (int, error) {
	value, ok := msg2Id[reflect.TypeOf(msg)]
	if ok {
		return value, nil
	} else {
		return 0, errors.New("not found")
	}
}

func GetMessageCmdFromType(typ reflect.Type) (int, error) {
	fmt.Println("type", typ)
	value, ok := msg2Id[typ]
	if ok {
		return value, nil
	} else {
		return 0, errors.New("not found")
	}
}

func GetMessageType(cmd int) (reflect.Type, error) {
	value, ok := id2Msg[cmd]
	if ok {
		return value, nil
	} else {
		return nil, errors.New("not found")
	}
}

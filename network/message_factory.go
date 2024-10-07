package network

import (
	"errors"
	"reflect"
	"strconv"
)

var (
	id2Msg map[int]reflect.Type = make(map[int]reflect.Type)
	msg2Id map[reflect.Type]int = make(map[reflect.Type]int)
)

func RegisterMessage(cmd int, msg any) {
	typeOf := reflect.TypeOf(msg)
	if typeOf.Kind() != reflect.Ptr {
		panic("msg must be ptr")
	}
	_, ok := id2Msg[cmd]
	if ok {
		panic("cmd duplicated: " + strconv.Itoa(cmd))
	}

	id2Msg[cmd] = typeOf
	msg2Id[typeOf] = cmd
}

func GetMessageCmd(msg any) (int, error) {
	value, ok := msg2Id[reflect.TypeOf(msg)]
	if ok {
		return value, nil
	} else {
		return 0, errors.New("GetMessageCmd not found")
	}
}

func GetMessageCmdFromType(typ reflect.Type) (int, error) {
	value, ok := msg2Id[typ]
	if ok {
		return value, nil
	} else {
		return 0, errors.New("GetMessageCmdFromType not found")
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

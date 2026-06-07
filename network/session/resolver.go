package session

import (
	"errors"
	"reflect"
)

type MessageResolver interface {
	GetMessageCmd(msg any) (int32, error)
	GetMsgName(cmd int32) (string, error)
	GetMessageType(cmd int32) (reflect.Type, error)
}

var messageResolver MessageResolver = noopMessageResolver{}

func SetMessageResolver(resolver MessageResolver) {
	if resolver == nil {
		messageResolver = noopMessageResolver{}
		return
	}
	messageResolver = resolver
}

type noopMessageResolver struct{}

func (noopMessageResolver) GetMessageCmd(msg any) (int32, error) {
	return 0, errors.New("message resolver is not configured")
}

func (noopMessageResolver) GetMsgName(cmd int32) (string, error) {
	return "", errors.New("message resolver is not configured")
}

func (noopMessageResolver) GetMessageType(cmd int32) (reflect.Type, error) {
	return nil, errors.New("message resolver is not configured")
}

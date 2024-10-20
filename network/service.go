package network

import (
	"fmt"
	"reflect"
)

type (
	//Handler represents a message.Message's handler's meta information.
	Handler struct {
		Receiver reflect.Value  // receiver of method
		Method   reflect.Method // method stub
		Type     reflect.Type   // arg type of method
	}

	MessageRoute struct {
		Handlers map[int]*Handler
	}
)

var (
	typeOfBytes   = reflect.TypeOf(([]byte)(nil))
	typeOfSession = reflect.TypeOf(&Session{})
)

func (self *MessageRoute) RegisterMessageHandlers(comp Module) error {
	clazz := reflect.TypeOf(comp)
	for m := 0; m < clazz.NumMethod(); m++ {
		method := clazz.Method(m)
		mt := method.Type
		if self.isHandlerMethod(method) {
			cmd, err := GetMessageCmdFromType(mt.In(2))
			if err != nil {
				return err
			}

			self.Handlers[cmd] = &Handler{Receiver: reflect.ValueOf(comp), Method: method, Type: mt.In(2)}
		}
	}
	return nil
}

// isHandlerMethod decide a method is suitable handler method
func (self *MessageRoute) isHandlerMethod(method reflect.Method) bool {
	mt := method.Type
	// Method must be exported.
	if method.PkgPath != "" {
		return false
	}
	// Method needs three ins: receiver, *Session, []byte or pointer.
	if mt.NumIn() != 3 {
		return false
	}
	// Method needs one outs: error
	if mt.NumOut() != 1 {
		return false
	}
	if t1 := mt.In(1); t1.Kind() != reflect.Ptr || t1 != typeOfSession {
		return false
	}
	if mt.In(2).Kind() != reflect.Ptr && mt.In(2) != typeOfBytes {
		return false
	}
	return true
}

func (self *MessageRoute) GetHandler(cmd int) (*Handler, error) {
	value, ok := self.Handlers[cmd]
	if ok {
		return value, nil
	} else {
		return nil, fmt.Errorf("cmd [%d] handler not found", cmd)
	}
}

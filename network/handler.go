package network

import (
	"fmt"
	"reflect"
	"strings"
)

type (
	//Handler represents a message.Message's handler's meta information.
	Handler struct {
		Receiver     reflect.Value  // receiver of method
		Method       reflect.Method // method stub
		Type         reflect.Type   // arg type of method
		Indindexed   bool
		NeedValidate bool // 是否需要参数校验
	}

	MessageRoute struct {
		Handlers map[int32]*Handler
	}
)

func NewMessageRoute() *MessageRoute {
	return &MessageRoute{Handlers: make(map[int32]*Handler)}
}

var (
	typeOfSession = reflect.TypeOf(&Session{})
)

func (r *MessageRoute) RegisterMessageHandlers(comp Module) error {
	clazz := reflect.TypeOf(comp)
	for m := 0; m < clazz.NumMethod(); m++ {
		method := clazz.Method(m)
		mt := method.Type
		if r.isHandlerMethod(method) {
			containsIndex := false
			cmdFieldIndex := 2
			if method.Type.NumIn() == 4 {
				containsIndex = true
				cmdFieldIndex = 3
			}
			cmd, err := GetMessageCmdFromType(mt.In(cmdFieldIndex))
			if err != nil {
				return err
			}

			needValidate := r.needValidation(method.Name, mt.In(cmdFieldIndex))

			r.Handlers[cmd] = &Handler{
				Receiver:     reflect.ValueOf(comp),
				Method:       method,
				Type:         mt.In(cmdFieldIndex),
				Indindexed:   containsIndex,
				NeedValidate: needValidate,
			}
		}
	}
	return nil
}

func (r *MessageRoute) needValidation(methodName string, msgType reflect.Type) bool {
	if strings.HasPrefix(methodName, "Validatable") {
		return true
	}
	return hasValidateTag(msgType)
}

func hasValidateTag(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return false
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Tag.Get("validate") != "" && field.Tag.Get("validate") != "-" {
			return true
		}
	}
	return false
}

// isHandlerMethod decide a method is suitable handler method
func (r *MessageRoute) isHandlerMethod(method reflect.Method) bool {
	mt := method.Type
	// Method must be exported.
	if method.PkgPath != "" {
		return false
	}
	// Method needs three ins: receiver, *Session, [index], pointer.
	if mt.NumIn() != 3 && mt.NumIn() != 4 {
		return false
	}
	// Method needs one outs: error
	// if mt.NumOut() != 1 {
	// 	return false
	// }
	if t1 := mt.In(1); t1.Kind() != reflect.Ptr || t1 != typeOfSession {
		return false
	}
	if mt.NumIn() == 3 {
		if mt.In(2).Kind() != reflect.Ptr {
			return false
		}
	}
	// 4个参数才有index int32字段
	if mt.NumIn() == 4 {
		// index must be int32
		if mt.In(2).Kind() != reflect.Int32 {
			panic(fmt.Sprintf("method %s is not a handler method, index must be int32", method.Name))
		}
		if mt.In(3).Kind() != reflect.Ptr {
			panic(fmt.Sprintf("method %s is not a handler method, arg must be pointer", method.Name))
		}
	}

	return true
}

func (r *MessageRoute) GetHandler(cmd int32) (*Handler, error) {
	value, ok := r.Handlers[cmd]
	if ok {
		return value, nil
	} else {
		return nil, fmt.Errorf("cmd [%d] handler not found", cmd)
	}
}

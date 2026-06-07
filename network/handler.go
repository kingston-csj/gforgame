package network

import (
	"fmt"
	"reflect"
	"strings"
)

type (
	//Handler represents a message.Message's handler's meta information.
	Handler struct {
		Receiver   reflect.Value  // receiver of method
		Method     reflect.Method // method stub
		Type       reflect.Type   // arg type of method
		Indindexed bool
		HasPlayer  bool
		HasSession bool
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
	typeOfString  = reflect.TypeOf("")
)

func (r *MessageRoute) RegisterMessageHandlers(comp any) error {
	clazz := reflect.TypeOf(comp)
	for m := 0; m < clazz.NumMethod(); m++ {
		method := clazz.Method(m)
		mt := method.Type
		if r.isHandlerMethod(method) {
			hasPlayer, hasSession, containsIndex, cmdFieldIndex := parseHandlerSignature(mt)
			cmd, err := GetMessageCmdFromType(mt.In(cmdFieldIndex))
			if err != nil {
				return err
			}
			needValidate := r.needValidation(method.Name, mt.In(cmdFieldIndex))
			r.Handlers[cmd] = &Handler{
				Receiver:   reflect.ValueOf(comp),
				Method:     method,
				Type:       mt.In(cmdFieldIndex),
				Indindexed: containsIndex,
				HasPlayer:  hasPlayer,
				HasSession: hasSession,
				NeedValidate: needValidate,
			}
		}  
	}
	return nil
}

// isHandlerMethod decide a method is suitable handler method
func (r *MessageRoute) isHandlerMethod(method reflect.Method) bool {
	mt := method.Type
	// Method must be exported.
	if method.PkgPath != "" {
		return false
	}
	// 兼容签名（receiver后参数）：
	// 1) *Session, *Req
	// 2) *Session, index, *Req
	// 3) playerId, *Req
	// 4) playerId, index, *Req
	// 5) playerId, *Session, *Req
	// 6) playerId, *Session, index, *Req
	if mt.NumIn() != 3 && mt.NumIn() != 4 && mt.NumIn() != 5 {
		return false
	}
	// Method needs one outs: error
	// if mt.NumOut() != 1 {
	// 	return false
	// }
	_, _, _, reqIndex := parseHandlerSignature(mt)
	return reqIndex > 0
}

func parseHandlerSignature(mt reflect.Type) (hasPlayer bool, hasSession bool, hasIndex bool, reqIndex int) {
	if mt.NumIn() < 3 || mt.NumIn() > 5 {
		return false, false, false, -1
	}
	i := 1 // skip receiver
	last := mt.NumIn() - 1

	if mt.In(i) == typeOfString {
		hasPlayer = true
		i++
		if i > last {
			return false, false, false, -1
		}
	}
	if mt.In(i) == typeOfSession {
		hasSession = true
		i++
		if i > last {
			return false, false, false, -1
		}
	}
	if i < last {
		if mt.In(i).Kind() != reflect.Int32 {
			return false, false, false, -1
		}
		hasIndex = true
		i++
	}
	if i != last {
		return false, false, false, -1
	}
	if mt.In(i).Kind() != reflect.Ptr {
		panic(fmt.Sprintf("method %s is not a handler method, arg must be pointer", mt.String()))
	}
	if !hasPlayer && !hasSession {
		return false, false, false, -1
	}
	return hasPlayer, hasSession, hasIndex, i
}

func (r *MessageRoute) GetHandler(cmd int32) (*Handler, error) {
	value, ok := r.Handlers[cmd]
	if ok {
		return value, nil
	} else {
		return nil, fmt.Errorf("cmd [%d] handler not found", cmd)
	}
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

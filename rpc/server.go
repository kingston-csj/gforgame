package rpc

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/protocol"
)

// ErrorResponse must be registered with MessageFactory.RegisterMessage on both peers.
// The selected codec must support this struct (for example codec/json).
type ErrorResponse struct {
	RequestCmd int32  `json:"requestCmd"`
	Message    string `json:"message"`
}

func (e *ErrorResponse) Error() string { return "rpc: remote error: " + e.Message }

type method struct {
	receiver    reflect.Value
	function    reflect.Method
	request     reflect.Type
	withSession bool
}

// Server handles requests in the calling dispatch goroutine. Long running handlers
// block that session's business dispatch. Register services before serving traffic.
type Server struct {
	factory *protocol.MessageFactory
	mu      sync.RWMutex
	methods map[int32]method
}

func NewServer(factory *protocol.MessageFactory) *Server {
	return &Server{factory: factory, methods: make(map[int32]method)}
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()
var sessionType = reflect.TypeOf((*network.Session)(nil)).Elem()

// CommandMessage supplies the application-assigned command for automatic message
// registration. RPCCommand must return a stable value without side effects.
// Existing explicit factory registrations remain supported.
type CommandMessage interface {
	RPCCommand() int32
}

// messageCommand validates and stages missing registrations. Registration is
// committed only after the entire service (or client response pair) is valid.
func messageCommand(factory *protocol.MessageFactory, typ reflect.Type, staged map[int32]reflect.Type) (int32, error) {
	registered, err := factory.GetMessageCmdFromType(typ)
	message, automatic := reflect.New(typ.Elem()).Interface().(CommandMessage)
	if !automatic {
		return registered, err
	}
	cmd := message.RPCCommand()
	if err == nil {
		if registered != cmd {
			return 0, fmt.Errorf("rpc: %s declares cmd %d but is registered as %d", typ, cmd, registered)
		}
		return cmd, nil
	}
	if other, err := factory.GetMessageType(cmd); err == nil && other != typ {
		return 0, fmt.Errorf("rpc: cmd %d is already registered for %s", cmd, other)
	}
	if other, exists := staged[cmd]; exists && other != typ {
		return 0, fmt.Errorf("rpc: cmd %d is shared by %s and %s", cmd, other, typ)
	}
	staged[cmd] = typ
	return cmd, nil
}

func registerMessages(factory *protocol.MessageFactory, staged map[int32]reflect.Type) {
	for cmd, typ := range staged {
		factory.RegisterMessage(cmd, reflect.New(typ.Elem()).Interface())
	}
}

func messagePair(factory *protocol.MessageFactory, req, res reflect.Type, staged map[int32]reflect.Type) (int32, int32, error) {
	for _, t := range []reflect.Type{req, res} {
		if t == nil || t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
			return 0, 0, errors.New("rpc: request and response must be struct pointers")
		}
		if t == reflect.TypeOf(&ErrorResponse{}) {
			return 0, 0, errors.New("rpc: ErrorResponse is reserved")
		}
	}
	q, err := messageCommand(factory, req, staged)
	if err != nil {
		return 0, 0, err
	}
	r, err := messageCommand(factory, res, staged)
	if err != nil {
		return 0, 0, err
	}
	if q == r {
		return 0, 0, errors.New("rpc: request and response must have different commands")
	}
	return q, r, nil
}

// Register discovers RPC_Method(*Req) *Res and RPC_Method(*Req) (*Res, error).
// Methods without the RPC_ prefix are ignored. Invalid RPC_ methods return an
// error. Failed registrations install no methods.
// An optional network.Session argument may precede *Req.
// Request and response types implementing CommandMessage are registered
// automatically; other types must already be registered in the factory.
func (s *Server) Register(service any) error {
	if isNil(service) {
		return errors.New("rpc: nil service")
	}
	if _, err := s.factory.GetMessageCmd(&ErrorResponse{}); err != nil {
		return fmt.Errorf("rpc: register ErrorResponse first: %w", err)
	}
	t := reflect.TypeOf(service)
	methods := make(map[int32]method)
	staged := make(map[int32]reflect.Type)
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if !strings.HasPrefix(m.Name, "RPC_") {
			continue
		}
		mt := m.Type
		withSession := mt.NumIn() == 3 && mt.In(1) == sessionType
		if (mt.NumIn() != 2 && !withSession) || (mt.NumOut() != 1 && mt.NumOut() != 2) {
			return fmt.Errorf("rpc: %s: expected RPC_Method([network.Session,] *Req) *Res or (*Res, error)", m.Name)
		}
		if mt.NumOut() == 2 && mt.Out(1) != errorType {
			return fmt.Errorf("rpc: %s: second return value must be error", m.Name)
		}
		requestType := mt.In(mt.NumIn() - 1)
		cmd, _, err := messagePair(s.factory, requestType, mt.Out(0), staged)
		if err != nil {
			return fmt.Errorf("rpc: %s: %w", m.Name, err)
		}
		if _, exists := methods[cmd]; exists {
			return fmt.Errorf("rpc: duplicate request cmd %d", cmd)
		}
		methods[cmd] = method{receiver: reflect.ValueOf(service), function: m, request: requestType, withSession: withSession}
	}
	if len(methods) == 0 {
		return errors.New("rpc: no RPC_ service methods")
	}
	for cmd := range methods {
		if _, exists := s.methods[cmd]; exists {
			return fmt.Errorf("rpc: duplicate request cmd %d", cmd)
		}
	}
	registerMessages(s.factory, staged)
	for cmd, m := range methods {
		s.methods[cmd] = m
	}
	return nil
}

// MessageReceived leaves unknown commands to later pipeline handlers.
func (s *Server) MessageReceived(session network.Session, frame *protocol.RequestDataFrame) bool {
	if frame == nil || isNil(session) {
		return true
	}
	s.mu.RLock()
	m, ok := s.methods[frame.Header.Cmd]
	s.mu.RUnlock()
	if !ok {
		return true
	}
	response, err := m.invoke(session, frame.Msg)
	if err != nil {
		response = &ErrorResponse{RequestCmd: frame.Header.Cmd, Message: err.Error()}
	}
	if err := session.Send(response, frame.Header.Index); err != nil {
		slog.Error("rpc response send failed", "cmd", frame.Header.Cmd, "index", frame.Header.Index, "error", err)
		session.Close()
	}
	return false
}

func (m method) invoke(session network.Session, request any) (response any, err error) {
	defer func() {
		if p := recover(); p != nil {
			slog.Error("rpc handler panic", "method", m.function.Name, "panic", p, "stack", string(debug.Stack()))
			response, err = nil, errors.New("internal server error")
		}
	}()
	if isNil(request) || reflect.TypeOf(request) != m.request {
		return nil, errors.New("invalid request type")
	}
	args := []reflect.Value{m.receiver}
	if m.withSession {
		args = append(args, reflect.ValueOf(session))
	}
	args = append(args, reflect.ValueOf(request))
	values := m.function.Func.Call(args)
	if len(values) == 2 && !values[1].IsNil() {
		return nil, values[1].Interface().(error)
	}
	if values[0].IsNil() {
		return nil, errors.New("handler returned nil response")
	}
	return values[0].Interface(), nil
}

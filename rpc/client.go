package rpc

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/protocol"
)

type CallbackTimeoutException struct{}

func (*CallbackTimeoutException) Error() string { return "rpc: request timeout" }

var (
	ErrTimeout       = &CallbackTimeoutException{}
	ErrSessionClosed = errors.New("rpc: session closed")
	ErrClientClosed  = errors.New("rpc: client closed")
)

type RequestCallback interface {
	OnSuccess(response any)
	OnFail(err error)
}

type RpcClient interface {
	Callback(session network.Session, request any, callback RequestCallback) error
	Request(session network.Session, request any) (any, error)
}

type result struct {
	response any
	err      error
}
type pending struct {
	die                     <-chan bool
	requestCmd, responseCmd int32
	result                  chan result
}

// Client implements RpcClient and dispatch.MessageHandler. Add it before business
// handlers. Configure protocol message registrations before serving traffic.
type Client struct {
	factory      *protocol.MessageFactory
	mu           sync.Mutex
	timeout      time.Duration
	next         int32
	closed       bool
	responses    map[int32]int32
	responseCmds map[int32]bool
	pending      map[int32]*pending
}

var _ RpcClient = (*Client)(nil)

func NewClient(timeout time.Duration, factory *protocol.MessageFactory) *Client {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &Client{factory: factory, timeout: timeout, responses: make(map[int32]int32), responseCmds: make(map[int32]bool), pending: make(map[int32]*pending)}
}

// GetMessageFactory returns the RPC registry for configuring transport sessions.
// Complete all registrations before opening sessions.
func (c *Client) GetMessageFactory() *protocol.MessageFactory { return c.factory }

// RegisterResponse associates request and response types, automatically
// registering missing types that implement CommandMessage.
func (c *Client) RegisterResponse(request, response any) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	staged := make(map[int32]reflect.Type)
	req, res, err := messagePair(c.factory, reflect.TypeOf(request), reflect.TypeOf(response), staged)
	if err != nil {
		return err
	}
	if _, err := c.factory.GetMessageCmd(&ErrorResponse{}); err != nil {
		return fmt.Errorf("rpc: register ErrorResponse first: %w", err)
	}
	if _, exists := c.responses[req]; exists {
		return fmt.Errorf("rpc: duplicate request cmd %d", req)
	}
	if c.responseCmds[req] {
		return fmt.Errorf("rpc: cmd %d is already a response", req)
	}
	if _, exists := c.responses[res]; exists {
		return fmt.Errorf("rpc: cmd %d is already a request", res)
	}
	registerMessages(c.factory, staged)
	c.responses[req] = res
	c.responseCmds[res] = true
	return nil
}

func (c *Client) begin(session network.Session, request any) (int32, *pending, error) {
	if isNil(session) || isNil(request) {
		return 0, nil, errors.New("rpc: nil session or request")
	}
	cmd, err := c.factory.GetMessageCmd(request)
	if err != nil {
		return 0, nil, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return 0, nil, ErrClientClosed
	}
	res, ok := c.responses[cmd]
	if !ok {
		return 0, nil, fmt.Errorf("rpc: no response registered for cmd %d", cmd)
	}
	select {
	case <-session.DieChan():
		return 0, nil, ErrSessionClosed
	default:
	}
	for {
		if c.next == 1<<31-1 {
			c.next = 0
		}
		c.next++
		if c.pending[c.next] == nil {
			break
		}
	}
	p := &pending{die: session.DieChan(), requestCmd: cmd, responseCmd: res, result: make(chan result, 1)}
	c.pending[c.next] = p
	return c.next, p, nil
}

// complete removes the request before delivering its single result.
func (c *Client) complete(index int32, p *pending, r result) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.pending[index] == p {
		delete(c.pending, index)
		p.result <- r
	}
}

func (c *Client) await(index int32, p *pending, timer *time.Timer) result {
	defer timer.Stop()
	select {
	case r := <-p.result:
		return r
	case <-timer.C:
		c.complete(index, p, result{err: ErrTimeout})
	case <-p.die:
		c.complete(index, p, result{err: ErrSessionClosed})
	}
	return <-p.result
}

func (c *Client) Request(session network.Session, request any) (any, error) {
	index, p, err := c.begin(session, request)
	if err != nil {
		return nil, err
	}
	timer := time.NewTimer(c.timeout)
	if err := session.Send(request, index); err != nil {
		c.complete(index, p, result{err: err})
	}
	r := c.await(index, p, timer)
	return r.response, r.err
}

// Callback returns immediate validation/send errors without invoking callback.
// After a nil return, callback runs exactly once in its own goroutine.
func (c *Client) Callback(session network.Session, request any, callback RequestCallback) error {
	if isNil(callback) {
		return errors.New("rpc: nil callback")
	}
	index, p, err := c.begin(session, request)
	if err != nil {
		return err
	}
	timer := time.NewTimer(c.timeout)
	if err := session.Send(request, index); err != nil {
		timer.Stop()
		c.complete(index, p, result{err: err})
		return err
	}
	go func() {
		r := c.await(index, p, timer)
		if r.err != nil {
			callback.OnFail(r.err)
		} else {
			callback.OnSuccess(r.response)
		}
	}()
	return nil
}

// MessageReceived returns false for responses, including late/duplicate ones.
func (c *Client) MessageReceived(session network.Session, frame *protocol.RequestDataFrame) bool {
	if frame == nil || isNil(session) {
		return true
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	errorCmd, errorRegistered := c.factory.GetMessageCmd(&ErrorResponse{})
	isError := errorRegistered == nil && frame.Header.Cmd == errorCmd
	if !isError && !c.responseCmds[frame.Header.Cmd] {
		return true
	}
	p := c.pending[frame.Header.Index]
	if p == nil || p.die != session.DieChan() {
		return false
	}
	r := result{response: frame.Msg}
	if isError {
		e, ok := frame.Msg.(*ErrorResponse)
		if !ok || e == nil || e.RequestCmd != p.requestCmd {
			return false
		}
		r = result{err: e}
	} else {
		if frame.Header.Cmd != p.responseCmd || isNil(frame.Msg) {
			return false
		}
		cmd, err := c.factory.GetMessageCmd(frame.Msg)
		if err != nil || cmd != p.responseCmd {
			return false
		}
	}
	delete(c.pending, frame.Header.Index)
	p.result <- r
	return false
}

// Close fails outstanding calls without closing their network sessions.
func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
	for index, p := range c.pending {
		delete(c.pending, index)
		p.result <- result{err: ErrClientClosed}
	}
}

func isNil(v any) bool {
	if v == nil {
		return true
	}
	r := reflect.ValueOf(v)
	switch r.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan:
		return r.IsNil()
	}
	return false
}

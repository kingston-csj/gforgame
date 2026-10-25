package rpc

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	jsoncodec "github.com/forfun/gforgame/codec/json"
	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/protocol"
	"github.com/forfun/gforgame/network/session"
)

type testRequest struct {
	Value int
	Mode  string
}
type testResponse struct{ Value int }
type simpleRequest struct{ Value int }
type simpleResponse struct{ Value int }
type testService struct{}

func (*testService) RPC_Handle(req *testRequest) (*testResponse, error) {
	switch req.Mode {
	case "error":
		return nil, errors.New("invalid value")
	case "panic":
		panic("private details")
	case "nil":
		return nil, nil
	}
	return &testResponse{Value: req.Value * 2}, nil
}
func (*testService) RPC_Simple(req *simpleRequest) *simpleResponse {
	return &simpleResponse{Value: req.Value}
}

func testMessageFactory() *protocol.MessageFactory {
	factory := protocol.NewMessageFactory()
	factory.RegisterMessage(1700000001, &testRequest{})
	factory.RegisterMessage(1700000002, &testResponse{})
	factory.RegisterMessage(1700000003, &simpleRequest{})
	factory.RegisterMessage(1700000004, &simpleResponse{})
	factory.RegisterMessage(1700000005, &ErrorResponse{})
	return factory
}

func testClient(t *testing.T, timeout time.Duration) *Client {
	t.Helper()
	c := NewClient(timeout, testMessageFactory())
	t.Cleanup(c.Close)
	if err := c.RegisterResponse(&testRequest{}, &testResponse{}); err != nil {
		t.Fatal(err)
	}
	if err := c.RegisterResponse(&simpleRequest{}, &simpleResponse{}); err != nil {
		t.Fatal(err)
	}
	return c
}

type fakeSession struct {
	network.Session
	die  chan bool
	once sync.Once
	send func(any, int32) error
}

func newFakeSession() *fakeSession {
	return &fakeSession{die: make(chan bool), send: func(any, int32) error { return nil }}
}
func (s *fakeSession) DieChan() <-chan bool      { return s.die }
func (s *fakeSession) Send(v any, i int32) error { return s.send(v, i) }
func (s *fakeSession) Close()                    { s.once.Do(func() { close(s.die) }) }

type testCallback struct{ results chan result }

func (c *testCallback) OnSuccess(v any) { c.results <- result{response: v} }
func (c *testCallback) OnFail(e error)  { c.results <- result{err: e} }
func receive(t *testing.T, ch <-chan result) result {
	t.Helper()
	select {
	case r := <-ch:
		return r
	case <-time.After(3 * time.Second):
		t.Fatal("call did not finish")
		return result{}
	}
}
func assertEmpty(t *testing.T, c *Client) {
	t.Helper()
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.pending) != 0 {
		t.Fatalf("pending leaked: %d", len(c.pending))
	}
}

func TestSessionRoundTrip(t *testing.T) {
	c := testClient(t, 3*time.Second)
	server := NewServer(testMessageFactory())
	if err := server.Register(&testService{}); err != nil {
		t.Fatal(err)
	}
	a, b := net.Pipe()
	left := session.NewSession(a, jsoncodec.NewSerializer(), testMessageFactory())
	right := session.NewSession(b, jsoncodec.NewSerializer(), testMessageFactory())
	var loops sync.WaitGroup
	start := func(f func()) { loops.Add(1); go func() { defer loops.Done(); f() }() }
	start(left.Read)
	start(left.Write)
	start(right.Read)
	start(right.Write)
	start(func() {
		for {
			select {
			case frame := <-left.DataReceivedChan():
				c.MessageReceived(left, frame)
			case <-left.DieChan():
				return
			}
		}
	})
	start(func() {
		for {
			select {
			case frame := <-right.DataReceivedChan():
				server.MessageReceived(right, frame)
			case <-right.DieChan():
				return
			}
		}
	})
	t.Cleanup(func() { left.Close(); right.Close(); loops.Wait() })
	var calls sync.WaitGroup
	for i := 0; i < 40; i++ {
		calls.Add(1)
		go func(n int) {
			defer calls.Done()
			v, err := c.Request(left, &testRequest{Value: n})
			if err != nil {
				t.Error(err)
				return
			}
			if v.(*testResponse).Value != n*2 {
				t.Errorf("response mismatch: %v", v)
			}
		}(i)
	}
	calls.Wait()
	cb := &testCallback{results: make(chan result, 2)}
	if err := c.Callback(left, &simpleRequest{Value: 7}, cb); err != nil {
		t.Fatal(err)
	}
	if r := receive(t, cb.results); r.err != nil || r.response.(*simpleResponse).Value != 7 {
		t.Fatalf("callback: %+v", r)
	}
	for _, mode := range []string{"error", "panic", "nil"} {
		_, err := c.Request(left, &testRequest{Mode: mode})
		var remote *ErrorResponse
		if !errors.As(err, &remote) {
			t.Fatalf("%s: expected remote error, got %v", mode, err)
		}
		if mode == "panic" && remote.Message != "internal server error" {
			t.Fatal("panic details leaked")
		}
	}
	assertEmpty(t, c)
}

func TestImmediateResponseAndFiltering(t *testing.T) {
	c := testClient(t, time.Second)
	s, other := newFakeSession(), newFakeSession()
	s.send = func(_ any, index int32) error {
		frame := &protocol.RequestDataFrame{Header: protocol.MessageHeader{Cmd: 1700000002, Index: index}, Msg: &testResponse{Value: 9}}
		c.MessageReceived(other, frame)
		c.mu.Lock()
		remaining := len(c.pending)
		c.mu.Unlock()
		if remaining != 1 {
			t.Error("wrong session completed call")
		}
		wrong := &protocol.RequestDataFrame{Header: protocol.MessageHeader{Cmd: 1700000004, Index: index}, Msg: &simpleResponse{}}
		c.MessageReceived(s, wrong)
		c.MessageReceived(s, frame)
		if c.MessageReceived(s, frame) {
			t.Error("duplicate response passed to business handler")
		}
		return nil
	}
	v, err := c.Request(s, &testRequest{})
	if err != nil || v.(*testResponse).Value != 9 {
		t.Fatalf("result %v, %v", v, err)
	}
	if !c.MessageReceived(s, &protocol.RequestDataFrame{Header: protocol.MessageHeader{Cmd: 1700000001}}) {
		t.Fatal("request intercepted as response")
	}
	assertEmpty(t, c)
}

func TestCallFailures(t *testing.T) {
	for _, mode := range []string{"timeout", "disconnect", "close", "send"} {
		t.Run(mode, func(t *testing.T) {
			c := testClient(t, 20*time.Millisecond)
			s := newFakeSession()
			sent := make(chan struct{})
			sendErr := errors.New("send failed")
			s.send = func(any, int32) error {
				close(sent)
				if mode == "send" {
					return sendErr
				}
				return nil
			}
			done := make(chan result, 1)
			go func() { v, e := c.Request(s, &testRequest{}); done <- result{v, e} }()
			<-sent
			want := error(ErrTimeout)
			switch mode {
			case "disconnect":
				s.Close()
				want = ErrSessionClosed
			case "close":
				c.Close()
				want = ErrClientClosed
			case "send":
				want = sendErr
			}
			if r := receive(t, done); !errors.Is(r.err, want) {
				t.Fatalf("want %v, got %+v", want, r)
			}
			assertEmpty(t, c)
			if !c.closed {
				if c.MessageReceived(s, &protocol.RequestDataFrame{Header: protocol.MessageHeader{Cmd: 1700000002, Index: 1}, Msg: &testResponse{}}) {
					t.Fatal("late response was forwarded")
				}
			}
		})
	}
}

func TestAsyncFailureAndExactlyOnce(t *testing.T) {
	c := testClient(t, 10*time.Millisecond)
	s := newFakeSession()
	cb := &testCallback{results: make(chan result, 4)}
	var index int32
	s.send = func(_ any, i int32) error { index = i; return nil }
	if err := c.Callback(s, &testRequest{}, cb); err != nil {
		t.Fatal(err)
	}
	if r := receive(t, cb.results); !errors.Is(r.err, ErrTimeout) {
		t.Fatal(r.err)
	}
	c.MessageReceived(s, &protocol.RequestDataFrame{Header: protocol.MessageHeader{Cmd: 1700000002, Index: index}, Msg: &testResponse{}})
	assertEmpty(t, c)
	select {
	case r := <-cb.results:
		t.Fatalf("duplicate callback: %+v", r)
	default:
	}
	s.send = func(any, int32) error { return errors.New("send failed") }
	if err := c.Callback(s, &testRequest{}, cb); err == nil {
		t.Fatal("expected immediate error")
	}
	select {
	case <-cb.results:
		t.Fatal("callback after immediate error")
	default:
	}
	assertEmpty(t, c)
}

func TestConcurrentCompletion(t *testing.T) {
	c := testClient(t, time.Second)
	for i := 0; i < 100; i++ {
		s := newFakeSession()
		index, p, err := c.begin(s, &testRequest{})
		if err != nil {
			t.Fatal(err)
		}
		var wg sync.WaitGroup
		for j := 0; j < 3; j++ {
			wg.Add(1)
			go func(n int) { defer wg.Done(); c.complete(index, p, result{err: fmt.Errorf("winner %d", n)}) }(j)
		}
		wg.Wait()
		receive(t, p.result)
		select {
		case <-p.result:
			t.Fatal("more than one result")
		default:
		}
	}
	assertEmpty(t, c)
}

func TestRegistrationAndValidation(t *testing.T) {
	c := testClient(t, time.Second)
	if err := c.RegisterResponse(&testRequest{}, &testResponse{}); err == nil {
		t.Fatal("duplicate accepted")
	}
	if err := c.RegisterResponse(&testResponse{}, &testRequest{}); err == nil {
		t.Fatal("ambiguous response accepted")
	}
	if _, err := c.Request(nil, &testRequest{}); err == nil {
		t.Fatal("nil session accepted")
	}
	s := newFakeSession()
	if err := c.Callback(s, &testRequest{}, nil); err == nil {
		t.Fatal("nil callback accepted")
	}
	if _, err := c.Request(s, (*testRequest)(nil)); err == nil {
		t.Fatal("typed nil accepted")
	}
	s.Close()
	if _, err := c.Request(s, &testRequest{}); !errors.Is(err, ErrSessionClosed) {
		t.Fatal(err)
	}
	server := NewServer(testMessageFactory())
	if err := server.Register(&testService{}); err != nil {
		t.Fatal(err)
	}
	if err := server.Register(&testService{}); err == nil {
		t.Fatal("duplicate service accepted")
	}
	if err := server.Register((*testService)(nil)); err == nil {
		t.Fatal("nil service accepted")
	}
	c.Close()
	if _, err := c.Request(newFakeSession(), &testRequest{}); !errors.Is(err, ErrClientClosed) {
		t.Fatal(err)
	}
}

// Ordinary exported methods must not be discovered, even with RPC-like signatures.
func (*testService) Handle(*testRequest) *testResponse { panic("ordinary method must not run") }
func (*testService) Helper()                           {}

type ordinaryService struct{}

func (*ordinaryService) Handle(*testRequest) *testResponse { return &testResponse{} }

type invalidRPCArgs struct{}

func (*invalidRPCArgs) RPC_Handle() *testResponse { return &testResponse{} }

type invalidRPCError struct{}

func (*invalidRPCError) RPC_Handle(*testRequest) (*testResponse, bool) { return nil, false }

type invalidRPCRequest struct{}

func (*invalidRPCRequest) RPC_Handle(testRequest) *testResponse { return &testResponse{} }

type invalidRPCResponse struct{}

func (*invalidRPCResponse) RPC_Handle(*testRequest) testResponse { return testResponse{} }

type partialRPCService struct{}

func (*partialRPCService) RPC_AValid(*testRequest) *testResponse { return &testResponse{} }
func (*partialRPCService) RPC_ZInvalid()                         {}

func TestRPCMethodPrefix(t *testing.T) {
	server := NewServer(testMessageFactory())
	if err := server.Register(&testService{}); err != nil {
		t.Fatal(err)
	}
	if len(server.methods) != 2 {
		t.Fatalf("expected only two RPC_ methods, got %d", len(server.methods))
	}
	for _, tc := range []struct {
		name    string
		service any
	}{
		{"ordinary", &ordinaryService{}},
		{"arguments", &invalidRPCArgs{}},
		{"error return", &invalidRPCError{}},
		{"request type", &invalidRPCRequest{}},
		{"response type", &invalidRPCResponse{}},
		{"atomic registration", &partialRPCService{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := NewServer(testMessageFactory())
			if err := s.Register(tc.service); err == nil {
				t.Fatal("expected registration error")
			}
			if len(s.methods) != 0 {
				t.Fatal("failed registration installed methods")
			}
		})
	}
}

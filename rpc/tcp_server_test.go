package rpc

import (
	"testing"

	jsoncodec "github.com/forfun/gforgame/codec/json"
	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/protocol"
)

func TestTCPServerDispatchAndLifecycle(t *testing.T) {
	created, closed := 0, 0
	server := NewTCPServer("", testMessageFactory(), jsoncodec.NewSerializer(),
		WithSessionHooks(func(network.Session) { created++ }, func(network.Session) { closed++ }))
	if err := server.Register(&testService{}); err != nil {
		t.Fatal(err)
	}
	session := newFakeSession()
	server.dispatcher.OnSessionCreated(session)
	var response any
	var responseIndex int32
	session.send = func(msg any, index int32) error {
		response, responseIndex = msg, index
		return nil
	}
	server.dispatcher.OnMessageReceived(session, &protocol.RequestDataFrame{
		Header: protocol.MessageHeader{Cmd: 1700000001, Index: 42},
		Msg:    &testRequest{Value: 7},
	})
	if result, ok := response.(*testResponse); !ok || result.Value != 14 || responseIndex != 42 {
		t.Fatalf("RPC service was not dispatched: response=%v index=%d", response, responseIndex)
	}
	server.Stop()
	server.Stop()
	select {
	case <-session.DieChan():
	default:
		t.Fatal("Stop left an RPC session open")
	}
	server.dispatcher.OnSessionClosed(session)
	if created != 1 || closed != 1 {
		t.Fatalf("unexpected lifecycle hooks: created=%d closed=%d", created, closed)
	}
	late := newFakeSession()
	server.dispatcher.OnSessionCreated(late)
	server.dispatcher.OnSessionClosed(late)
	select {
	case <-late.DieChan():
	default:
		t.Fatal("session accepted during shutdown remained open")
	}
	if created != 1 || closed != 1 {
		t.Fatal("rejected session should not trigger application hooks")
	}
	if err := server.Start(); err == nil {
		t.Fatal("a stopped server must not restart")
	}
}

func TestTCPServerStartFailure(t *testing.T) {
	server := NewTCPServer("invalid-address", testMessageFactory(), jsoncodec.NewSerializer())
	defer server.Stop()
	if err := server.Start(); err == nil {
		t.Fatal("expected listener startup error")
	}
}

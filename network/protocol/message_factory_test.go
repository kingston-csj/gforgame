package protocol

import (
	"reflect"
	"testing"
)

type factoryRequest struct{}
type factoryResponse struct{}

func TestMessageFactoryIsolation(t *testing.T) {
	var first MessageFactory
	second := NewMessageFactory()
	first.RegisterMessage(1, &factoryRequest{})
	second.RegisterMessage(1, &factoryResponse{})
	second.RegisterMessage(2, &factoryRequest{})
	for _, input := range []any{factoryRequest{}, &factoryRequest{}} {
		if cmd, err := first.GetMessageCmd(input); err != nil || cmd != 1 {
			t.Fatalf("first cmd = %d, %v", cmd, err)
		}
		if cmd, err := second.GetMessageCmd(input); err != nil || cmd != 2 {
			t.Fatalf("second cmd = %d, %v", cmd, err)
		}
	}
	if typ, err := second.GetMessageType(1); err != nil || typ != reflect.TypeOf(&factoryResponse{}) {
		t.Fatalf("second type = %v, %v", typ, err)
	}
	snapshot := first.GetMsgName2IdMapper()
	snapshot["factoryRequest"] = 99
	if first.GetMsgName2IdMapper()["factoryRequest"] != 1 {
		t.Fatal("caller mutated registry via snapshot")
	}
}

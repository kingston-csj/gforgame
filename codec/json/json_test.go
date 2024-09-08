package json

import (
	"reflect"
	"testing"
)

type Message struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

func TestJsonCodec(t *testing.T) {
	m := Message{1, "hello world"}
	s := &JsonCodec{}
	b, err := s.Encode(m)
	if err != nil {
		t.Fail()
	}

	m2 := Message{}
	if err := s.Decode(b, &m2); err != nil {
		t.Fail()
	}
	if !reflect.DeepEqual(m, m2) {
		t.Fail()
	}
}

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
	s := &Codec{}
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

func BenchmarkJsonCodec(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := Message{1, "hello world"}
		s := &Codec{}
		body, err := s.Encode(m)
		if err != nil {
			b.Fail()
		}
		m2 := Message{}
		if err := s.Decode(body, &m2); err != nil {
			b.Fail()
		}
		if !reflect.DeepEqual(m, m2) {
			b.Fail()
		}
	}
}

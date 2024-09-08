package json

import (
	"reflect"
	"testing"
)

type Message struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

func TestSerializer_Serialize(t *testing.T) {
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

func BenchmarkSerializer_Deserialize(b *testing.B) {
	m := &Message{100, "hell world"}
	s := &JsonCodec{}

	d, err := s.Encode(m)
	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		m1 := &Message{}
		if err := s.Decode(d, m1); err != nil {
			b.Fatalf("unmarshal failed: %v", err)
		}
	}
}

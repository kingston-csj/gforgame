package structcodec

import (
	"reflect"
	"testing"
)

type Address struct {
	City string
	Zip  int32
}

type Person struct {
	Name   string
	Age    int32
	Active bool
	Height float64
	Scores []int32
	Tags   map[string]string
	Addr   Address
}

func TestStructRoundTripSimple(t *testing.T) {
	in := Person{
		Name:   "测试🌟",
		Age:    28,
		Active: true,
		Height: 1.75,
		Scores: []int32{100, 98, 95},
		Tags:   map[string]string{"role": "admin", "lang": "zh"},
		Addr:   Address{City: "Shanghai", Zip: 200000},
	}
	ser := NewSerializer()
	data, err := ser.Encode(in)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}
	var out Person
	if err := ser.Decode(data, &out); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if !reflect.DeepEqual(in, out) {
		t.Fatalf("roundtrip mismatch\nin:  %+v\nout: %+v", in, out)
	}
	print(out.Name)
}

func TestStringCodecUtf8(t *testing.T) {
	ser := NewSerializer()
	in := struct{ S string }{S: "你好，世界🌍"}
	data, err := ser.Encode(in)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}
	var out struct{ S string }
	if err := ser.Decode(data, &out); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if in.S != out.S {
		t.Fatalf("string mismatch: in=%q out=%q", in.S, out.S)
	}
}
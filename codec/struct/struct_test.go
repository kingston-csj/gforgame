package structcodec_test

import (
	structcodec "io/github/gforgame/codec/struct"
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
	Animals []Animal
}

type Animal interface {
	Sound() string
}

type Dog struct{
	Name string
}
func (d *Dog) Sound() string { return d.Name + "汪汪" }

type Cat struct{
	Name2 string
}
func (c *Cat) Sound() string { return c.Name2 + "喵喵" }

var animalSlice []Animal = []Animal{&Dog{Name: "旺财"}, &Cat{Name2: "咪咪"}}


func TestStructRoundTripSimple(t *testing.T) {
	in := Person{
		Name:   "测试🌟",
		Age:    28,
		Active: true,
		Height: 1.75,
		Scores: []int32{100, 98, 95},
		Tags:   map[string]string{"role": "admin", "lang": "zh"},
		Addr:   Address{City: "Shanghai", Zip: 200000},
		Animals: animalSlice,
	}
	ser := structcodec.NewSerializer()
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
	for _, animal := range out.Animals {
		println(animal.Sound())
	}
}

// func TestStringCodecUtf8(t *testing.T) {
// 	ser := structcodec.NewSerializer()
// 	in := struct{ S string }{S: "你好，世界🌍"}
// 	data, err := ser.Encode(in)
// 	if err != nil {
// 		t.Fatalf("encode error: %v", err)
// 	}
// 	var out struct{ S string }
// 	if err := ser.Decode(data, &out); err != nil {
// 		t.Fatalf("decode error: %v", err)
// 	}
// 	if in.S != out.S {
// 		t.Fatalf("string mismatch: in=%q out=%q", in.S, out.S)
// 	}
// }

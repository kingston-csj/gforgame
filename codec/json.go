package codec

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type JsonCodec struct {
}

func (*JsonCodec) Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Decode 将byte数组反序列化为bean
func (*JsonCodec) Decode(data []byte, v any) error {
	// 反射用于检查接口类型并进行解码
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		// 确保解码到指针类型
		return json.Unmarshal(data, v)
	}
	// 如果不是指针类型，返回错误
	return fmt.Errorf("decode need a pointer type")
}

func NewJsonCodec() JsonCodec {
	return JsonCodec{}
}

package structcodec

import (
	"bytes"
	"errors"
	"reflect"
)

// StringCodec 字符串编解码器，使用 UTF-8，并以 uint16 长度前缀表示
type StringCodec struct{}

// Decode 读取 uint16 长度与 UTF-8 字节，返回 string
func (*StringCodec) Decode(r *bytes.Reader, typ reflect.Type) (any, error) {
    s, err := readUtf8(r)
    if err != nil {
        return nil, err
    }
    return s, nil
}

// Encode 将字符串按 UTF-8 写出，带 uint16 长度前缀
func (*StringCodec) Encode(w *bytes.Buffer, value any) error {
    s, ok := value.(string)
    if !ok {
        return errors.New("string codec on wrong type")
    }
    return writeUtf8(w, s)
}
package structcodec

import (
	"bytes"
	"errors"
	"reflect"
)

// BoolCodec 布尔类型编解码器，使用 1/0 表示 true/false
type BoolCodec struct{}

// Decode 从 1 字节读取布尔值
func (*BoolCodec) Decode(r *bytes.Reader, typ reflect.Type) (any, error) {
    b := make([]byte, 1)
    if _, err := r.Read(b); err != nil {
        return nil, err
    }
    return b[0] == 1, nil
}

// Encode 将布尔值写为 1 或 0
func (*BoolCodec) Encode(w *bytes.Buffer, value any) error {
    if b, ok := value.(bool); ok {
        if b {
            return w.WriteByte(1)
        }
        return w.WriteByte(0)
    }
    return errors.New("bool codec on wrong type")
}
package structcodec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
)

// Int32Codec 32 位整型编解码器，使用大端序
type Int32Codec struct{}

// Decode 以大端序读取 int32
func (*Int32Codec) Decode(r *bytes.Reader, typ reflect.Type) (any, error) {
    var v int32
    if err := binary.Read(r, binary.BigEndian, &v); err != nil {
        return nil, err
    }
    return v, nil
}

// Encode 以大端序写入 int32
func (*Int32Codec) Encode(w *bytes.Buffer, value any) error {
    v, ok := toInt32(value)
    if !ok {
        return errors.New("int32 codec on wrong type")
    }
    return binary.Write(w, binary.BigEndian, v)
}
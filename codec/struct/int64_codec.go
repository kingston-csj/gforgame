package structcodec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
)

// Int64Codec 64 位整型编解码器，使用大端序
type Int64Codec struct{}

// Decode 以大端序读取 int64
func (*Int64Codec) Decode(r *bytes.Reader, typ reflect.Type) (any, error) {
    var v int64
    if err := binary.Read(r, binary.BigEndian, &v); err != nil {
        return nil, err
    }
    return v, nil
}

// Encode 以大端序写入 int64
func (*Int64Codec) Encode(w *bytes.Buffer, value any) error {
    v, ok := toInt64(value)
    if !ok {
        return errors.New("int64 codec on wrong type")
    }
    return binary.Write(w, binary.BigEndian, v)
}

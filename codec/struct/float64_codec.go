package structcodec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"reflect"
)

// Float64Codec 64 位浮点编解码器，按 IEEE 754 大端序编码
type Float64Codec struct{}

// Decode 读取 uint64 位模式并转换为 float64
func (*Float64Codec) Decode(r *bytes.Reader, typ reflect.Type) (any, error) {
    var u uint64
    if err := binary.Read(r, binary.BigEndian, &u); err != nil {
        return nil, err
    }
    return math.Float64frombits(u), nil
}

// Encode 将 float64 转为位模式并以大端序写入
func (*Float64Codec) Encode(w *bytes.Buffer, value any) error {
    f, ok := toFloat64(value)
    if !ok {
        return errors.New("float64 codec on wrong type")
    }
    u := math.Float64bits(f)
    return binary.Write(w, binary.BigEndian, u)
}
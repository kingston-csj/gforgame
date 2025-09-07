package structcodec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"reflect"
)

// Float32Codec 32 位浮点编解码器，按 IEEE 754 大端序编码
type Float32Codec struct{}

// Decode 读取 uint32 位模式并转换为 float32
func (*Float32Codec) Decode(r *bytes.Reader, typ reflect.Type) (any, error) {
    var u uint32
    if err := binary.Read(r, binary.BigEndian, &u); err != nil {
        return nil, err
    }
    return math.Float32frombits(u), nil
}

// Encode 将 float32 转为位模式并以大端序写入
func (*Float32Codec) Encode(w *bytes.Buffer, value any) error {
    f, ok := toFloat32(value)
    if !ok {
        return errors.New("float32 codec on wrong type")
    }
    u := math.Float32bits(f)
    return binary.Write(w, binary.BigEndian, u)
}
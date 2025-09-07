package structcodec

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
)

// FieldCodec 定义字段级编解码器接口，针对某一具体类型实现值的读写
type FieldCodec interface {
    // Decode 从字节流中读取并还原为给定类型的值
    Decode(r *bytes.Reader, typ reflect.Type) (any, error)
    // Encode 将值按该类型的编码规则写入字节缓冲区
    Encode(w *bytes.Buffer, value any) error
}



// FieldCodecMeta 缓存结构体字段的类型与对应编解码器，降低重复反射成本
type FieldCodecMeta struct {
    index int
    typ   reflect.Type
    codec FieldCodec
}

// valueOfField 为结构体字段选择合适的 FieldCodec
func valueOfField(sf reflect.StructField) (FieldCodecMeta, error) {
    typ := sf.Type
    c, err := getFieldCodec(typ)
    if err != nil {
        return FieldCodecMeta{}, err
    }
    return FieldCodecMeta{typ: typ, codec: c}, nil
}

// getFieldCodec 根据反射类型返回对应的字段编解码器
func getFieldCodec(t reflect.Type) (FieldCodec, error) {
    switch t.Kind() {
    case reflect.Bool:
        return &BoolCodec{}, nil
    case reflect.String:
        return &StringCodec{}, nil
    case reflect.Int:
        return &Int64Codec{}, nil
    case reflect.Int8, reflect.Int16, reflect.Int32:
        return &Int32Codec{}, nil
    case reflect.Int64:
        return &Int64Codec{}, nil
    case reflect.Uint, reflect.Uint64:
        return &Int64Codec{}, nil
    case reflect.Uint8, reflect.Uint16, reflect.Uint32:
        return &Int32Codec{}, nil
    case reflect.Float32:
        return &Float32Codec{}, nil
    case reflect.Float64:
        return &Float64Codec{}, nil
    case reflect.Slice, reflect.Array:
        return &ArrayCodec{}, nil
    case reflect.Map:
        return &MapCodec{}, nil
    case reflect.Struct:
        return &BeanCodec{}, nil
    default:
        return nil, errors.New("unsupported type: " + t.String())
    }
}

// writeShort 以大端序写入 uint16
func writeShort(w *bytes.Buffer, v uint16) error {
    return binary.Write(w, binary.BigEndian, v)
}

// readShort 以大端序读取 uint16
func readShort(r *bytes.Reader) (uint16, error) {
    var v uint16
    err := binary.Read(r, binary.BigEndian, &v)
    return v, err
}

// writeUtf8 以 uint16 长度前缀写入 UTF-8 字符串
func writeUtf8(w *bytes.Buffer, s string) error {
    b := []byte(s)
    if len(b) > int(^uint16(0)) {
        return errors.New("string length exceed max uint16 value")
    }
    if err := writeShort(w, uint16(len(b))); err != nil {
        return err
    }
    _, err := w.Write(b)
    return err
}

// readUtf8 读取以 uint16 长度前缀编码的 UTF-8 字符串
func readUtf8(r *bytes.Reader) (string, error) {
    n, err := readShort(r)
    if err != nil {
        return "", err
    }
    if n == 0 {
        return "", nil
    }
    b := make([]byte, int(n))
    if _, err := r.Read(b); err != nil {
        return "", err
    }
    return string(b), nil
}

// toInt32 将多种整型转换为 int32（编解码内部使用）
func toInt32(v any) (int32, bool) {
    switch x := v.(type) {
    case int32:
        return x, true
    case int:
        return int32(x), true
    case uint32:
        return int32(x), true
    case uint16:
        return int32(x), true
    case int16:
        return int32(x), true
    case int8:
        return int32(x), true
    case uint8:
        return int32(x), true
    default:
        return 0, false
    }
}

// toInt64 将多种整型转换为 int64（编解码内部使用）
func toInt64(v any) (int64, bool) {
    switch x := v.(type) {
    case int64:
        return x, true
    case int:
        return int64(x), true
    case uint64:
        return int64(x), true
    case int32:
        return int64(x), true
    case uint32:
        return int64(x), true
    case int16:
        return int64(x), true
    case uint16:
        return int64(x), true
    case int8:
        return int64(x), true
    case uint8:
        return int64(x), true
    default:
        return 0, false
    }
}

// toFloat32 将 float64/float32 转换为 float32（编解码内部使用）
func toFloat32(v any) (float32, bool) {
    switch x := v.(type) {
    case float32:
        return x, true
    case float64:
        return float32(x), true
    default:
        return 0, false
    }
}

// toFloat64 将 float64/float32 转换为 float64（编解码内部使用）
func toFloat64(v any) (float64, bool) {
    switch x := v.(type) {
    case float64:
        return x, true
    case float32:
        return float64(x), true
    default:
        return 0, false
    }
}
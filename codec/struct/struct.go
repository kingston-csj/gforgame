package structcodec

import (
	"bytes"
	"errors"
	"reflect"
)

// 基于struct定义的消息编码解码器
// 按照struct申明的字段类型以及相应的顺序，将一个struct与byte[]相互转换
type Codec struct{}

func NewSerializer() *Codec {
    return &Codec{}
}

func (c *Codec) Encode(v any) ([]byte, error) {
    buf := &bytes.Buffer{}
    t := reflect.TypeOf(v)
    fc, err := getFieldCodec(t)
    if err != nil {
        return nil, err
    }
    if err := fc.Encode(buf, v); err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}

func (c *Codec) Decode(data []byte, v any) error {
    r := bytes.NewReader(data)
    rv := reflect.ValueOf(v)
    if rv.Kind() != reflect.Ptr || rv.IsNil() {
        return errors.New("decode need a non-nil pointer")
    }
    et := rv.Type().Elem()
    fc, err := getFieldCodec(et)
    if err != nil {
        return err
    }
    val, err := fc.Decode(r, et)
    if err != nil {
        return err
    }
    rv.Elem().Set(reflect.ValueOf(val))
    return nil
}
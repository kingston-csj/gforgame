package structcodec

import (
	"bytes"
	"errors"
	"reflect"
)

// MapCodec 映射编解码器，仅支持 map[string]V；长度以 uint16 存储
type MapCodec struct{}

// Decode 读取长度后按 key(string) + value(V) 的顺序解码
func (*MapCodec) Decode(r *bytes.Reader, typ reflect.Type) (any, error) {
    if typ.Key().Kind() != reflect.String {
        return nil, errors.New("map key must be string")
    }
    n, err := readShort(r)
    if err != nil {
        return nil, err
    }
    size := int(n)
    vt := typ.Elem()
    fc, err := getFieldCodec(vt)
    if err != nil {
        return nil, err
    }
    mv := reflect.MakeMapWithSize(typ, size)
    for i := 0; i < size; i++ {
        k, err := readUtf8(r)
        if err != nil {
            return nil, err
        }
        v, err := fc.Decode(r, vt)
        if err != nil {
            return nil, err
        }
        mv.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
    }
    return mv.Interface(), nil
}

// Encode 写入长度并依次写入 key(string) 与 value(V)
func (*MapCodec) Encode(w *bytes.Buffer, value any) error {
    if value == nil {
        return writeShort(w, 0)
    }
    rv := reflect.ValueOf(value)
    if rv.Kind() != reflect.Map || rv.Type().Key().Kind() != reflect.String {
        return errors.New("map codec requires map[string]V")
    }
    keys := rv.MapKeys()
    size := len(keys)
    if size > int(^uint16(0)) {
        return errors.New("collection size exceed max uint16 value")
    }
    if err := writeShort(w, uint16(size)); err != nil {
        return err
    }
    vt := rv.Type().Elem()
    fc, err := getFieldCodec(vt)
    if err != nil {
        return err
    }
    for _, k := range keys {
        if err := writeUtf8(w, k.String()); err != nil {
            return err
        }
        v := rv.MapIndex(k).Interface()
        if err := fc.Encode(w, v); err != nil {
            return err
        }
    }
    return nil
}
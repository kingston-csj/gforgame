package structcodec

import (
	"bytes"
	"errors"
	"reflect"
)

// ArrayCodec 数组/切片编解码器；长度以 uint16 存储，最大 65535
type ArrayCodec struct{}

// Decode 读取长度后逐元素解码为切片，元素类型由 typ.Elem() 决定
func (*ArrayCodec) Decode(r *bytes.Reader, typ reflect.Type) (any, error) {
    n, err := readShort(r)
    if err != nil {
        return nil, err
    }
    size := int(n)
    if size < 0 {
        return nil, errors.New("array size less than zero")
    }
    elemType := typ.Elem()
    fc, err := getFieldCodec(elemType)
    if err != nil {
        return nil, err
    }
    slice := reflect.MakeSlice(reflect.SliceOf(elemType), size, size)
    for i := 0; i < size; i++ {
        ev, err := fc.Decode(r, elemType)
        if err != nil {
            return nil, err
        }
        slice.Index(i).Set(reflect.ValueOf(ev))
    }
    return slice.Interface(), nil
}

// Encode 写入长度并逐元素编码
func (*ArrayCodec) Encode(w *bytes.Buffer, value any) error {
    if value == nil {
        return writeShort(w, 0)
    }
    rv := reflect.ValueOf(value)
    switch rv.Kind() {
    case reflect.Array, reflect.Slice:
    default:
        return errors.New("array codec on non array/slice")
    }
    size := rv.Len()
    if size > int(^uint16(0)) {
        return errors.New("collection size exceed max uint16 value")
    }
    if err := writeShort(w, uint16(size)); err != nil {
        return err
    }
    for i := 0; i < size; i++ {
        elem := rv.Index(i).Interface()
        et := rv.Type().Elem()
        fc, err := getFieldCodec(et)
        if err != nil {
            return err
        }
        if err := fc.Encode(w, elem); err != nil {
            return err
        }
    }
    return nil
}

package structcodec

import (
	"bytes"
	"errors"
	"reflect"
	"sync"
)

// BeanCodec 结构体（bean）编解码器，按导出字段声明顺序处理
type BeanCodec struct{}

// beanMetaCache 缓存类型到字段元信息的映射，避免重复构建
var beanMetaCache sync.Map

// buildBeanMeta 构建并缓存结构体字段的编解码元信息
func buildBeanMeta(t reflect.Type) ([]FieldCodecMeta, error) {
    if t.Kind() != reflect.Struct {
        return nil, errors.New("bean codec on non struct")
    }
    if v, ok := beanMetaCache.Load(t); ok {
        return v.([]FieldCodecMeta), nil
    }
    var metas []FieldCodecMeta
    for i := 0; i < t.NumField(); i++ {
        sf := t.Field(i)
        if sf.PkgPath != "" {
            continue
        }
        m, err := valueOfField(sf)
        if err != nil {
            return nil, err
        }
        m.index = i
        metas = append(metas, m)
    }
    beanMetaCache.Store(t, metas)
    return metas, nil
}

// Decode 逐字段解码并返回结构体值
func (*BeanCodec) Decode(r *bytes.Reader, typ reflect.Type) (any, error) {
    metas, err := buildBeanMeta(typ)
    if err != nil {
        return nil, err
    }
    ptr := reflect.New(typ)
    for _, m := range metas {
        val, err := m.codec.Decode(r, m.typ)
        if err != nil {
            return nil, err
        }
        ptr.Elem().Field(m.index).Set(reflect.ValueOf(val))
    }
    return ptr.Elem().Interface(), nil
}

// Encode 逐字段编码结构体值
func (*BeanCodec) Encode(w *bytes.Buffer, value any) error {
    rv := reflect.ValueOf(value)
    if rv.Kind() != reflect.Struct {
        return errors.New("bean codec on non struct")
    }
    metas, err := buildBeanMeta(rv.Type())
    if err != nil {
        return err
    }
    for _, m := range metas {
        fv := rv.Field(m.index).Interface()
        if err := m.codec.Encode(w, fv); err != nil {
            return err
        }
    }
    return nil
}

package data

import (
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type Converter interface {
	Convert(source string) (any, error)
}

type ConvertFunc func(source string) (any, error)

func (f ConvertFunc) Convert(source string) (any, error) {
	return f(source)
}

type ConvertRegistry struct {
	mu         sync.RWMutex
	converters map[reflect.Type]Converter
}

var (
	ConvertRegistryInstance = NewConvertRegistry()
)

func NewConvertRegistry() *ConvertRegistry {
	return &ConvertRegistry{
		converters: make(map[reflect.Type]Converter),
	}
}

func (r *ConvertRegistry) Register(targetType reflect.Type, converter Converter) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.converters[targetType] = converter
}

func (r *ConvertRegistry) RegisterFunc(targetType reflect.Type, convertFunc func(string) (any, error)) {
	r.Register(targetType, ConvertFunc(convertFunc))
}

func (r *ConvertRegistry) Convert(source string, targetType reflect.Type) (any, error) {
	r.mu.RLock()
	converter, ok := r.converters[targetType]
	r.mu.RUnlock()

	if ok {
		return converter.Convert(source)
	}

	if targetType.Kind() == reflect.Slice {
		return convertSlice(source, targetType)
	}

	return nil, nil
}

func convertSlice(source string, targetType reflect.Type) (any, error) {
	if source == "" {
		return reflect.MakeSlice(targetType, 0, 0).Interface(), nil
	}

	elements := strings.Split(source, ",")
	elemType := targetType.Elem()
	result := reflect.MakeSlice(targetType, len(elements), len(elements))

	for i, elem := range elements {
		elem = strings.TrimSpace(elem)
		var val reflect.Value

		switch elemType.Kind() {
		case reflect.String:
			val = reflect.ValueOf(elem)
		case reflect.Int:
			v, err := strconv.Atoi(elem)
			if err != nil {
				return nil, err
			}
			val = reflect.ValueOf(v)
		case reflect.Int8:
			v, err := strconv.ParseInt(elem, 10, 8)
			if err != nil {
				return nil, err
			}
			val = reflect.ValueOf(int8(v))
		case reflect.Int16:
			v, err := strconv.ParseInt(elem, 10, 16)
			if err != nil {
				return nil, err
			}
			val = reflect.ValueOf(int16(v))
		case reflect.Int32:
			v, err := strconv.ParseInt(elem, 10, 32)
			if err != nil {
				return nil, err
			}
			val = reflect.ValueOf(int32(v))
		case reflect.Int64:
			v, err := strconv.ParseInt(elem, 10, 64)
			if err != nil {
				return nil, err
			}
			val = reflect.ValueOf(v)
		case reflect.Uint:
			v, err := strconv.ParseUint(elem, 10, 0)
			if err != nil {
				return nil, err
			}
			val = reflect.ValueOf(uint(v))
		case reflect.Uint8:
			v, err := strconv.ParseUint(elem, 10, 8)
			if err != nil {
				return nil, err
			}
			val = reflect.ValueOf(uint8(v))
		case reflect.Uint16:
			v, err := strconv.ParseUint(elem, 10, 16)
			if err != nil {
				return nil, err
			}
			val = reflect.ValueOf(uint16(v))
		case reflect.Uint32:
			v, err := strconv.ParseUint(elem, 10, 32)
			if err != nil {
				return nil, err
			}
			val = reflect.ValueOf(uint32(v))
		case reflect.Uint64:
			v, err := strconv.ParseUint(elem, 10, 64)
			if err != nil {
				return nil, err
			}
			val = reflect.ValueOf(v)
		case reflect.Float32:
			v, err := strconv.ParseFloat(elem, 32)
			if err != nil {
				return nil, err
			}
			val = reflect.ValueOf(float32(v))
		case reflect.Float64:
			v, err := strconv.ParseFloat(elem, 64)
			if err != nil {
				return nil, err
			}
			val = reflect.ValueOf(v)
		case reflect.Bool:
			v, err := strconv.ParseBool(elem)
			if err != nil {
				return nil, err
			}
			val = reflect.ValueOf(v)
		default:
			return nil, nil
		}

		result.Index(i).Set(val)
	}

	return result.Interface(), nil
}

func (r *ConvertRegistry) MustConvert(source string, targetType reflect.Type) any {
	result, err := r.Convert(source, targetType)
	if err != nil {
		panic(err)
	}
	return result
}
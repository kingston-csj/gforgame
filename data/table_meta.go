package data

import (
	"fmt"
	"reflect"
)

type TableMeta struct {
	TableName     string            // 表名
	IDField       string            // ID 字段名
	IndexFuncs    map[string]string // 索引字段名 -> 索引名称
	RecordType    reflect.Type      // 记录类型
	ContainerType reflect.Type      // 容器类型
}

func ProcessTable(reader *ExcelDataReader, filePath string, config TableMeta) (interface{}, error) {
	// 读取 Excel 文件
	records, err := reader.Read(filePath, reflect.New(config.RecordType).Interface())
	if err != nil {
		return nil, fmt.Errorf("failed to read table %s: %v", config.TableName, err)
	}

	// 创建容器实例
	var container interface{}
	if config.ContainerType != nil {
		containerValue := reflect.New(config.ContainerType.Elem())
		container = containerValue.Interface()

		// 调用Init和AfterLoad方法
		if initializer, ok := container.(IContainer); ok {
			initializer.Init()
		}
	} else {
		// 获取ID字段的类型
		field, ok := config.RecordType.FieldByName(config.IDField)
		if !ok {
			return nil, fmt.Errorf("field %s not found in type %s", config.IDField, config.RecordType.Name())
		}

		// 根据字段类型创建对应的容器
		switch field.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			container = NewContainer[int64, any]()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			container = NewContainer[uint64, any]()
		case reflect.String:
			container = NewContainer[string, any]()
		default:
			return nil, fmt.Errorf("unsupported ID field type: %v", field.Type)
		}
	}

	// 创建 ID 获取函数
	getIdFunc := func(record any) any {
		val := reflect.ValueOf(record)
		if val.Kind() == reflect.Interface {
			val = val.Elem()
		}
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		return val.FieldByName(config.IDField).Interface()
	}

	// 创建索引函数
	indexFuncs := make(map[string]func(any) any)
	if config.IndexFuncs != nil {
		for indexName, fieldName := range config.IndexFuncs {
			indexFuncs[indexName] = func(record any) any {
				val := reflect.ValueOf(record)
				if val.Kind() == reflect.Interface {
					val = val.Elem()
				}
				if val.Kind() == reflect.Ptr {
					val = val.Elem()
				}
				return val.FieldByName(fieldName).Interface()
			}
		}
	}

	// 注入数据到容器
	containerValue := reflect.ValueOf(container)
	injectMethod := containerValue.MethodByName("Inject")
	injectMethod.Call([]reflect.Value{
		reflect.ValueOf(records),
		reflect.ValueOf(getIdFunc),
		reflect.ValueOf(indexFuncs),
	})

	if initializer, ok := container.(IContainer); ok {
		initializer.AfterLoad()
	}

	return container, nil
}

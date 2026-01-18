package data

import (
	"fmt"
	"io/github/gforgame/logger"
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
	logger.Info(fmt.Sprintf("processed table %s, %d records", config.TableName, len(records)))

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
		if field.Type.Kind() != reflect.Int32 {
			return nil, fmt.Errorf("ID field %s must be int32 type", config.IDField)
		}
		container = NewContainer[int32, any]()
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

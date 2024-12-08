package data

import (
	"fmt"
	"reflect"
)

type TableMeta struct {
	TableName  string            // 表名
	IDField    string            // ID 字段名
	IndexFuncs map[string]string // 索引字段名 -> 索引名称
	RecordType reflect.Type      // 记录类型
}

func ProcessTable(reader *ExcelDataReader, filePath string, config TableMeta) (*Container[int64, interface{}], error) {
	// 读取 Excel 文件
	records, err := reader.Read(filePath, reflect.New(config.RecordType).Interface())
	if err != nil {
		return nil, fmt.Errorf("failed to read table %s: %v", config.TableName, err)
	}

	// 创建 Container
	container := NewContainer[int64, interface{}]()

	// 定义 ID 获取函数
	getIdFunc := func(record interface{}) int64 {
		val := reflect.ValueOf(record)
		// 如果 record 是指针，则调用 Elem() 获取实际值
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		field := val.FieldByName(config.IDField)
		return field.Int()
	}

	// 定义索引函数
	indexFuncs := make(map[string]func(interface{}) interface{})
	if config.IndexFuncs != nil {
		for indexName, fieldName := range config.IndexFuncs {
			indexFuncs[indexName] = func(record interface{}) interface{} {
				val := reflect.ValueOf(record)
				// 如果 record 是指针，则调用 Elem() 获取实际值
				if val.Kind() == reflect.Ptr {
					val = val.Elem()
				}
				field := val.FieldByName(fieldName)
				return field.Interface()
			}
		}
	}

	// 将记录注入容器
	container.Inject(records, getIdFunc, indexFuncs)

	return container, nil
}

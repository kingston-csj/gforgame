package data

import (
	"fmt"
	"io/github/gforgame/logger"
	"io/github/gforgame/util/jsonutil"
	"os"
	"path/filepath"
	"reflect"
)

var (
	// 容器类型
	EnableJSONOutput bool   = false
	JSONOutputDir    string = "json" // JSON输出目录
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

	// 输出为 JSON 文件
	if EnableJSONOutput {
		if err := writeRecordsToJSON(records, config.TableName, JSONOutputDir); err != nil {
			logger.Error3(fmt.Sprintf("failed to write JSON for table %s: %v", config.TableName, err))
		}
	}

	// 创建容器实例
	var container interface{}
	if config.ContainerType != nil {
		containerValue := reflect.New(config.ContainerType.Elem())
		container = containerValue.Interface()

		// 调用Init和AfterLoad方法
		if initializer, ok := container.(IBaseContainer); ok {
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

	if initializer, ok := container.(IBaseContainer); ok {
		initializer.AfterLoad()
	}

	return container, nil
}

// writeRecordsToJSON 将记录写入 JSON 文件
func writeRecordsToJSON(records []interface{}, tableName, outputDir string) error {
	// 确保输出目录存在
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// 生成 JSON 字符串
	jsonStr, err := jsonutil.StructToPrettyJSON(records)
	if err != nil {
		return fmt.Errorf("failed to convert records to JSON: %v", err)
	}

	// 构建输出文件路径
	outputPath := filepath.Join(outputDir, tableName+".json")

	// 写入文件
	if err := os.WriteFile(outputPath, []byte(jsonStr), 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %v", err)
	}

	logger.Info(fmt.Sprintf("written table %s to JSON file: %s", tableName, outputPath))
	return nil
}

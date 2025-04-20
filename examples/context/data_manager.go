package context

import (
	"fmt"
	"reflect"
	"sync"

	"io/github/gforgame/data"
	domain "io/github/gforgame/examples/domain/config"
)

// GetConfigRecordAs returns a record of specific type
func GetConfigRecordAs[T any](name string, id int64) *T {
	record := GetDataManager().GetRecord(name, id)
	if record == nil {
		return nil
	}
	if result, ok := record.(*T); ok {
		return result
	}
	result := new(T)
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(record))
	return result
}

// DataManager[T any] is a generic data manager that can handle records of type T
type DataManager[T any] struct {
	containers map[string]*data.Container[int64, T]
}

var (
	instance *DataManager[any]
	once     sync.Once
)

func GetDataManager() *DataManager[any] {
	once.Do(func() {
		instance = &DataManager[any]{}

		// 创建 ExcelDataReader
		reader := data.NewExcelDataReader(true)

		// 定义表配置
		tableConfigs := []data.TableMeta{
			// 道具表
			{
				TableName:  "item",
				IDField:    "Id",
				RecordType: reflect.TypeOf(domain.ItemData{}),
			},
			// 英雄表
			{
				TableName:  "hero",
				IDField:    "Id",
				RecordType: reflect.TypeOf(domain.HeroData{}),
			},
			// 英雄等级表
			{
				TableName:  "herolevel",
				IDField:    "Id",
				RecordType: reflect.TypeOf(domain.HeroLevelData{}),
			},
			// 英雄阶段表
			{
				TableName:  "herostage",
				IDField:    "Id",
				RecordType: reflect.TypeOf(domain.HeroStageData{}),
			},
			// 技能表
			{
				TableName:  "skill",
				IDField:    "Id",
				RecordType: reflect.TypeOf(domain.SkillData{}),
			},
			// buff表
			{
				TableName:  "buff",
				IDField:    "Id",
				RecordType: reflect.TypeOf(domain.BuffData{}),
			},
		}

		// 处理每张表
		containers := make(map[string]*data.Container[int64, any])
		for _, config := range tableConfigs {
			container, err := data.ProcessTable(reader, config.TableName+".xlsx", config)
			if err != nil {
				fmt.Printf("Failed to process table %s: %v\n", config.TableName, err)
				panic(err)
			}
			containers[config.TableName] = container
		}

		instance.containers = containers
	})
	return instance
}

func (dm *DataManager[T]) GetRecord(name string, id int64) any {
	container := dm.containers[name]
	if container == nil {
		return nil
	}
	record, ok := container.GetRecord(id)
	if !ok {
		return nil
	}
	return record
}

func (dm *DataManager[T]) GetRecords(name string) []any {
	container := dm.containers[name]
	if container == nil {
		return nil
	}
	records := container.GetAllRecords()
	result := make([]any, len(records))
	for i, record := range records {
		result[i] = record
	}
	return result
}

func (dm *DataManager[T]) GetRecordsByIndex(configName string, indexName string, indexValue any) []any {
	container := dm.containers[configName]
	if container == nil {
		return nil
	}
	records := container.GetRecordsBy(indexName, indexValue)
	result := make([]any, len(records))
	for i, record := range records {
		result[i] = record
	}
	return result
}

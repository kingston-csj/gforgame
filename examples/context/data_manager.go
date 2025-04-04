package context

import (
	"fmt"
	"reflect"
	"sync"

	"io/github/gforgame/data"
	domain "io/github/gforgame/examples/domain/config"
)

type DataManager struct {
	containers map[string]*data.Container[int64, any]
}

var (
	instance *DataManager
	once     sync.Once
)

func GetDataManager() *DataManager {
	once.Do(func() {
		instance = &DataManager{}

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

func (dm *DataManager) GetRecord(name string, id int64) any {
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

func (dm *DataManager) GetRecords(name string) []any {
	container := dm.containers[name]
	if container == nil {
		return nil
	}
	return container.GetAllRecords()
}

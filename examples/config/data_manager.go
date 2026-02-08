package config

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"io/github/gforgame/data"
	"io/github/gforgame/examples/config/container"

	domain "io/github/gforgame/examples/domain/config"
)

// DataManager 配置数据管理器
type DataManager struct {
	containers map[string]any
}

var (
	instance *DataManager
	once     sync.Once
)

// table名称对应Meta
var tableConfigMap map[string]data.TableMeta
// 容器类型对应表名
var containerKeys map[reflect.Type]string

func init() {
	tableConfigMap = make(map[string]data.TableMeta)
	containerKeys = make(map[reflect.Type]string)
	// 定义表配置
	tableConfigs := []data.TableMeta{
		// 公共配置表
		{
			RecordType: reflect.TypeOf(domain.CommonData{}),
			ContainerType: reflect.TypeOf(&container.CommonContainer{}),
		},
		// 道具表
		{
			RecordType: reflect.TypeOf(domain.PropData{}),
		},
		// 英雄表
		{
			RecordType:    reflect.TypeOf(domain.HeroData{}),
			ContainerType: reflect.TypeOf(&data.Container[int32, domain.HeroData]{}),
		},
		// 英雄等级表
		{
			RecordType:    reflect.TypeOf(domain.HeroLevelData{}),
			ContainerType: reflect.TypeOf(&container.HeroLevelContainer{}),
		},
		// 英雄阶段表
		{
			RecordType:    reflect.TypeOf(domain.HeroStageData{}),
			ContainerType: reflect.TypeOf(&container.HeroStageContainer{}),
		},
		// 技能表
		{
			RecordType:    reflect.TypeOf(domain.SkillData{}),
			ContainerType: reflect.TypeOf(&data.Container[int32, domain.SkillData]{}),
		},
		// quest表
		{
			RecordType:    reflect.TypeOf(domain.QuestData{}),
			ContainerType: reflect.TypeOf(&container.QuestContainer{}),
			IndexFuncs:    map[string]string{"Category": "Category"},
		},
		// 活动表
		{
			RecordType:    reflect.TypeOf(domain.ActivityData{}),
		},
	

	}

	for _, config := range tableConfigs {
		if config.IDField == "" {
			config.IDField = "Id"
		}
		if config.TableName == "" {
			// 去掉"Data"后缀
			config.TableName = strings.ToLower(strings.ReplaceAll(config.RecordType.Name(), "Data", ""))
		}
		tableConfigMap[config.TableName] = config
		if config.ContainerType != nil {
			containerKeys[config.ContainerType] = config.TableName
		}
	}
}

// GetDataManager 获取单例实例
func GetDataManager() *DataManager {
	once.Do(func() {
		instance = &DataManager{
			containers: make(map[string]interface{}),
		}

		// 创建 ExcelDataReader
		reader := data.NewExcelDataReader(true)

		// 处理每张表
		for name, config := range tableConfigMap {
			container, err := data.ProcessTable(reader, name+".xlsx", config)
			if err != nil {
				fmt.Printf("Failed to process table %s: %v\n", name, err)
				panic(err)
			}
			instance.containers[name] = container
		}
	})
	return instance
}

// GetContainer 获取原始容器
func GetContainer(name string) interface{} {
	return GetDataManager().containers[name]
}
 
// GetSpecificContainer 获取特定类型的容器
func GetSpecificContainer[C any]() C {
	tableName := containerKeys[reflect.TypeOf((*C)(nil)).Elem()]
	if tableName == "" {
		var zero C
		return zero
	}
	container := GetContainer(tableName)
	if container == nil {
		var zero C
		return zero
	}
	if specific, ok := container.(C); ok {
		return specific
	}
	var zero C
	return zero
}


// QueryAll 查询指定类型的所有记录
func QueryAll[V any]() []*V {
	tableName := getTableName[V]()
	container := GetContainer(tableName)
	if container == nil {
		return nil
	}
	// 1. 尝试直接匹配泛型容器 (针对自定义容器，如 CommonContainer)
	if c, ok := container.(data.IContainer[int32, V]); ok {
		return c.GetAllRecords()
	}

	// 2. 尝试作为 IAnyContainer 处理 (针对默认容器 Container[int32, any])
	// 注意：这里必须断言为 IContainer[int32, any] 才能调用 GetAllRecords
	if c, ok := container.(data.IContainer[int32, any]); ok {
		anyRecords := c.GetAllRecords() // 返回 []*any
		results := make([]*V, 0, len(anyRecords))

		for _, ptrAny := range anyRecords {
			if ptrAny == nil {
				continue
			}
			val := *ptrAny // 获取 interface{}，内部可能是 V 或 *V

			// 类型断言
			if v, ok := val.(V); ok {
				results = append(results, &v)
			} else if vPtr, ok := val.(*V); ok {
				results = append(results, vPtr)
			}
		}
		return results
	}

	return nil
}
// QueryById 根据ID查询指定类型的记录
// 这段恶心的代码先凑合着用，后续再干掉
func QueryById[V any](id int32) *V {
	tableName := getTableName[V]()
	container := GetContainer(tableName)
	if container == nil {
		return nil
	}
	if c, ok := container.(data.IContainer[int32, V]); ok {
		return c.GetRecord(id)
	}

	// 尝试作为 IAnyContainer 处理 (兼容 Container[int32, any])
	if c, ok := container.(data.IAnyContainer); ok {
		val := c.GetRecordAny(id)
		if val == nil {
			return nil
		}

		// 1. 如果容器本身存储的就是目标类型的指针 (Container[int32, V])
		// 虽然前面的 IContainer 检查应该已经涵盖了这种情况，但为了保险起见保留
		if v, ok := val.(*V); ok {
			return v
		}

		// 2. 如果容器是 Container[int32, any]，则 val 是 *any
		if ptrAny, ok := val.(*any); ok {
			if ptrAny == nil {
				return nil
			}
			inner := *ptrAny // 获取 any 内部持有的值

			// 如果内部值是目标类型 V (struct)
			if v, ok := inner.(V); ok {
				return &v
			}
			// 如果内部值是目标类型的指针 *V
			if v, ok := inner.(*V); ok {
				return v
			}
		}
	}

	return nil

}

// QueryContainer 获取指定类型的容器
func QueryContainer[V any, C any]() C {
	tableName := getTableName[V]()
	container := GetContainer(tableName)
	if container == nil {
		var zero C
		return zero
	}
	if specific, ok := container.(C); ok {
		return specific
	}
	var zero C
	return zero
}

// getTableName 根据类型获取表名
func getTableName[V any]() string {
	t := reflect.TypeOf((*V)(nil)).Elem()
	// 移除Data后缀
	name := t.Name()
	if strings.HasSuffix(name, "Data") {
		name = name[:len(name)-4]
	}
	return strings.ToLower(name)
}

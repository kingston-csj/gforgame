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
	containers map[string]interface{}
}

var (
	instance *DataManager
	once     sync.Once
)

var tableConfigMap map[string]data.TableMeta

func init() {
	tableConfigMap = make(map[string]data.TableMeta)
	// 定义表配置
	tableConfigs := []data.TableMeta{
		// 公共配置表
		{
			TableName:  "common",
			IDField:    "Id",
			RecordType: reflect.TypeOf(domain.CommonData{}),
			ContainerType: reflect.TypeOf(&container.CommonContainer{}),
		},
		// 道具表
		{
			TableName:  "prop",
			IDField:    "Id",
			RecordType: reflect.TypeOf(domain.PropData{}),
		},
		// 英雄表
		{
			TableName:     "hero",
			IDField:       "Id",
			RecordType:    reflect.TypeOf(domain.HeroData{}),
			ContainerType: reflect.TypeOf(&data.Container[int32, domain.HeroData]{}),
		},
		// 英雄等级表
		{
			TableName:     "herolevel",
			IDField:       "Id",
			RecordType:    reflect.TypeOf(domain.HeroLevelData{}),
			ContainerType: reflect.TypeOf(&container.HeroLevelContainer{}),
		},
		// 英雄阶段表
		{
			TableName:     "herostage",
			IDField:       "Id",
			RecordType:    reflect.TypeOf(domain.HeroStageData{}),
			ContainerType: reflect.TypeOf(&container.HeroStageContainer{}),
		},
		// 技能表
		{
			TableName:     "skill",
			IDField:       "Id",
			RecordType:    reflect.TypeOf(domain.SkillData{}),
			ContainerType: reflect.TypeOf(&data.Container[int32, domain.SkillData]{}),
		},
		// quest表
		{
			TableName:     "quest",
			IDField:       "Id",
			RecordType:    reflect.TypeOf(domain.QuestData{}),
			ContainerType: reflect.TypeOf(&container.QuestContainer{}),
			IndexFuncs:    map[string]string{"Category": "Category"},
		},
		// 活动表
		{
			TableName:     "activity",
			IDField:       "Id",
			RecordType:    reflect.TypeOf(domain.ActivityData{}),
			ContainerType: reflect.TypeOf(&data.Container[int32, domain.ActivityData]{}),
		},
		// // 每日签到奖励表
		// {
		// 	TableName:     "signin",
		// 	IDField:       "Id",
		// 	RecordType:    reflect.TypeOf(domain.SigninData{}),
		// },
		// // 充值表
		// {
		// 	TableName:     "recharge",
		// 	IDField:       "Id",
		// 	RecordType:    reflect.TypeOf(domain.RechargeData{}),
		// },
		// // 商城表
		// {
		// 	TableName:     "mall",
		// 	IDField:       "Id",
		// 	RecordType:    reflect.TypeOf(domain.MallData{}),
		// },
		// // 月卡表
		// {
		// 	TableName:     "monthlycard",
		// 	IDField:       "Id",
		// 	RecordType:    reflect.TypeOf(domain.MonthlyCardData{}),
		// },
		// // 抽奖表
		// {
		// 	TableName:     "gacha",
		// 	IDField:       "Id",
		// 	RecordType:    reflect.TypeOf(domain.GachaData{}),
		// 	ContainerType: reflect.TypeOf(&container.GachaContainer{}),
		// },

	}

	for _, config := range tableConfigs {
		tableConfigMap[config.TableName] = config
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
func (dm *DataManager) GetContainer(name string) interface{} {
	return dm.containers[name]
}

// GetSpecificContainer 获取特定类型的容器
func GetSpecificContainer[C any](name string) *C {
	container := GetDataManager().GetContainer(name)
	if container == nil {
		return nil
	}
	if specific, ok := container.(*C); ok {
		return specific
	}
	return nil
}

// QueryById 根据ID查询指定类型的记录
// 这段恶心的代码先凑合着用，后续再干掉
func QueryById[V any](id int32) *V {
	tableName := getTableName[V]()
	container := GetDataManager().GetContainer(tableName)
	if container == nil {
		return nil
	}
	// 尝试调用GetRecord方法
	if method := reflect.ValueOf(container).MethodByName("GetRecord"); method.IsValid() {
		// 获取方法的参数类型
		methodType := method.Type()
		if methodType.NumIn() != 1 {
			return nil
		}
		paramType := methodType.In(0)

		// 转换id到正确的类型
		idValue := reflect.ValueOf(id)
		var convertedValue reflect.Value

		if idValue.Type().ConvertibleTo(paramType) {
			convertedValue = idValue.Convert(paramType)
		}

		if !convertedValue.IsValid() {
			return nil
		}

		results := method.Call([]reflect.Value{convertedValue})
		if len(results) > 0 && !results[0].IsNil() {
			if record := results[0].Interface(); record != nil {
				// 如果record是*any类型，需要先获取其指向的值
				recordValue := reflect.ValueOf(record)
				if recordValue.Kind() == reflect.Ptr && recordValue.Elem().Type() == reflect.TypeOf((*any)(nil)).Elem() {
					// 获取*any指向的实际值
					actualValue := recordValue.Elem().Interface()
					if actualValue != nil {
						// 将actualValue转换为*V
						actualValuePtr := reflect.ValueOf(actualValue)
						if actualValuePtr.Kind() == reflect.Ptr {
							if actualValuePtr.Type().Elem() == reflect.TypeOf((*V)(nil)).Elem() {
								return actualValue.(*V)
							}
						}
						// 如果不是指针，尝试看看是不是直接就是V类型
						if actualValuePtr.Type() == reflect.TypeOf((*V)(nil)).Elem() {
							// 将actualValue转换为V类型，然后获取地址
							v := actualValue.(V)
							return &v
						}
					}
				} else {
					// 直接尝试类型转换
					if v, ok := record.(*V); ok {
						return v
					}
				}
			}
		}
	}

	return nil
}

// QueryContainer 获取指定类型的容器
func QueryContainer[V any, C any]() C {
	tableName := getTableName[V]()
	container := GetDataManager().GetContainer(tableName)
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

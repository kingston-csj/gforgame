package config

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/data"
	"github.com/forfun/gforgame/internal/config/container"

	domain "github.com/forfun/gforgame/internal/domain/config"
)

// DataManager 配置数据管理器
type DataManager struct {
	containers map[string]any
}

type configTableLoader struct {
	tableName     string
	containerType reflect.Type
	load          func(reader data.DataReader) (any, error)
}

// 表名称对应强类型加载器
var tableLoaders map[string]configTableLoader

// 容器类型对应表名
var containerKeys map[reflect.Type]string

// 全局唯一实例（包私有）
var global *DataManager

func init() {
	tableLoaders = make(map[string]configTableLoader)
	containerKeys = make(map[reflect.Type]string)
	// TODO 联合索引
	tableConfigs := []configTableLoader{
		newContainerTableLoader[domain.CommonData](func() *container.CommonContainer { return &container.CommonContainer{} }, nil),
		newDefaultTableLoader[domain.PropData](nil, nil),
		newDefaultTableLoader[domain.HeroData](nil, nil),

		newDefaultTableLoader[domain.PlayerLevelData](nil, map[string]string{"player_level": "Id,Level"}),
		newContainerTableLoader[domain.HeroLevelData](func() *container.HeroLevelContainer { return &container.HeroLevelContainer{} }, nil),
		newDefaultTableLoader[domain.HeroStageData](nil, map[string]string{"hero_stage": "TargetId,Stage"}),

		newDefaultTableLoader[domain.SkillData](reflect.TypeOf(&data.Container[int32, domain.SkillData]{}), nil),
		newContainerTableLoader[domain.QuestData](func() *container.QuestContainer { return &container.QuestContainer{} }, map[string]string{"Category": "Category"}),
		newDefaultTableLoader[domain.ActivityData](nil, nil),
		newDefaultTableLoader[domain.ActivityRewardData](nil, map[string]string{"ActivityId": "ActivityId"}),
		newDefaultTableLoader[domain.SigninData](nil, nil),
		newDefaultTableLoader[domain.RechargeData](nil, nil),
		newDefaultTableLoader[domain.MallData](nil, nil),
		newDefaultTableLoader[domain.MailData](nil, nil),
		newDefaultTableLoader[domain.MonthlyCardData](nil, nil),
		newContainerTableLoader[domain.GachaData](func() *container.GachaContainer { return &container.GachaContainer{} }, nil),
		newContainerTableLoader[domain.NameData](func() *container.NameContainer { return &container.NameContainer{} }, nil),
		newDefaultTableLoader[domain.RuneData](nil, map[string]string{"quality": "Quality"}),
		newDefaultTableLoader[domain.VipData](nil, nil),
	}

	for _, config := range tableConfigs {
		tableLoaders[config.tableName] = config
		if config.containerType != nil {
			containerKeys[config.containerType] = config.tableName
		}
	}
}

func InitConfig() {
	mgr := &DataManager{
		containers: make(map[string]any),
	}
	reader := data.NewExcelDataReader(true)

	tableNames := make([]string, 0, len(tableLoaders))
	for name := range tableLoaders {
		tableNames = append(tableNames, name)
	}
	sort.Strings(tableNames)

	for _, name := range tableNames {
		config := tableLoaders[name]
		container, err := config.load(reader)
		if err != nil {
			fmt.Printf("Failed to process table %s: %v\n", name, err)
			continue
		}
		mgr.containers[name] = container
	}

	global = mgr
}

func GetContainer(name string) any {
	return global.containers[name]
}

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

func QueryAll[V any]() []*V {
	tableName := getTableName[V]()
	container := GetContainer(tableName)
	if container == nil {
		return nil
	}
	if c, ok := container.(data.IContainer[int32, V]); ok {
		return c.GetAllRecords()
	}
	return nil
}

func QueryById[V any](id int32) *V {
	tableName := getTableName[V]()
	container := GetContainer(tableName)
	if container == nil {
		return nil
	}
	if c, ok := container.(data.IContainer[int32, V]); ok {
		return c.GetRecord(id)
	}
	return nil
}

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

func getTableName[V any]() string {
	t := reflect.TypeOf((*V)(nil)).Elem()
	name := strings.TrimSuffix(t.Name(), "Data")
	return strings.ToLower(name)
}

func newDefaultTableLoader[T any](containerType reflect.Type, indexFields map[string]string) configTableLoader {
	tableName := getTableName[T]()
	filePath := tableName + ".xlsx"
	return configTableLoader{
		tableName:     tableName,
		containerType: containerType,
		load: func(reader data.DataReader) (any, error) {
			return data.ProcessTableTyped(reader, tableName, filePath, func(record *T) int32 {
				return mustGetInt32Field(record, "Id")
			}, buildFieldIndexFuncs[T](indexFields))
		},
	}
}

func newContainerTableLoader[T any, C interface {
	data.IBaseContainer
	data.IDataInjector
}](newContainer func() C, indexFields map[string]string) configTableLoader {
	tableName := getTableName[T]()
	filePath := tableName + ".xlsx"
	return configTableLoader{
		tableName:     tableName,
		containerType: reflect.TypeOf(newContainer()),
		load: func(reader data.DataReader) (any, error) {
			records, err := data.ReadTyped[T](reader, filePath)
			if err != nil {
				return nil, err
			}
			container := newContainer()
			container.Init()
			if err := container.Inject(tableName, records, func(record *T) int32 {
				return mustGetInt32Field(record, "Id")
			}, buildFieldIndexFuncs[T](indexFields)); err != nil {
				logger.ErrorNoStack(err)
				// 配置错误，直接panic
				panic(err)
			}
			container.AfterLoad()
			logger.Info(fmt.Sprintf("Loaded table [%s] with %d records", tableName, len(records)))
			return container, nil
		},
	}
}

func buildFieldIndexFuncs[T any](indexFields map[string]string) map[string]func(*T) any {
	if len(indexFields) == 0 {
		return nil
	}
	indexFuncs := make(map[string]func(*T) any, len(indexFields))
	for indexName, fieldExpr := range indexFields {
		// 支持联合索引：value 用逗号分隔多字段名，例如 "TargetId,LevelStart"。
		// 单字段时返回字段值本身；多字段时返回 []any，由 data 层 indexKey 拼成 name@v1_v2。
		fieldNames := strings.Split(fieldExpr, ",")
		for i, f := range fieldNames {
			fieldNames[i] = strings.TrimSpace(f)
		}
		if len(fieldNames) == 1 {
			fieldName := fieldNames[0]
			indexFuncs[indexName] = func(record *T) any {
				return mustGetFieldValue(record, fieldName)
			}
		} else {
			names := fieldNames
			indexFuncs[indexName] = func(record *T) any {
				values := make([]any, len(names))
				for i, fieldName := range names {
					values[i] = mustGetFieldValue(record, fieldName)
				}
				return values
			}
		}
	}
	return indexFuncs
}

func mustGetInt32Field[T any](record *T, fieldName string) int32 {
	value := mustGetFieldValue(record, fieldName)
	resultValue := reflect.ValueOf(value)
	if !resultValue.IsValid() || !resultValue.Type().ConvertibleTo(reflect.TypeOf(int32(0))) {
		panic(fmt.Errorf("field %s value type %T cannot convert to int32", fieldName, value))
	}
	return int32(resultValue.Convert(reflect.TypeOf(int32(0))).Int())
}

func mustGetFieldValue[T any](record *T, fieldName string) any {
	recordValue := reflect.ValueOf(record)
	if !recordValue.IsValid() || recordValue.IsNil() {
		panic(fmt.Errorf("record is nil when accessing field %s", fieldName))
	}
	value := recordValue.Elem()
	field := value.FieldByName(fieldName)
	if !field.IsValid() {
		panic(fmt.Errorf("field %s not found in %v", fieldName, value.Type()))
	}
	return field.Interface()
}

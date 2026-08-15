package data

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// IBaseContainer 定义容器的基础接口，不包含泛型方法
type IBaseContainer interface {
	Init()
	AfterLoad()
}

// IDataInjector 定义容器注入能力，避免通过反射按方法名调用。
type IDataInjector interface {
	Inject(tableName string, records any, getIdFunc any, indexFuncs any) error
}

// IContainer 定义容器的泛型接口
type IContainer[K int32, V any] interface {
	IBaseContainer
	GetRecord(id int32) *V
	GetAllRecords() []*V
	GetRecordsByIndex(name string, index any) []*V
	GetUniqueRecordByIndex(name string, index any) *V
	// GetRecordsByJointIndex 按联合索引查询，indexes 为联合索引各字段值，顺序与声明时一致。
	GetRecordsByJointIndex(name string, indexes ...any) []*V
	// GetUniqueRecordByJointIndex 按联合索引查询单条记录。
	GetUniqueRecordByJointIndex(name string, indexes ...any) *V
}

// Container 是一个通用的数据容器，支持按 ID 查询、按索引查询和查询所有记录
type Container[K int32, V any] struct {
	data        map[K]*V        // 存储 ID 到记录指针的映射
	indexMapper map[string][]*V // 存储索引到记录指针的映射
}

// NewContainer 创建一个新的 Container 实例
func NewContainer[K int32, V any]() *Container[K, V] {
	return &Container[K, V]{
		data:        make(map[K]*V),
		indexMapper: make(map[string][]*V),
	}
}

// BuildContainer 用强类型方式构建容器，避免调用方手写反射注入逻辑。
func BuildContainer[V any](tableName string, records []*V, getID func(*V) int32, indexFuncs map[string]func(*V) any) (*Container[int32, V], error) {
	c := NewContainer[int32, V]()
	if err := c.Inject(tableName, records, getID, indexFuncs); err != nil {
		return nil, err
	}
	return c, nil
}

// Init 初始化容器，子类可以重写此方法
func (c *Container[int32, V]) Init() {
	if c.data == nil {
		c.data = make(map[int32]*V)
	}
	if c.indexMapper == nil {
		c.indexMapper = make(map[string][]*V)
	}
}

// AfterLoad 数据加载后的处理，子类可以重写此方法
func (c *Container[int32, V]) AfterLoad() {
}

// GetRecord 根据 ID 获取单个记录
func (c *Container[int32, V]) GetRecord(id int32) *V {
	record, exists := c.data[id]
	if !exists {
		return nil
	}
	return record
}

// GetAllRecords 获取所有记录
func (c *Container[int32, V]) GetAllRecords() []*V {
	ids := make([]int, 0, len(c.data))
	for id := range c.data {
		ids = append(ids, int(id))
	}
	sort.Ints(ids)

	records := make([]*V, 0, len(ids))
	for _, id := range ids {
		record := c.data[int32(id)]
		if record != nil {
			records = append(records, record)
		}
	}
	return records
}

// GetRecordsBy 根据索引名称和索引值获取记录
func (c *Container[int32, V]) GetRecordsByIndex(name string, index any) []*V {
	key := indexKey(name, index)
	records := c.indexMapper[key]
	ids := make([]int, 0, len(records))
	for _, record := range records {
		// record 是 *V（指向类型参数的指针），不是接口类型，
		// 需先转为 any 才能进行运行时类型断言。
		if ct, ok := any(record).(ConfigTable); ok {
			ids = append(ids, int(ct.GetId()))
		}
	}
	sort.Ints(ids)
	records = make([]*V, 0, len(ids))
	for _, id := range ids {
		record := c.data[int32(id)]
		if record != nil {
			records = append(records, record)
		}
	}
	return records
}

// GetUniqueRecordByIndex 按索引查询单条记录，索引无匹配或声明为非 unique 时返回首条。
func (c *Container[int32, V]) GetUniqueRecordByIndex(name string, index any) *V {
	records := c.GetRecordsByIndex(name, index)
	if len(records) == 0 {
		return nil
	}
	return records[0]
}

// GetRecordsByJointIndex 按联合索引查询，indexes 为联合索引各字段值，顺序与声明时一致。
func (c *Container[int32, V]) GetRecordsByJointIndex(name string, indexes ...any) []*V {
	return c.GetRecordsByIndex(name, any(indexes))
}

// GetUniqueRecordByJointIndex 按联合索引查询单条记录。
func (c *Container[int32, V]) GetUniqueRecordByJointIndex(name string, indexes ...any) *V {
	return c.GetUniqueRecordByIndex(name, any(indexes))
}

// Inject 将数据注入容器，并构建索引。
func (c *Container[int32, V]) Inject(tableName string, records any, getIdFunc any, indexFuncs any) error {
	// 确保 maps 已初始化
	if c.data == nil {
		c.data = make(map[int32]*V)
	}
	if c.indexMapper == nil {
		c.indexMapper = make(map[string][]*V)
	}

	// 获取记录切片的值
	recordsValue := reflect.ValueOf(records)
	if !recordsValue.IsValid() || recordsValue.Kind() != reflect.Slice {
		return fmt.Errorf("records must be a slice, got %T", records)
	}

	// 创建正确类型的记录切片
	var typedRecords []*V
	recordType := reflect.TypeOf((*V)(nil)).Elem()

	for i := 0; i < recordsValue.Len(); i++ {
		recordValue := recordsValue.Index(i)
		if recordValue.Kind() == reflect.Interface {
			if recordValue.IsNil() {
				return fmt.Errorf("record at index %d is nil", i)
			}
			recordValue = recordValue.Elem()
		}

		var ptr *V
		if recordValue.Type().AssignableTo(reflect.PtrTo(recordType)) {
			// 如果已经是正确类型的指针
			ptr = recordValue.Interface().(*V)
		} else if recordValue.Type().AssignableTo(recordType) {
			// 如果是正确类型的值
			newPtr := reflect.New(recordType)
			newPtr.Elem().Set(recordValue)
			ptr = newPtr.Interface().(*V)
		} else {
			return fmt.Errorf("record at index %d has incompatible type: got %v, want %v or %v", i, recordValue.Type(), recordType, reflect.PtrTo(recordType))
		}

		if ptr == nil {
			return fmt.Errorf("failed to create pointer for record at index %d", i)
		}
		typedRecords = append(typedRecords, ptr)
	}

	// 转换 ID 获取函数
	idFunc := reflect.ValueOf(getIdFunc)
	if !idFunc.IsValid() || idFunc.Kind() != reflect.Func {
		return fmt.Errorf("getIdFunc must be a function, got %T", getIdFunc)
	}
	if idFunc.Type().NumIn() != 1 || idFunc.Type().NumOut() != 1 {
		return fmt.Errorf("getIdFunc must accept exactly 1 argument and return 1 value, got %s", idFunc.Type())
	}
	idInputType := idFunc.Type().In(0)
	getTypedId := func(v *V) (int32, error) {
		arg, err := buildRecordCallArg(v, idInputType)
		if err != nil {
			return 0, err
		}
		results, err := callFuncSafely(idFunc, []reflect.Value{arg})
		if err != nil {
			return 0, err
		}
		result := results[0].Interface()

		// 使用反射进行类型转换
		resultValue := reflect.ValueOf(result)
		if !resultValue.Type().ConvertibleTo(reflect.TypeOf(*new(int32))) {
			return 0, fmt.Errorf("ID function returned %T which cannot be converted to type %T", result, *new(int32))
		}

		converted := resultValue.Convert(reflect.TypeOf(*new(int32))).Interface().(int32)
		return converted, nil
	}

	// 转换索引函数
	indexFuncsMap := make(map[string]func(*V) (any, error))
	if indexFuncs != nil {
		indexFuncsValue := reflect.ValueOf(indexFuncs)
		if indexFuncsValue.Kind() != reflect.Map {
			return fmt.Errorf("indexFuncs must be a map, got %T", indexFuncs)
		}
		iter := indexFuncsValue.MapRange()
		for iter.Next() {
			name := iter.Key().String()
			fn := iter.Value()
			if fn.Kind() != reflect.Func {
				return fmt.Errorf("index func %s must be a function, got %s", name, fn.Kind())
			}
			if fn.Type().NumIn() != 1 || fn.Type().NumOut() != 1 {
				return fmt.Errorf("index func %s must accept exactly 1 argument and return 1 value, got %s", name, fn.Type())
			}
			indexInputType := fn.Type().In(0)
			indexFuncsMap[name] = func(v *V) (any, error) {
				arg, err := buildRecordCallArg(v, indexInputType)
				if err != nil {
					return nil, err
				}
				results, err := callFuncSafely(fn, []reflect.Value{arg})
				if err != nil {
					return nil, err
				}
				return results[0].Interface(), nil
			}
		}
	}

	// 注入数据
	for i, record := range typedRecords {
		id, err := getTypedId(record)
		if err != nil {
			return fmt.Errorf("failed to build record id at index %d: %w", i, err)
		}
		if c.data[id] != nil {
			return fmt.Errorf("[config] table %s record with ID %d already exists", tableName, id)
		}
		// 构建索引
		c.data[id] = record
		// 构建索引
		for name, indexFunc := range indexFuncsMap {
			indexValue, err := indexFunc(record)
			if err != nil {
				return fmt.Errorf("[config] table %s failed to build index %s for record ID %d: %w", tableName, name, id, err)
			}
			key := indexKey(name, indexValue)
			c.indexMapper[key] = append(c.indexMapper[key], record)
		}
	}
	c.AfterLoad()
	return nil
}

// indexKey 生成索引键。index 为 []any 时按联合索引拼接为 name@v1_v2_...
func indexKey(name string, index any) string {
	return fmt.Sprintf("%s@%s", name, formatIndexValue(index))
}

// formatIndexValue 将索引值格式化为索引键中的值部分。
// 单值直接 %v；[]any（联合索引）用 "_" 连接各分量的 %v 表示。
func formatIndexValue(index any) string {
	switch v := index.(type) {
	case nil:
		return ""
	case []any:
		parts := make([]string, len(v))
		for i, item := range v {
			parts[i] = fmt.Sprintf("%v", item)
		}
		return strings.Join(parts, "_")
	default:
		return fmt.Sprintf("%v", v)
	}
}

func buildRecordCallArg[V any](record *V, inputType reflect.Type) (reflect.Value, error) {
	recordValue := reflect.ValueOf(record)
	value := recordValue.Elem()
	switch {
	case value.Type().AssignableTo(inputType):
		return value, nil
	case recordValue.Type().AssignableTo(inputType):
		return recordValue, nil
	case value.Kind() == reflect.Interface && !value.IsNil() && value.Elem().Type().AssignableTo(inputType):
		return value.Elem(), nil
	case recordValue.Type().ConvertibleTo(inputType):
		return recordValue.Convert(inputType), nil
	case value.Type().ConvertibleTo(inputType):
		return value.Convert(inputType), nil
	default:
		return reflect.Value{}, fmt.Errorf("record type %v is not assignable to function input %v", recordValue.Type(), inputType)
	}
}

func callFuncSafely(fn reflect.Value, args []reflect.Value) (results []reflect.Value, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("function call panic: %v", recovered)
		}
	}()
	return fn.Call(args), nil
}

package data

import (
	"fmt"
	"reflect"
)

// Container 是一个通用的数据容器，支持按 ID 查询、按索引查询和查询所有记录
type Container[K comparable, V any] struct {
	data        map[K]*V        // 存储 ID 到记录指针的映射
	indexMapper map[string][]*V // 存储索引到记录指针的映射
}

// IContainer 定义容器的接口
type IContainer interface {
	AfterLoad()
	Init()
}

// NewContainer 创建一个新的 Container 实例
func NewContainer[K comparable, V any]() *Container[K, V] {
	return &Container[K, V]{
		data:        make(map[K]*V),
		indexMapper: make(map[string][]*V),
	}
}

// Init 初始化容器，子类可以重写此方法
func (c *Container[K, V]) Init() {
	if c.data == nil {
		c.data = make(map[K]*V)
	}
	if c.indexMapper == nil {
		c.indexMapper = make(map[string][]*V)
	}
}

// AfterLoad 数据加载后的处理，子类可以重写此方法
func (c *Container[K, V]) AfterLoad() {
}

// GetRecord 根据 ID 获取单个记录
func (c *Container[K, V]) GetRecord(id K) (*V, bool) {
	record, exists := c.data[id]
	return record, exists
}

// GetAllRecords 获取所有记录
func (c *Container[K, V]) GetAllRecords() []*V {
	records := make([]*V, 0, len(c.data))
	for _, record := range c.data {
		records = append(records, record)
	}
	return records
}

// GetRecordsBy 根据索引名称和索引值获取记录
func (c *Container[K, V]) GetRecordsBy(name string, index any) []*V {
	key := indexKey(name, index)
	return c.indexMapper[key]
}

// Values 返回所有记录的map
func (c *Container[K, V]) Values() map[K]*V {
	return c.data
}

// Inject 将数据注入容器，并构建索引
func (c *Container[K, V]) Inject(records interface{}, getIdFunc interface{}, indexFuncs interface{}) {
	// 确保 maps 已初始化
	if c.data == nil {
		c.data = make(map[K]*V)
	}
	if c.indexMapper == nil {
		c.indexMapper = make(map[string][]*V)
	}

	// 获取记录切片的值
	recordsValue := reflect.ValueOf(records)
	if recordsValue.Kind() != reflect.Slice {
		panic("records must be a slice")
	}

	// 创建正确类型的记录切片
	var typedRecords []*V
	recordType := reflect.TypeOf((*V)(nil)).Elem()

	for i := 0; i < recordsValue.Len(); i++ {
		recordValue := recordsValue.Index(i)
		if recordValue.Kind() == reflect.Interface {
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
			panic(fmt.Sprintf("record at index %d has incompatible type: got %v, want %v", i, recordValue.Type(), recordType))
		}

		if ptr == nil {
			panic(fmt.Sprintf("failed to create pointer for record at index %d", i))
		}
		typedRecords = append(typedRecords, ptr)
	}

	// 转换 ID 获取函数
	idFunc := reflect.ValueOf(getIdFunc)
	getTypedId := func(v *V) K {
		// 确保传递给 ID 函数的是解引用后的值
		val := reflect.ValueOf(v).Elem().Interface()
		results := idFunc.Call([]reflect.Value{reflect.ValueOf(val)})
		result := results[0].Interface()

		// 使用反射进行类型转换
		resultValue := reflect.ValueOf(result)
		if !resultValue.Type().ConvertibleTo(reflect.TypeOf(*new(K))) {
			panic(fmt.Sprintf("ID function returned %T which cannot be converted to type %T", result, *new(K)))
		}

		converted := resultValue.Convert(reflect.TypeOf(*new(K))).Interface().(K)
		return converted
	}

	// 转换索引函数
	indexFuncsMap := make(map[string]func(*V) any)
	if indexFuncs != nil {
		indexFuncsValue := reflect.ValueOf(indexFuncs)
		iter := indexFuncsValue.MapRange()
		for iter.Next() {
			name := iter.Key().String()
			fn := iter.Value()
			indexFuncsMap[name] = func(v *V) any {
				// 确保传递给索引函数的是解引用后的值
				val := reflect.ValueOf(v).Elem().Interface()
				results := fn.Call([]reflect.Value{reflect.ValueOf(val)})
				return results[0].Interface()
			}
		}
	}

	// 注入数据
	for _, record := range typedRecords {
		id := getTypedId(record)
		c.data[id] = record

		// 构建索引
		for name, indexFunc := range indexFuncsMap {
			indexValue := indexFunc(record)
			key := indexKey(name, indexValue)
			c.indexMapper[key] = append(c.indexMapper[key], record)
		}
	}
	c.AfterLoad()
}

// indexKey 生成索引键
func indexKey(name string, index any) string {
	return fmt.Sprintf("%s@%v", name, index)
}

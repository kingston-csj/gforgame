package data

import "fmt"

// Container 是一个通用的数据容器，支持按 ID 查询、按索引查询和查询所有记录
type Container[K comparable, V any] struct {
	data        map[K]V        // 存储 ID 到记录的映射
	indexMapper map[string][]V // 存储索引到记录的映射
}

// NewContainer 创建一个新的 Container 实例
func NewContainer[K comparable, V any]() *Container[K, V] {
	return &Container[K, V]{
		data:        make(map[K]V),
		indexMapper: make(map[string][]V),
	}
}

// Inject 将数据注入容器，并构建索引
func (c *Container[K, V]) Inject(records []V, getIdFunc func(V) K, indexFuncs map[string]func(V) interface{}) {
	for _, record := range records {
		id := getIdFunc(record)
		c.data[id] = record

		// 构建索引
		for name, indexFunc := range indexFuncs {
			indexValue := indexFunc(record)
			key := indexKey(name, indexValue)
			c.indexMapper[key] = append(c.indexMapper[key], record)
		}
	}
}

// GetRecord 根据 ID 获取单个记录
func (c *Container[K, V]) GetRecord(id K) (V, bool) {
	record, exists := c.data[id]
	return record, exists
}

// GetAllRecords 获取所有记录
func (c *Container[K, V]) GetAllRecords() []V {
	records := make([]V, 0, len(c.data))
	for _, record := range c.data {
		records = append(records, record)
	}
	return records
}

// GetRecordsBy 根据索引名称和索引值获取记录
func (c *Container[K, V]) GetRecordsBy(name string, index any) []V {
	key := indexKey(name, index)
	return c.indexMapper[key]
}

// indexKey 生成索引键
func indexKey(name string, index interface{}) string {
	return fmt.Sprintf("%s@%v", name, index)
}

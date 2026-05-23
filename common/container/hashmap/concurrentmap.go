package hashmap

import (
	"fmt"
	"hash/fnv"
	"sync"
)

// 默认分片数量，可以根据实际业务和 CPU 核心数调整
const DefaultShardCount = 32

// ConcurrentMap 对外暴露的并发安全 Map
type ConcurrentMap[K comparable, V any] struct {
	shards []*ConcurrentMapShared[K, V] // 分片数组
	count  int                          // 分片数量
}

// ConcurrentMapShared 每一个具体的分片，包含一个小 map 和一个读写锁
type ConcurrentMapShared[K comparable, V any] struct {
	items map[K]V
	mu    sync.RWMutex
}

// NewConcurrentMap 创建一个新的并发 Map
func NewConcurrentMap[K comparable, V any]() *ConcurrentMap[K, V] {
	m := &ConcurrentMap[K, V]{
		shards: make([]*ConcurrentMapShared[K, V], DefaultShardCount),
		count:  DefaultShardCount,
	}
	for i := 0; i < DefaultShardCount; i++ {
		m.shards[i] = &ConcurrentMapShared[K, V]{
			items: make(map[K]V),
		}
	}
	return m
}

// getShard 根据 key 计算出对应的分片索引
func (m *ConcurrentMap[K, V]) getShard(key K) *ConcurrentMapShared[K, V] {
	// 使用 FNV-1a 哈希算法计算 hash 值
	h := fnv.New32a()
	// 这里利用了 fmt 的格式化将任意可比较的 key 转为字节切片进行哈希
	fmt.Fprint(h, key)
	hashVal := h.Sum32()
	return m.shards[uint(hashVal)%uint(m.count)]
}

// Set 插入或更新键值对
func (m *ConcurrentMap[K, V]) Set(key K, value V) {
	shard := m.getShard(key)
	shard.mu.Lock()         // 写操作加写锁
	defer shard.mu.Unlock()
	shard.items[key] = value
}

// Get 获取指定 key 的值
func (m *ConcurrentMap[K, V]) Get(key K) (V, bool) {
	shard := m.getShard(key)
	shard.mu.RLock()        // 读操作加读锁
	defer shard.mu.RUnlock()
	val, ok := shard.items[key]
	return val, ok
}

// Remove 删除指定 key
func (m *ConcurrentMap[K, V]) Remove(key K) {
	shard := m.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	delete(shard.items, key)
}

// Count 获取 Map 中的元素总数（需要遍历所有分片）
func (m *ConcurrentMap[K, V]) Count() int {
	count := 0
	for i := 0; i < m.count; i++ {
		shard := m.shards[i]
		shard.mu.RLock()
		count += len(shard.items)
		shard.mu.RUnlock()
	}
	return count
}

func (m *ConcurrentMap[K, V]) Values() []V {
	values := make([]V, 0)
	for i := 0; i < m.count; i++ {
		shard := m.shards[i]
		shard.mu.RLock()
		for _, v := range shard.items {
			values = append(values, v)
		}
		shard.mu.RUnlock()
	}
	return values
}

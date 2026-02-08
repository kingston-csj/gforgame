package hashmap

import (
	"errors"
	"fmt"
)

// DualHashMap 双向哈希映射，K和V需为可比较类型（满足map键要求）
// 保证K和V一一对应，无重复的K或V
type DualHashMap[K comparable, V comparable] struct {
	keyToVal map[K]V // key到value的映射
	valToKey map[V]K // value到key的映射
}

// NewDualHashMap 创建一个空的双向哈希映射
func NewDualHashMap[K comparable, V comparable]() *DualHashMap[K, V] {
	return &DualHashMap[K, V]{
		keyToVal: make(map[K]V),
		valToKey: make(map[V]K),
	}
}

// Put 添加key-value映射，若key或value已存在则返回错误
func (d *DualHashMap[K, V]) Put(key K, value V) error {
	// 检查key是否已存在
	if _, exists := d.keyToVal[key]; exists {
		return errors.New(fmt.Sprintf("key %v 已存在", key))
	}
	// 检查value是否已存在（保证一一对应）
	if _, exists := d.valToKey[value]; exists {
		return errors.New(fmt.Sprintf("value %v 已存在", value))
	}
	// 同时更新两个map
	d.keyToVal[key] = value
	d.valToKey[value] = key
	return nil
}

// GetByKey 通过key查找value，返回value和是否存在的布尔值
func (d *DualHashMap[K, V]) GetByKey(key K) (V, bool) {
	val, exists := d.keyToVal[key]
	return val, exists
}

// GetByValue 通过value查找key，返回key和是否存在的布尔值
func (d *DualHashMap[K, V]) GetByValue(value V) (K, bool) {
	key, exists := d.valToKey[value]
	return key, exists
}

// DeleteByKey 通过key删除映射（同时删除value对应的反向映射）
func (d *DualHashMap[K, V]) DeleteByKey(key K) {
	if val, exists := d.keyToVal[key]; exists {
		delete(d.keyToVal, key)
		delete(d.valToKey, val)
	}
}

// DeleteByValue 通过value删除映射（同时删除key对应的正向映射）
func (d *DualHashMap[K, V]) DeleteByValue(value V) {
	if key, exists := d.valToKey[value]; exists {
		delete(d.valToKey, value)
		delete(d.keyToVal, key)
	}
}

// Len 返回映射的元素个数
func (d *DualHashMap[K, V]) Len() int {
	// 两个map长度始终一致，取其一即可
	return len(d.keyToVal)
}

// Clear 清空所有映射
func (d *DualHashMap[K, V]) Clear() {
	d.keyToVal = make(map[K]V)
	d.valToKey = make(map[V]K)
}
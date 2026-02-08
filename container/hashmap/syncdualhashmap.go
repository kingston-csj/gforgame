package hashmap

import "sync"

// ---------------------- 线程安全版（并发场景使用） ----------------------
// SyncDualHashMap 线程安全的双向哈希映射
type SyncDualHashMap[K comparable, V comparable] struct {
	dualMap *DualHashMap[K, V]
	mu      sync.RWMutex // 读写锁，读多写少场景更高效
}

// NewSyncDualHashMap 创建线程安全的双向哈希映射
func NewSyncDualHashMap[K comparable, V comparable]() *SyncDualHashMap[K, V] {
	return &SyncDualHashMap[K, V]{
		dualMap: NewDualHashMap[K, V](),
	}
}

// Put 线程安全的添加操作
func (s *SyncDualHashMap[K, V]) Put(key K, value V) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.dualMap.Put(key, value)
}

// GetByKey 线程安全的按key查询
func (s *SyncDualHashMap[K, V]) GetByKey(key K) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.dualMap.GetByKey(key)
}

// GetByValue 线程安全的按value查询
func (s *SyncDualHashMap[K, V]) GetByValue(value V) (K, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.dualMap.GetByValue(value)
}

// DeleteByKey 线程安全的按key删除
func (s *SyncDualHashMap[K, V]) DeleteByKey(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dualMap.DeleteByKey(key)
}

// DeleteByValue 线程安全的按value删除
func (s *SyncDualHashMap[K, V]) DeleteByValue(value V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dualMap.DeleteByValue(value)
}

// Len 线程安全的获取长度
func (s *SyncDualHashMap[K, V]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.dualMap.Len()
}

// Clear 线程安全的清空操作
func (s *SyncDualHashMap[K, V]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dualMap.Clear()
}
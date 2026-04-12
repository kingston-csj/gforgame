package set

import (
	"encoding/json"
)

// Set 是一个基于泛型的集合实现，非线程安全
// 支持 JSON 序列化为数组，反序列化时自动去重
type Set[T comparable] struct {
	m map[T]struct{}
}

// NewSet 创建一个新的 Set
func NewSet[T comparable](items ...T) *Set[T] {
	s := &Set[T]{
		m: make(map[T]struct{}),
	}
	for _, item := range items {
		s.Add(item)
	}
	return s
}

// Add 添加元素，如果元素已存在返回 false，否则返回 true
func (s *Set[T]) Add(item T) bool {
	if s.m == nil {
		s.m = make(map[T]struct{})
	}
	if _, ok := s.m[item]; ok {
		return false
	}
	s.m[item] = struct{}{}
	return true
}

// Remove 移除元素
func (s *Set[T]) Remove(item T) {
	if s.m == nil {
		return
	}
	delete(s.m, item)
}

// Contains 是否包含元素
func (s *Set[T]) Contains(item T) bool {
	if s.m == nil {
		return false
	}
	_, ok := s.m[item]
	return ok
}

// Len 返回元素个数
func (s *Set[T]) Len() int {
	if s.m == nil {
		return 0
	}
	return len(s.m)
}

// ToSlice 转换为切片
func (s *Set[T]) ToSlice() []T {
	if s.m == nil {
		return []T{}
	}
	slice := make([]T, 0, len(s.m))
	for item := range s.m {
		slice = append(slice, item)
	}
	return slice
}

// MarshalJSON 实现 json.Marshaler 接口
func (s *Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.ToSlice())
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (s *Set[T]) UnmarshalJSON(data []byte) error {
	var slice []T
	if err := json.Unmarshal(data, &slice); err != nil {
		return err
	}
	s.m = make(map[T]struct{}, len(slice))
	for _, item := range slice {
		s.m[item] = struct{}{}
	}
	return nil
}

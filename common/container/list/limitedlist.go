package list

import (
	"container/list"
	"sync"
)

// LimitedList 带最大长度限制的泛型链表，超过长度时自动删除最早元素。
// 所有公共方法均为并发安全。
type LimitedList[T any] struct {
	mu      sync.RWMutex
	list    *list.List
	maxSize int // 最大允许的元素数量
}

// NewLimitedList 创建一个带长度限制的泛型链表
func NewLimitedList[T any](maxSize int) *LimitedList[T] {
	return &LimitedList[T]{
		list:    list.New(),
		maxSize: maxSize,
	}
}

// Push 添加元素到链表尾部，若超过最大长度则删除头部（最早的元素）
func (l *LimitedList[T]) Push(v T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.list.PushBack(v)
	// 超过最大长度时，移除头部元素
	if l.list.Len() > l.maxSize {
		l.list.Remove(l.list.Front())
	}
}

// Len 返回当前元素数量
func (l *LimitedList[T]) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.list.Len()
}

// Front 返回头部元素（最早添加的），若链表为空则返回零值和false
func (l *LimitedList[T]) Front() (T, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if elem := l.list.Front(); elem != nil {
		return elem.Value.(T), true
	}
	var zero T
	return zero, false
}

// PopFront 弹出头部元素（最早添加的），若链表为空则返回零值和false
func (l *LimitedList[T]) PopFront() (T, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if elem := l.list.Front(); elem != nil {
		l.list.Remove(elem)
		return elem.Value.(T), true
	}
	var zero T
	return zero, false
}

// Back 返回尾部元素（最新添加的），若链表为空则返回零值和false
func (l *LimitedList[T]) Back() (T, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if elem := l.list.Back(); elem != nil {
		return elem.Value.(T), true
	}
	var zero T
	return zero, false
}

// Each 遍历所有元素（从早到晚）。fn 在锁保护下执行，fn 内不要再次调用本链表方法（否则死锁）。
func (l *LimitedList[T]) Each(fn func(T)) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	for elem := l.list.Front(); elem != nil; elem = elem.Next() {
		fn(elem.Value.(T))
	}
}

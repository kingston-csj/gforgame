package list

import (
	"container/list"
)

// LimitedList 带最大长度限制的泛型链表，超过长度时自动删除最早元素
type LimitedList[T any] struct {
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
	l.list.PushBack(v)
	// 超过最大长度时，移除头部元素
	if l.list.Len() > l.maxSize {
		l.list.Remove(l.list.Front())
	}
}

// Len 返回当前元素数量
func (l *LimitedList[T]) Len() int {
	return l.list.Len()
}

// Front 返回头部元素（最早添加的），若链表为空则返回零值和false
func (l *LimitedList[T]) Front() (T, bool) {
	if elem := l.list.Front(); elem != nil {
		return elem.Value.(T), true
	}
	var zero T
	return zero, false
}

// Back 返回尾部元素（最新添加的），若链表为空则返回零值和false
func (l *LimitedList[T]) Back() (T, bool) {
	if elem := l.list.Back(); elem != nil {
		return elem.Value.(T), true
	}
	var zero T
	return zero, false
}

// Each 遍历所有元素（从早到晚）
func (l *LimitedList[T]) Each(fn func(T)) {
	for elem := l.list.Front(); elem != nil; elem = elem.Next() {
		fn(elem.Value.(T))
	}
}

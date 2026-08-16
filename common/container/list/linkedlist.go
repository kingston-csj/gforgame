package list

import (
	"errors"
)

// 对标 Java LinkedList
// 双向链表，非线程安全
type LinkedList[T any] struct {
	first *Node[T] // 头节点
	last  *Node[T] // 尾节点
	size  int      // 元素数量
}

// 链表节点
type Node[T any] struct {
	item T
	prev *Node[T]
	next *Node[T]
}

// NewLinkedList 创建空链表
func NewLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{}
}

// Add 添加元素到尾部（对标 Java add）
func (l *LinkedList[T]) Add(t T) {
	l.linkLast(t)
}

// AddFirst 添加到头部
func (l *LinkedList[T]) AddFirst(t T) {
	l.linkFirst(t)
}

// AddLast 添加到尾部
func (l *LinkedList[T]) AddLast(t T) {
	l.linkLast(t)
}

// RemoveFirst 删除头部并返回
func (l *LinkedList[T]) RemoveFirst() (T, error) {
	if l.size == 0 {
		var zero T
		return zero, errors.New("no such element")
	}
	return l.unlinkFirst(), nil
}

// RemoveLast 删除尾部并返回
func (l *LinkedList[T]) RemoveLast() (T, error) {
	if l.size == 0 {
		var zero T
		return zero, errors.New("no such element")
	}
	return l.unlinkLast(), nil
}

// Get 获取指定下标元素
func (l *LinkedList[T]) Get(index int) (T, error) {
	if !l.isIndexValid(index) {
		var zero T
		return zero, errors.New("index out of bounds")
	}
	return l.node(index).item, nil
}

// RemoveAt 删除指定下标
func (l *LinkedList[T]) RemoveAt(index int) (T, error) {
	if !l.isIndexValid(index) {
		var zero T
		return zero, errors.New("index out of bounds")
	}
	return l.unlink(l.node(index)), nil
}

// Size 返回大小
func (l *LinkedList[T]) Size() int {
	return l.size
}

// IsEmpty 是否为空
func (l *LinkedList[T]) IsEmpty() bool {
	return l.size == 0
}

// Clear 清空链表
func (l *LinkedList[T]) Clear() {
	// 帮助 GC
	for n := l.first; n != nil; {
		next := n.next
		var zero T
		n.item = zero
		n.prev, n.next = nil, nil
		n = next
	}
	l.first, l.last = nil, nil
	l.size = 0
}

// 内部：连接到头部
func (l *LinkedList[T]) linkFirst(t T) {
	f := l.first
	newNode := &Node[T]{item: t, next: f}
	l.first = newNode
	if f == nil {
		l.last = newNode
	} else {
		f.prev = newNode
	}
	l.size++
}

// 内部：连接到尾部
func (l *LinkedList[T]) linkLast(t T) {
	last := l.last
	newNode := &Node[T]{item: t, prev: last}
	l.last = newNode
	if last == nil {
		l.first = newNode
	} else {
		last.next = newNode
	}
	l.size++
}

// 内部：解除头节点
func (l *LinkedList[T]) unlinkFirst() T {
	n := l.first
	val := n.item
	next := n.next
	var zero T
	n.item = zero
	n.next = nil // GC

	l.first = next
	if next == nil {
		l.last = nil
	} else {
		next.prev = nil
	}
	l.size--
	return val
}

// 内部：解除尾节点
func (l *LinkedList[T]) unlinkLast() T {
	n := l.last
	val := n.item
	prev := n.prev
	var zero T
	n.item = zero
	n.prev = nil // GC

	l.last = prev
	if prev == nil {
		l.first = nil
	} else {
		prev.next = nil
	}
	l.size--
	return val
}

// 内部：解除任意节点
func (l *LinkedList[T]) unlink(n *Node[T]) T {
	val := n.item
	prev := n.prev
	next := n.next

	if prev == nil {
		l.first = next
	} else {
		prev.next = next
		n.prev = nil
	}

	if next == nil {
		l.last = prev
	} else {
		next.prev = prev
		n.next = nil
	}

	var zero T
	n.item = zero
	l.size--
	return val
}

// 内部：找到指定下标的节点（前后半分治查找）
func (l *LinkedList[T]) node(index int) *Node[T] {
	if index < l.size/2 {
		n := l.first
		for i := 0; i < index; i++ {
			n = n.next
		}
		return n
	} else {
		n := l.last
		for i := l.size - 1; i > index; i-- {
			n = n.prev
		}
		return n
	}
}

// 内部：下标是否合法
func (l *LinkedList[T]) isIndexValid(index int) bool {
	return index >= 0 && index < l.size
}

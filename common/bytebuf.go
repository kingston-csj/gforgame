package common

import (
	"errors"
	"fmt"
)

// 自带的ByteBuf，功能实在是太局限了，这里提供一个高级点的
// ByteBuffer 自定义字节缓冲区（支持标记/重置、动态扩容、读索引回退）
type ByteBuffer struct {
	buffer            []byte // 内部存储字节数据的缓冲区
	writeIndex        int    // 下一次写入的起始位置（从0开始）
	readIndex         int    // 下一次读取的起始位置（从0开始）
	markedReadIndex   int    // 标记的读取位置（用于 MarkRead/ResetRead）
	markedWriteIndex  int    // 标记的写入位置（用于 MarkWrite/ResetWrite）
	initialCapacity   int    // 初始容量（用于 Clear 时恢复初始大小）
	maxCapacity       int    // 最大容量（避免无限制扩容，0 表示无限制）
}

// 预定义错误
var (
	ErrBufferOverflow = errors.New("buffer overflow: exceed max capacity")
	ErrReadOutOfRange = errors.New("read out of range: no enough data")
	ErrInvalidIndex   = errors.New("invalid index: readIndex > writeIndex")
)

// NewByteBuffer 创建 ByteBuffer 实例
// initialCapacity: 初始容量（建议根据业务设置，如 4096）
// maxCapacity: 最大容量（0 表示无限制）
func NewByteBuffer(initialCapacity, maxCapacity int) *ByteBuffer {
	if initialCapacity <= 0 {
		initialCapacity = 4096 // 默认初始容量
	}
	return &ByteBuffer{
		buffer:          make([]byte, initialCapacity),
		initialCapacity: initialCapacity,
		maxCapacity:     maxCapacity,
		// 初始时 readIndex、writeIndex、标记位置均为 0
	}
}

// -------------------------- 核心写入操作 --------------------------
// Write 写入字节切片到缓冲区（自动扩容）
func (b *ByteBuffer) Write(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	// 检查是否需要扩容
	requiredCapacity := b.writeIndex + len(data)
	if requiredCapacity > len(b.buffer) {
		if err := b.expand(requiredCapacity); err != nil {
			return err
		}
	}

	// 拷贝数据到缓冲区
	copy(b.buffer[b.writeIndex:], data)
	b.writeIndex += len(data)
	return nil
}

// WriteByte 写入单个字节
func (b *ByteBuffer) WriteByte(c byte) error {
	return b.Write([]byte{c})
}

// -------------------------- 核心读取操作 --------------------------
// Read 读取指定长度的字节到目标切片（会移动 readIndex）
// 返回实际读取的字节数和错误
func (b *ByteBuffer) Read(dst []byte) (int, error) {
	if len(dst) == 0 {
		return 0, nil
	}

	// 检查可读字节数
	available := b.writeIndex - b.readIndex
	if available <= 0 {
		return 0, ErrReadOutOfRange
	}

	// 实际读取长度（取目标长度和可用字节数的最小值）
	readLen := len(dst)
	if readLen > available {
		readLen = available
		dst = dst[:readLen] // 截取目标切片，避免越界
	}

	// 拷贝数据到目标切片
	copy(dst, b.buffer[b.readIndex:b.readIndex+readLen])
	b.readIndex += readLen
	return readLen, nil
}

// Next 读取指定长度的字节切片（会移动 readIndex）
// 注意：返回的是缓冲区的拷贝，避免外部修改内部数据
func (b *ByteBuffer) Next(n int) ([]byte, error) {
	if n <= 0 {
		return []byte{}, nil
	}

	available := b.writeIndex - b.readIndex
	if n > available {
		return nil, fmt.Errorf("%w: need %d, available %d", ErrReadOutOfRange, n, available)
	}

	// 拷贝数据（避免外部修改内部 buffer）
	result := make([]byte, n)
	copy(result, b.buffer[b.readIndex:b.readIndex+n])
	b.readIndex += n
	return result, nil
}

// Peek 偷看指定长度的字节（不移动 readIndex）
func (b *ByteBuffer) Peek(n int) ([]byte, error) {
	if n <= 0 {
		return []byte{}, nil
	}

	available := b.writeIndex - b.readIndex
	if n > available {
		return nil, fmt.Errorf("%w: need %d, available %d", ErrReadOutOfRange, n, available)
	}

	// 返回拷贝，不移动 readIndex
	result := make([]byte, n)
	copy(result, b.buffer[b.readIndex:b.readIndex+n])
	return result, nil
}

// -------------------------- 标记/重置操作（你指定的核心功能） --------------------------
// MarkReadIndex 标记当前读取位置（保存到 markedReadIndex）
func (b *ByteBuffer) MarkReadIndex() {
	b.markedReadIndex = b.readIndex
}

// ResetReadIndex 重置读取位置到标记的位置
func (b *ByteBuffer) ResetReadIndex() error {
	if b.markedReadIndex < 0 || b.markedReadIndex > b.writeIndex {
		return ErrInvalidIndex
	}
	b.readIndex = b.markedReadIndex
	return nil
}

// MarkWriteIndex 标记当前写入位置（保存到 markedWriteIndex）
func (b *ByteBuffer) MarkWriteIndex() {
	b.markedWriteIndex = b.writeIndex
}

// ResetWriteIndex 重置写入位置到标记的位置
func (b *ByteBuffer) ResetWriteIndex() error {
	if b.markedWriteIndex < b.readIndex || b.markedWriteIndex > len(b.buffer) {
		return ErrInvalidIndex
	}
	b.writeIndex = b.markedWriteIndex
	return nil
}

// -------------------------- 索引操作（协议解析关键） --------------------------
// UnreadBytes 回退读索引 n 个字节（解析 header 后包不完整时使用）
func (b *ByteBuffer) UnreadBytes(n int) error {
	if n <= 0 {
		return nil
	}
	if b.readIndex-n < 0 {
		return ErrInvalidIndex
	}
	b.readIndex -= n
	return nil
}

// SetReadIndex 直接设置读索引（慎用，需确保索引合法）
func (b *ByteBuffer) SetReadIndex(index int) error {
	if index < 0 || index > b.writeIndex {
		return ErrInvalidIndex
	}
	b.readIndex = index
	return nil
}

// SetWriteIndex 直接设置写索引（慎用）
func (b *ByteBuffer) SetWriteIndex(index int) error {
	if index < b.readIndex || index > len(b.buffer) {
		return ErrInvalidIndex
	}
	b.writeIndex = index
	return nil
}

// -------------------------- 缓冲区状态查询 --------------------------
// Len 未读字节数（writeIndex - readIndex）
func (b *ByteBuffer) Len() int {
	return b.writeIndex - b.readIndex
}

// Capacity 缓冲区总容量
func (b *ByteBuffer) Capacity() int {
	return len(b.buffer)
}

// RemainingWrite 剩余可写字节数
func (b *ByteBuffer) RemainingWrite() int {
	return len(b.buffer) - b.writeIndex
}

// IsEmpty 是否无未读数据
func (b *ByteBuffer) IsEmpty() bool {
	return b.Len() == 0
}

// -------------------------- 缓冲区管理 --------------------------
// Clear 清空缓冲区（重置索引，不清除数据，复用缓冲区）
func (b *ByteBuffer) Clear() {
	b.readIndex = 0
	b.writeIndex = 0
	b.markedReadIndex = 0
	b.markedWriteIndex = 0
	// 缩容到初始容量（避免长期占用过大内存）
	if len(b.buffer) > b.initialCapacity {
		b.buffer = make([]byte, b.initialCapacity)
	}
}

// Compact 压缩缓冲区（将未读数据移到缓冲区开头，释放尾部空间）
func (b *ByteBuffer) Compact() {
	if b.readIndex == 0 {
		return // 无需压缩
	}

	// 把未读数据移到开头
	copy(b.buffer[0:], b.buffer[b.readIndex:b.writeIndex])
	// 更新索引
	b.writeIndex -= b.readIndex
	b.readIndex = 0
	// 重置标记位置
	if b.markedReadIndex > 0 {
		b.markedReadIndex = 0
	}
	if b.markedWriteIndex > b.writeIndex {
		b.markedWriteIndex = b.writeIndex
	}
}

// -------------------------- 内部扩容逻辑 --------------------------
func (b *ByteBuffer) expand(requiredCapacity int) error {
	// 检查是否超过最大容量
	if b.maxCapacity > 0 && requiredCapacity > b.maxCapacity {
		return fmt.Errorf("%w: required %d, max %d", ErrBufferOverflow, requiredCapacity, b.maxCapacity)
	}

	// 扩容策略：当前容量翻倍，直到满足需求（至少为 requiredCapacity）
	newCapacity := len(b.buffer)
	for newCapacity < requiredCapacity {
		newCapacity *= 2
		// 避免扩容过大（如果有最大容量限制）
		if b.maxCapacity > 0 && newCapacity > b.maxCapacity {
			newCapacity = b.maxCapacity
			break
		}
	}

	// 创建新缓冲区并拷贝原数据
	newBuffer := make([]byte, newCapacity)
	copy(newBuffer[0:b.writeIndex], b.buffer[0:b.writeIndex])
	b.buffer = newBuffer
	return nil
}

// String 调试用：打印缓冲区状态（不打印具体字节数据）
func (b *ByteBuffer) String() string {
	return fmt.Sprintf(
		"ByteBuffer[read=%d, write=%d, len=%d, capacity=%d, max=%d]",
		b.readIndex, b.writeIndex, b.Len(), len(b.buffer), b.maxCapacity,
	)
}
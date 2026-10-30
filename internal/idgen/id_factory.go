package idgen

import (
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	serverconfig "github.com/forfun/gforgame/internal/serverconfig"
)

// IDGenerator 全局唯一ID生成器
// 高16位为serverId
// 中32位为系统秒数
// 低16位为自增长号
type IDGenerator struct {
	serverID int64
	sequence atomic.Int64
}

// NewIDGenerator 创建一个新的ID生成器实例
func NewIDGenerator(serverID int64) *IDGenerator {
	return &IDGenerator{
		serverID: serverID,
	}
}

// NextID 生成下一个全局唯一id
func (g *IDGenerator) NextID() string {
	currentTimeSeconds := time.Now().Unix()
	id := (g.serverID << 48) |
		(currentTimeSeconds << 16) |
		g.sequence.Add(1)&0xFFFF

	return strconv.FormatInt(id, 10)
}

// DefaultIDGenerator 默认的全局ID生成器实例（懒初始化，首次使用时创建）
var (
	defaultIDGen     *IDGenerator
	defaultIDGenOnce sync.Once
)

// GetNextID 使用默认生成器生成ID的便捷方法
func GetNextID() string {
	defaultIDGenOnce.Do(func() {
		defaultIDGen = NewIDGenerator(int64(serverconfig.GetCurrentNode().ServerId))
	})
	return defaultIDGen.NextID()
}

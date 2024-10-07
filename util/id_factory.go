package util

import (
	"io/github/gforgame/config"
	"strconv"
	"sync/atomic"
	"time"
)

// IdGenerator 生成全局唯一ID
type IdGenerator struct {
	ServerId int64
}

// generator 用于生成自增长号
var (
	generator atomic.Int64
	idFactory *IdGenerator
)

func init() {
	sid := config.ServerConfig.ServerId
	idFactory = &IdGenerator{ServerId: int64(sid)}
}

// GetNextId 生成全局唯一id
func GetNextId() string {
	// 高16位为serverId
	// 中32位为系统秒数
	// 低16位为自增长号
	// 获取当前时间的秒数
	currentTimeSeconds := time.Now().Unix()
	// 生成ID
	id := (idFactory.ServerId << 48) |
		(currentTimeSeconds << 16) |
		generator.Add(1)&0xFFFF

	return strconv.FormatInt(id, 10)
}

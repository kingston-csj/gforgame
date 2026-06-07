package set

import (
	"sync"
	"time"
)

// ExpireSet 带自动过期的去重容器，用于日志限流
type ExpireSet struct {
	mu       sync.RWMutex
	items    map[string]int64 // key: 唯一标识，value: 过期时间戳
	expire   time.Duration    // 过期时间
	stopChan chan struct{}
}

// NewExpireSet 创建一个自动去重、自动过期的容器
// expire: 同一个key在这段时间内重复都会被拦截
func NewExpireSet(expire time.Duration) *ExpireSet {
	es := &ExpireSet{
		items:    make(map[string]int64),
		expire:   expire,
		stopChan: make(chan struct{}),
	}
	// 启动后台定时清理（每10秒清一次过期key）
	go es.startCleanup(10 * time.Second)
	return es
}

// IsExists 判断key是否存在（存在=true表示需要限流）
// 不存在则自动插入，并返回false
func (es *ExpireSet) IsExists(key string) bool {
	es.mu.Lock()
	defer es.mu.Unlock()

	now := time.Now().Unix()

	// 存在 且 没过期 → 返回true（重复，要限流）
	if expireAt, ok := es.items[key]; ok && expireAt > now {
		return true
	}

	// 不存在/已过期 → 插入，返回false
	es.items[key] = now + int64(es.expire.Seconds())
	return false
}

// 后台清理过期key
func (es *ExpireSet) startCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			es.cleanExpired()
		case <-es.stopChan:
			return
		}
	}
}

// 真正清理逻辑
func (es *ExpireSet) cleanExpired() {
	es.mu.Lock()
	defer es.mu.Unlock()

	now := time.Now().Unix()
	for key, expireAt := range es.items {
		if expireAt <= now {
			delete(es.items, key)
		}
	}
}

// Close 停止清理协程
func (es *ExpireSet) Close() {
	close(es.stopChan)
}
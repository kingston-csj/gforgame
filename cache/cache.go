package cache

import (
	"sync"
	"time"
)

// CacheItem 表示缓存条目
type CacheItem struct {
	Value      interface{}
	LastAccess time.Time
}

// Cache 表示缓存
type Cache struct {
	mu              sync.RWMutex
	items           map[string]*CacheItem
	expiry          time.Duration
	cleanupInterval time.Duration
	loader          func(key string) (interface{}, error)
}

// NewCache 创建一个新的缓存实例
func NewCache(expiry time.Duration, cleanupInterval time.Duration, loader func(key string) (interface{}, error)) *Cache {
	cache := &Cache{
		items:           make(map[string]*CacheItem),
		expiry:          expiry,
		cleanupInterval: cleanupInterval,
		loader:          loader,
	}
	go cache.cleanup() // 启动定期清理线程
	return cache
}

// Get 从缓存中获取数据
func (c *Cache) Get(key string) (interface{}, error) {
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()

	if found {
		// 更新访问时间
		c.mu.Lock()
		item.LastAccess = time.Now()
		c.mu.Unlock()
		return item.Value, nil
	}

	// 如果缓存未命中，从数据库加载数据
	value, err := c.loader(key)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.items[key] = &CacheItem{
		Value:      value,
		LastAccess: time.Now(),
	}
	c.mu.Unlock()

	return value, nil
}

// Set 更新缓存中的数据
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	c.items[key] = &CacheItem{
		Value:      value,
		LastAccess: time.Now(),
	}
	c.mu.Unlock()
}

// cleanup 定期清理沉默缓存
func (c *Cache) cleanup() {
	for {
		time.Sleep(c.cleanupInterval) // 以指定的清理间隔进行清理

		c.mu.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.LastAccess.Add(c.expiry)) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

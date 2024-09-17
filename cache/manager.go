package cache

import (
	"fmt"
	"sync"
	"time"
)

var (
	loaders = map[string]func(key string) (interface{}, error){}
)

func init() {
	loaders = map[string]func(string) (interface{}, error){}
}

type CacheManager struct {
	caches map[string]*Cache
	mu     sync.Mutex
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		caches: make(map[string]*Cache),
	}
}

func (cm *CacheManager) Register(table string, loader func(key string) (interface{}, error)) {
	loaders[table] = loader
}

func (cm *CacheManager) GetCache(table string) (*Cache, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cache, exists := cm.caches[table]; exists {
		return cache, nil
	}

	dbLoader, ok := loaders[table]
	if !ok {
		return nil, fmt.Errorf("cache table %s not found", table)
	}
	cache := NewCache(5*time.Second, 10*time.Second, dbLoader)
	cm.caches[table] = cache
	return cache, nil
}

package cache

import (
	"fmt"
	"sync"
	"time"
)

var (
	loaders map[string]func(key string) (interface{}, error)
)

func init() {
	loaders = make(map[string]func(key string) (interface{}, error))
}

type Manager struct {
	caches map[string]*Cache
	mu     sync.Mutex
}

func NewCacheManager() *Manager {
	return &Manager{
		caches: make(map[string]*Cache),
	}
}

func (cm *Manager) Register(table string, loader func(key string) (interface{}, error)) {
	loaders[table] = loader
}

func (cm *Manager) GetCache(table string) (*Cache, error) {
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

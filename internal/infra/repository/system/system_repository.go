package system

import (
	"github.com/forfun/gforgame/cache"
	"go.uber.org/dig"
)

const systemParameterCacheTable = "systemparameter"

// SystemRepository 系统参数的缓存代理：对外保留原类型名，消费方零改动。
// 内部委托 MySQLSystemRepository 做 db 读写，自身只管 cache 与编排出参。
type SystemRepository struct {
	cacheManager *cache.Manager
	inner        *MySQLSystemRepository
}

type SystemRepositoryParams struct {
	dig.In

	CacheManager *cache.Manager
	Inner        *MySQLSystemRepository
}

func NewSystemRepository(params SystemRepositoryParams) *SystemRepository {
	s := &SystemRepository{
		cacheManager: params.CacheManager,
		inner:        params.Inner,
	}
	s.init()
	return s
}

// init 注册系统参数缓存的回源逻辑（cache miss 时委托 inner 从 db 加载）。
func (r *SystemRepository) init() {
	r.cacheManager.Register(systemParameterCacheTable, func(key string) (any, error) {
		return r.inner.GetRecord(key), nil
	})
}

// LoadValue 读取系统参数原始字符串值。
func (r *SystemRepository) LoadValue(id string) string {
	record := r.getRecord(id)
	if record == nil {
		return ""
	}
	return record.Data
}

// SaveValue 保存系统参数原始字符串值：先更新缓存，再异步落库。
func (r *SystemRepository) SaveValue(id string, value string) {
	record := r.getOrCreateRecord(id)
	record.Data = value

	cache, _ := r.cacheManager.GetCache(systemParameterCacheTable)
	cache.Set(id, record)
	r.inner.SaveRecord(record)
}

// WarmUp 预热指定参数，确保缓存和回源链路已建立。
func (r *SystemRepository) WarmUp(id string) {
	_ = r.LoadValue(id)
}

func (r *SystemRepository) getRecord(id string) *SystemParameterEnt {
	cache, _ := r.cacheManager.GetCache(systemParameterCacheTable)
	cacheEntity, err := cache.Get(id)
	if err != nil || cacheEntity == nil {
		return nil
	}
	record, _ := cacheEntity.(*SystemParameterEnt)
	return record
}

func (r *SystemRepository) getOrCreateRecord(id string) *SystemParameterEnt {
	record := r.getRecord(id)
	if record == nil {
		record = &SystemParameterEnt{}
		record.Id = id
	}
	return record
}

package system

import (
	mysqldb "io/github/gforgame/db"
	"io/github/gforgame/examples/context"
	"sync"

	"gorm.io/gorm"
)

type SystemParameterService struct {
}

func (s *SystemParameterService) init() {
	// 缓存数据读取
	dbLoader := func(key string) (interface{}, error) {
		var p SystemParameterEnt
		result := mysqldb.Db.First(&p, "id=?", key)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				// 未找到记录
				return nil, nil
			}
		}

		return &p, nil
	}
	context.CacheManager.Register("systemparameter", dbLoader)
}

var once sync.Once
var instance *SystemParameterService = &SystemParameterService{}

func GetSystemParameterService() *SystemParameterService {
	return instance
}

func (ps *SystemParameterService) GetSystemParameterRecord(id string) *SystemParameterEnt {
	cache, _ := context.CacheManager.GetCache("systemparameter")
	cacheEntity, err := cache.Get(string(id))
	if err != nil {
		return nil
	}
	if cacheEntity == nil {
		return nil
	}
	record, _ := cacheEntity.(*SystemParameterEnt)
	return record
}

func (ps *SystemParameterService) GetOrCreateSystemParameterRecord(id string) *SystemParameterEnt {
	record := ps.GetSystemParameterRecord(id)
	if record == nil {
		record = &SystemParameterEnt{}
		record.Id = string(id)
	}
	return record
}

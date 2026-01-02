package scene

import (
	mysqldb "io/github/gforgame/db"
	"io/github/gforgame/examples/context"
	"io/github/gforgame/examples/domain/player"
	playerdomain "io/github/gforgame/examples/domain/player"
	"sync"

	"gorm.io/gorm"
)

type SceneService struct {
}

var (
	instance *SceneService
	once     sync.Once
)

func GetSceneService() *SceneService {
	once.Do(func() {
		instance = &SceneService{}

	// 自动建表
	err := mysqldb.Db.AutoMigrate(&player.Scene{})
	if err != nil {
		panic(err)
	}

	// 缓存数据读取
	dbLoader := func(key string) (interface{}, error) {
		var p player.Scene
		result :=mysqldb.Db.First(&p, "id=?", key)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				// 未找到记录
				return nil, nil
			}
		}
		p.AfterFind(nil)
		return &p, nil
	}
	context.CacheManager.Register("scene", dbLoader)
	})
	return instance
}

func (ps *SceneService) GetSceneRecord(playerId string, sceneId string) *playerdomain.Scene {
	cache, _ := context.CacheManager.GetCache("scene")
	cacheEntity, err := cache.Get(playerId + "_" + sceneId)
	if err != nil {
		return nil
	}
	if cacheEntity == nil {
		return nil
	}
	scene, _ := cacheEntity.(*playerdomain.Scene)
	return scene
}

func (s *SceneService) GetOrCreatePlayer(playerId string, sceneId string) *playerdomain.Scene {
	var p playerdomain.Scene
	record := s.GetSceneRecord(playerId, sceneId)
	if record == nil {
		record = &playerdomain.Scene{}
		// 未找到记录
		key := playerId + "_" + sceneId
		record.Id = key
		mysqldb.Db.Create(&p)
	} else {
		p = *record
	}
	return &p
}

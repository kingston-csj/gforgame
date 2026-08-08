package scene

import (
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	mysqldb "github.com/forfun/gforgame/internal/infra/persistence"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

// MySQLSceneRepository 是场景数据的纯持久化仓储，只管 db 读写，不碰缓存。
// 读：同步查询；写：异步落库（经 AsyncDBService）。
type MySQLSceneRepository struct {
	db        *gorm.DB
	dbService *mysqldb.AsyncDBService
}

type MySQLSceneRepositoryParams struct {
	dig.In
	DB        *gorm.DB
	DbService *mysqldb.AsyncDBService
}

func NewMySQLSceneRepository(params MySQLSceneRepositoryParams) *MySQLSceneRepository {
	return &MySQLSceneRepository{
		db:        params.DB,
		dbService: params.DbService,
	}
}

// GetSceneByID 按主键从 db 加载场景数据，未命中返回 nil。
func (r *MySQLSceneRepository) GetSceneByID(key string) *playerdomain.Scene {
	var p playerdomain.Scene
	result := r.db.First(&p, "id=?", key)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		return nil
	}
	return &p
}

// SaveSceneToDb 异步落库场景数据。
func (r *MySQLSceneRepository) SaveSceneToDb(scene *playerdomain.Scene) {
	r.dbService.SaveToDb(scene)
}

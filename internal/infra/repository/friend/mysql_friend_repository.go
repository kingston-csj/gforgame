package friend

import (
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	mysqldb "github.com/forfun/gforgame/internal/infra/persistence"
	"github.com/forfun/gforgame/internal/infra/persistence/po"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

// MySQLFriendRepository 是好友聚合的纯持久化仓储，只管 db 读写，不碰缓存。
// 读：同步查询；写：异步落库（经 AsyncDBService）。
type MySQLFriendRepository struct {
	db        *gorm.DB
	dbService *mysqldb.AsyncDBService
}

type MySQLFriendRepositoryParams struct {
	dig.In
	DB        *gorm.DB
	DbService *mysqldb.AsyncDBService
}

func NewMySQLFriendRepository(params MySQLFriendRepositoryParams) *MySQLFriendRepository {
	return &MySQLFriendRepository{
		db:        params.DB,
		dbService: params.DbService,
	}
}

func (r *MySQLFriendRepository) GetFriendEnt(playerId string) *playerdomain.Friend {
	var p po.FriendPO
	result := r.db.Where("id=?", playerId).Limit(1).Find(&p)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil
	}
	friend, err := p.ToDomain()
	if err != nil {
		return nil
	}
	return friend
}

func (r *MySQLFriendRepository) SaveFriend(friend *playerdomain.Friend) {
	record, err := po.NewFriendPOFromDomain(friend)
	if err != nil {
		panic(err)
	}
	r.dbService.SaveToDb(record)
}

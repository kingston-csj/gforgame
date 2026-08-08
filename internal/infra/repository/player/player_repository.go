package player

import (
	"github.com/forfun/gforgame/cache"
	"github.com/forfun/gforgame/common/eventbus"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/events"
	"go.uber.org/dig"
)

// PlayerRepository 玩家聚合的缓存代理：对外保留原类型名，消费方零改动。
// 内部委托 MySQLPlayerRepository 做 db 读写，自身只管 cache 与事件发布。
// 玩家资料索引职责已拆分到 PlayerProfileService。
type PlayerRepository struct {
	cacheManager *cache.Manager
	inner        *MySQLPlayerRepository
}

type PlayerRepositoryParams struct {
	dig.In

	CacheManager *cache.Manager
	Inner        *MySQLPlayerRepository
}

func NewPlayerRepository(params PlayerRepositoryParams) *PlayerRepository {
	return &PlayerRepository{
		cacheManager: params.CacheManager,
		inner:        params.Inner,
	}
}

// Init 注册玩家缓存的回源逻辑（cache miss 时委托 inner 从 db 加载并发事件）。
func (r *PlayerRepository) Init() {
	r.cacheManager.Register("player", func(key string) (any, error) {
		player := r.inner.GetPlayerByID(key)
		if player == nil {
			return nil, nil
		}
		eventbus.Default().Publish(events.PlayerAfterLoad, player)
		return player, nil
	})
}

func (r *PlayerRepository) GetPlayer(playerId string) *playerdomain.Player {
	cache, _ := r.cacheManager.GetCache("player")
	cacheEntity, err := cache.Get(playerId)
	if err != nil {
		return nil
	}
	if cacheEntity == nil {
		return nil
	}
	player, _ := cacheEntity.(*playerdomain.Player)
	return player
}

func (r *PlayerRepository) SavePlayer(player *playerdomain.Player) {
	cache, _ := r.cacheManager.GetCache("player")
	cache.Set(player.GetId(), player)
	r.inner.SavePlayerToDb(player)
}

func (r *PlayerRepository) FindTopPlayersByFightPower() ([]*playerdomain.Player, error) {
	return r.inner.FindTopPlayersByFightPower()
}

func (r *PlayerRepository) FindTopPlayersByMaxGuanka() ([]*playerdomain.Player, error) {
	return r.inner.FindTopPlayersByMaxGuanka()
}

func (r *PlayerRepository) FindTopPlayersByArenaScore() ([]*playerdomain.Player, error) {
	return r.inner.FindTopPlayersByArenaScore()
}

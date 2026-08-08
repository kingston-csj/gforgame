package player

import (
	"go.uber.org/dig"
	"gorm.io/gorm"

	"github.com/forfun/gforgame/common/logger"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	mysqldb "github.com/forfun/gforgame/internal/infra/persistence"
	playerpo "github.com/forfun/gforgame/internal/infra/persistence/po"
)

// MySQLPlayerRepository 纯 db 仓储：只管玩家聚合的持久化读写，不碰 cache。
// toDomain 内含 AfterLoad，所需配置提供器以 ItemConfigProviders 形式注入。
type MySQLPlayerRepository struct {
	db        *gorm.DB
	dbService *mysqldb.AsyncDBService
	providers playerdomain.ItemConfigProviders
}

type MySQLPlayerRepositoryParams struct {
	dig.In

	DB        *gorm.DB
	DbService *mysqldb.AsyncDBService
	Providers playerdomain.ItemConfigProviders
}

func NewMySQLPlayerRepository(params MySQLPlayerRepositoryParams) *MySQLPlayerRepository {
	return &MySQLPlayerRepository{
		db:        params.DB,
		dbService: params.DbService,
		providers: params.Providers,
	}
}

// GetPlayerByID 按主键从 db 加载玩家聚合（含配置绑定与 AfterLoad），未命中返回 nil。
func (r *MySQLPlayerRepository) GetPlayerByID(playerId string) *playerdomain.Player {
	var record playerpo.PlayerPO
	result := r.db.First(&record, "id=?", playerId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		return nil
	}
	player, err := r.toDomain(&record)
	if err != nil {
		return nil
	}
	return player
}

// SavePlayerToDb 异步落库玩家聚合。
func (r *MySQLPlayerRepository) SavePlayerToDb(player *playerdomain.Player) {
	record, err := playerpo.NewPlayerPOFromDomain(player)
	if err != nil || record == nil {
		logger.Error("SavePlayer error: %v", err)
		return
	}
	r.dbService.SaveToDb(record)
}

// LoadPlayerProfilesFromDb 从 db 预加载玩家资料摘要列表。
func (r *MySQLPlayerRepository) LoadPlayerProfilesFromDb() []*playerdomain.PlayerProfile {
	var profiles []*playerdomain.PlayerProfile
	err := r.db.Model(&playerpo.PlayerPO{}).
		Select("id, name, level, head, camp, fight, guanka").
		Scan(&profiles).Error
	if err != nil {
		panic(err)
	}
	return profiles
}

// FindTopPlayersByFightPower 按战力排行查 db（返回部分字段的摘要）。
func (r *MySQLPlayerRepository) FindTopPlayersByFightPower() ([]*playerdomain.Player, error) {
	var records []*playerpo.PlayerPO
	err := r.db.Where("fight >?", 0).Order("level desc").Limit(50).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return r.toSummaries(records), nil
}

func (r *MySQLPlayerRepository) FindTopPlayersByMaxGuanka() ([]*playerdomain.Player, error) {
	var records []*playerpo.PlayerPO
	err := r.db.Where("guanka >?", 0).Order("guanka desc").Limit(50).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return r.toSummaries(records), nil
}

func (r *MySQLPlayerRepository) FindTopPlayersByArenaScore() ([]*playerdomain.Player, error) {
	var records []*playerpo.PlayerPO
	err := r.db.Where("arena_score >?", 0).Order("arena_score desc").Limit(50).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return r.toSummaries(records), nil
}

func (r *MySQLPlayerRepository) toDomain(record *playerpo.PlayerPO) (*playerdomain.Player, error) {
	player, err := record.ToDomain()
	if err != nil {
		return nil, err
	}
	if err = player.AfterLoad(r.providers); err != nil {
		return nil, err
	}
	return player, nil
}

func (r *MySQLPlayerRepository) toSummaries(records []*playerpo.PlayerPO) []*playerdomain.Player {
	result := make([]*playerdomain.Player, 0, len(records))
	for _, record := range records {
		if record == nil {
			continue
		}
		result = append(result, &playerdomain.Player{
			BaseEntity:  record.BaseEntity,
			Name:        record.Name,
			Head:        record.Head,
			Level:       record.Level,
			Guanka:      record.Guanka,
			Fight:       record.Fight,
			Camp:        record.Camp,
			ArenaScore:  record.ArenaScore,
			CreateTime:  record.CreateTime,
			VipLevel:    record.VipLevel,
			RechargeRmb: record.RechargeRmb,
		})
	}
	return result
}

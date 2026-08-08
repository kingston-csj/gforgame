package po

import (
	"encoding/json"

	util "github.com/forfun/gforgame/common/util/conv"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/persist"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// PlayerPO 表示玩家聚合的数据库持久化对象。
type PlayerPO struct {
	persist.BaseEntity
	Name              string `gorm:"column:name"`
	Head              int32  `gorm:"column:head;default:0"`
	ClientScore       int32  `gorm:"column:client_score"`
	ClientData        string `gorm:"column:client_data"`
	RechargeRmb       int32  `gorm:"column:recharge_rmb"`
	VipLevel          int32  `gorm:"column:vip_level"`
	CreateTime        int64  `gorm:"column:create_time"`
	Level             int32  `gorm:"column:level"`
	Guanka            int32  `gorm:"column:guanka"`
	Stage             int32  `gorm:"column:stage"`
	ArenaScore        int32  `gorm:"column:arena_score"`
	BackpackJson      string `gorm:"column:backpack"`
	RuneBackpackJson  string `gorm:"column:rune_backpack"`
	SceneBackpackJson string `gorm:"column:scene_backpack"`
	EquipBackpackJson string `gorm:"column:equip_backpack"`
	CardBackpackJson  string `gorm:"column:card_backpack"`
	HeroBoxJson       string `gorm:"column:hero_box"`
	PetBoxJson        string `gorm:"column:pet_box"`
	DailyResetJson    string `gorm:"column:daily_reset"`
	WeeklyResetJson   string `gorm:"column:weeklyreset"`
	MonthlyResetJson  string `gorm:"column:monthly_reset"`
	Fight             int32  `gorm:"column:fight"`
	Camp              int32  `gorm:"column:camp"`
	MailboxJson       string `gorm:"column:mailbox"`
	ExtendBoxJson     string `gorm:"column:extend_box"`
	QuestBoxJson      string `gorm:"column:quest_box"`
	RechargeBoxJson   string `gorm:"column:recharge_box"`
	ArenaBoxJson      string `gorm:"column:arena_box"`
	ActivityBoxJson   string `gorm:"column:activity_box"`
	FightBoxJson      string `gorm:"column:fight_box"`
	EquipBoxJson      string `gorm:"column:equip_box"`
}

func (p *PlayerPO) GetKey() string {
	return "Player_" + p.Id
}

func (p *PlayerPO) TableName() string {
	return schema.NamingStrategy{}.TableName("Player")
}

func (p *PlayerPO) BeforePersist() error {
	return nil
}

func (p *PlayerPO) AfterLoad() error {
	return nil
}

func (p *PlayerPO) BeforeSave(tx *gorm.DB) error {
	return p.BeforePersist()
}

func (p *PlayerPO) AfterFind(tx *gorm.DB) error {
	return p.AfterLoad()
}

// NewPlayerPOFromDomain 将运行时领域对象转换为持久化对象。
func NewPlayerPOFromDomain(player *playerdomain.Player) (*PlayerPO, error) {
	if player == nil {
		return nil, nil
	}
	result := &PlayerPO{
		BaseEntity:  player.BaseEntity,
		Name:        player.Name,
		Head:        player.Head,
		ClientScore: player.ClientScore,
		ClientData:  player.ClientData,
		RechargeRmb: player.RechargeRmb,
		VipLevel:    player.VipLevel,
		CreateTime:  player.CreateTime,
		Level:       player.Level,
		Guanka:      player.Guanka,
		Stage:       player.Stage,
		ArenaScore:  player.ArenaScore,
		Fight:       player.Fight,
		Camp:        player.Camp,
	}
	var err error
	if result.BackpackJson, err = marshalJSON(player.Backpack); err != nil {
		return nil, err
	}
	if result.RuneBackpackJson, err = marshalJSON(player.RuneBackpack); err != nil {
		return nil, err
	}
	if result.HeroBoxJson, err = marshalJSON(player.HeroBox); err != nil {
		return nil, err
	}
	if result.DailyResetJson, err = marshalJSON(player.DailyReset); err != nil {
		return nil, err
	}
	if result.WeeklyResetJson, err = marshalJSON(player.WeeklyReset); err != nil {
		return nil, err
	}
	if result.MonthlyResetJson, err = marshalJSON(player.MonthlyReset); err != nil {
		return nil, err
	}
	if result.MailboxJson, err = marshalJSON(player.Mailbox); err != nil {
		return nil, err
	}
	if result.ExtendBoxJson, err = marshalJSON(player.ExtendBox); err != nil {
		return nil, err
	}
	if result.QuestBoxJson, err = marshalJSON(player.QuestBox); err != nil {
		return nil, err
	}
	if result.RechargeBoxJson, err = marshalJSON(player.RechargeBox); err != nil {
		return nil, err
	}
	if result.ArenaBoxJson, err = marshalJSON(player.ArenaBox); err != nil {
		return nil, err
	}
	if result.ActivityBoxJson, err = marshalJSON(player.ActivityBox); err != nil {
		return nil, err
	}
	return result, nil
}

// ToDomain 将持久化对象恢复为领域对象。
func (p *PlayerPO) ToDomain() (*playerdomain.Player, error) {
	if p == nil {
		return nil, nil
	}
	result := &playerdomain.Player{
		BaseEntity:  p.BaseEntity,
		Name:        p.Name,
		Head:        p.Head,
		ClientScore: p.ClientScore,
		ClientData:  p.ClientData,
		RechargeRmb: p.RechargeRmb,
		VipLevel:    p.VipLevel,
		CreateTime:  p.CreateTime,
		Level:       p.Level,
		Guanka:      p.Guanka,
		Stage:       p.Stage,
		ArenaScore:  p.ArenaScore,
		Fight:       p.Fight,
		Camp:        p.Camp,
	}
	loadJSON(p.BackpackJson, &result.Backpack)
	loadJSON(p.RuneBackpackJson, &result.RuneBackpack)
	loadJSON(p.HeroBoxJson, &result.HeroBox)
	loadJSON(p.DailyResetJson, &result.DailyReset)
	loadJSON(p.WeeklyResetJson, &result.WeeklyReset)
	loadJSON(p.MonthlyResetJson, &result.MonthlyReset)
	loadJSON(p.MailboxJson, &result.Mailbox)
	loadJSON(p.ExtendBoxJson, &result.ExtendBox)
	loadJSON(p.QuestBoxJson, &result.QuestBox)
	loadJSON(p.RechargeBoxJson, &result.RechargeBox)
	loadJSON(p.ArenaBoxJson, &result.ArenaBox)
	loadJSON(p.ActivityBoxJson, &result.ActivityBox)
	return result, nil
}

func marshalJSON[T any](component *T) (string, error) {
	if component == nil {
		return "", nil
	}
	data, err := json.Marshal(component)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func loadJSON[T any](jsonStr string, target **T) {
	if util.IsEmptyString(jsonStr) {
		return
	}
	value := new(T)
	if err := json.Unmarshal([]byte(jsonStr), value); err != nil {
		return
	}
	*target = value
}

func (p *PlayerPO) SnapshotEntity() (persist.Entity, error) {
	cp := *p
	return &cp, nil
}

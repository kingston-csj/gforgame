package player

import (
	configcontract "github.com/forfun/gforgame/internal/config/contracts"
	"github.com/forfun/gforgame/internal/fight/attribute"
	"github.com/forfun/gforgame/internal/io"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/persist"
)

type Player struct {
	persist.BaseEntity
	Name         string             `gorm:"player's name"`
	Head         int32              `gorm:"player's head default:0"`
	ClientScore  int32              `gorm:"player's client score"`
	ClientData   string             `gorm:"player's client data"`
	RechargeRmb  int32              `gorm:"player's recharge rmb"`
	VipLevel     int32              `gorm:"player's vip level"`
	CreateTime   int64              `gorm:"player's create time"`
	Level        int32              `gorm:"player's' level"`
	Guanka       int32              `gorm:"player's' guanka"`
	Stage        int32              `gorm:"player's stage"`
	ArenaScore   int32              `gorm:"player's arena score"`
	Backpack     *Backpack          `gorm:"-"`
	RuneBackpack *Backpack          `gorm:"-"`
	HeroBox      *HeroBox           `gorm:"-"`
	Purse        *Purse             `gorm:"-"`
	AttrBox      *attribute.AttrBox `gorm:"-"`
	DailyReset   *DailyReset        `gorm:"-"`
	WeeklyReset  *WeeklyReset       `gorm:"-"`
	// 月度重置
	MonthlyReset *MonthlyResetBox `gorm:"-"`
	Fight        int32            `gorm:"player's fight"`
	Camp         int32            `gorm:"player's camp"`
	Mailbox      *Mailbox         `gorm:"-"`
	ExtendBox    *ExtendBox       `gorm:"-"`
	QuestBox     *QuestBox        `gorm:"-"`
	RechargeBox  *RechargeBox     `gorm:"-"`
	// 竞技场数据
	ArenaBox    *ArenaBox    `gorm:"-"`
	ActivityBox *ActivityBox `gorm:"-"`
}

// 数据重置，仅用于gm
func (p *Player) Reset() {
	p.Camp = 0
	p.Level = 0
	p.Stage = 0
	p.Fight = 0
	p.VipLevel = 0
	p.CreateTime = 0
	p.Name = ""
}

func (p *Player) AfterLoad(providers ItemConfigProviders) error {
	p.ensureBackpack(&p.Backpack, providers.Base, 9999)
	p.ensureBackpack(&p.RuneBackpack, providers.Rune, 9999)
	p.ensureComponent(&p.HeroBox, func() *HeroBox {
		return &HeroBox{
			Heros: make(map[int32]*Hero),
		}
	})

	p.ensureComponent(&p.DailyReset, func() *DailyReset {
		return &DailyReset{}
	})
	p.ensureComponent(&p.WeeklyReset, func() *WeeklyReset {
		return &WeeklyReset{}
	})
	p.ensureComponent(&p.MonthlyReset, func() *MonthlyResetBox {
		return &MonthlyResetBox{}
	})
	p.ensureComponent(&p.Mailbox, func() *Mailbox {
		return &Mailbox{
			Mails: make(map[string]*Mail),
		}
	})
	p.ensureComponent(&p.ExtendBox, func() *ExtendBox {
		return &ExtendBox{}
	})
	p.AttrBox = attribute.NewAttrBox()
	for _, hero := range p.HeroBox.Heros {
		hero.AttrBox = attribute.NewAttrBox()
	}
	p.ensureComponent(&p.QuestBox, func() *QuestBox {
		return &QuestBox{
			Doing:    make(map[int32]*Quest),
			Finished: make(map[int32]bool),
		}
	})
	p.ensureComponent(&p.RechargeBox, func() *RechargeBox {
		return &RechargeBox{}
	})
	p.ensureComponent(&p.ArenaBox, func() *ArenaBox {
		return &ArenaBox{}
	})
	p.ensureComponent(&p.ActivityBox, func() *ActivityBox {
		return &ActivityBox{}
	})

	return nil
}

func (p *Player) ensureBackpack(target **Backpack, provider configcontract.ItemConfigProvider, defaultCapacity int32) {
	if *target == nil {
		*target = &Backpack{
			Items:    make(map[string]*Item),
			Capacity: defaultCapacity,
		}
	}
	(*target).configProvider = provider
	if defaultCapacity > 0 && (*target).Capacity == 0 {
		(*target).Capacity = defaultCapacity
	}
	(*target).AfterLoad()
}

func (p *Player) ensureComponent(target any, factory any) {
	switch t := target.(type) {
	case **HeroBox:
		ensureValue(t, factory.(func() *HeroBox))
	case **DailyReset:
		ensureValue(t, factory.(func() *DailyReset))
	case **WeeklyReset:
		ensureValue(t, factory.(func() *WeeklyReset))
	case **MonthlyResetBox:
		ensureValue(t, factory.(func() *MonthlyResetBox))
	case **Mailbox:
		ensureValue(t, factory.(func() *Mailbox))
	case **ExtendBox:
		ensureValue(t, factory.(func() *ExtendBox))
	case **QuestBox:
		ensureValue(t, factory.(func() *QuestBox))
	case **RechargeBox:
		ensureValue(t, factory.(func() *RechargeBox))
	case **ArenaBox:
		ensureValue(t, factory.(func() *ArenaBox))
	case **ActivityBox:
		ensureValue(t, factory.(func() *ActivityBox))
	}
}

func ensureValue[T any](target **T, factory func() *T) {
	if *target == nil {
		*target = factory()
	}
	if loader, ok := any(*target).(PostLoader); ok {
		loader.AfterLoad()
	}
}

// PostLoader 接口，用于在加载 JSON 后进行额外的初始化操作
type PostLoader interface {
	AfterLoad()
}

func (p *Player) GetId() string {
	return p.Id
}

func (p *Player) GetName() string {
	return p.Name
}

func (p *Player) NotifyPurseChange() {
	resPurse := &protos.PushPurseInfo{}
	resPurse.Diamond = p.Purse.Diamond
	resPurse.Gold = p.Purse.Gold
	io.NotifyPlayer(p, resPurse)
}

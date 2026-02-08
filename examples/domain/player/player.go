package player

import (
	"encoding/json"

	"io/github/gforgame/db"
	"io/github/gforgame/examples/fight/attribute"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/protos"
	"io/github/gforgame/util"

	"gorm.io/gorm"
)

type Player struct {
	db.BaseEntity
	Name           string             `gorm:"player's name"`
	Head           int32              `gorm:"player's head default:0"`
	ClientScore    int32              `gorm:"player's client score"`
	ClientData     string             `gorm:"player's client data"`
	RechargeRmb    int32              `gorm:"player's recharge rmb"`
	VipLevel       int32              `gorm:"player's vip level"`
	CreateTime     int64              `gorm:"player's create time"`
	Level          int32              `gorm:"player's' level"`
	Stage          int32              `gorm:"player's stage"`
	ArenaScore int32              `gorm:"player's arena score"`
	Backpack       *Backpack          `gorm:"-"`
	BackpackJson   string             `gorm:"backpack"`
	RuneBackpack       *Backpack          `gorm:"-"`
	RuneBackpackJson   string             `gorm:"runeBackpack"`
	SceneBackpack       *Backpack          `gorm:"-"`
	SceneBackpackJson   string             `gorm:"sceneBackpack"`
	HeroBox        *HeroBox           `gorm:"-"`
	HeroBoxJson    string             `gorm:"herobox"`
	Purse          *Purse             `gorm:"-"`
	AttrBox        *attribute.AttrBox `gorm:"-"`
	PurseJson      string             `gorm:"purse"`
	DailyReset     *DailyReset        `gorm:"-"`
	DailyResetJson string             `gorm:"dailyreset"`
	WeeklyReset     *WeeklyReset        `gorm:"-"`
	WeeklyResetJson string             `gorm:"weeklyreset"`
	// 月度重置
	MonthlyReset     *MonthlyResetBox        `gorm:"-"`
	MonthlyResetJson string             `gorm:"monthlyreset"`
	Fight          int32              `gorm:"player's fight"`
	Camp           int32              `gorm:"player's camp"`
	Mailbox        *Mailbox           `gorm:"-"`
	MailboxJson    string             `gorm:"mailbox"`
	ExtendBox      *ExtendBox         `gorm:"-"`
	ExtendBoxJson  string             `gorm:"extendbox"`
	QuestBox       *QuestBox          `gorm:"-"`
	QuestBoxJson   string             `gorm:"questbox"`
	RechargeBox    *RechargeBox       `gorm:"-"`
	RechargeBoxJson string             `gorm:"rechargebox"`
	// 竞技场数据
	ArenaBox       *ArenaBox          `gorm:"-"`
	ArenaBoxJson   string             `gorm:"arenabox"`
}

func (p *Player) BeforeSave(tx *gorm.DB) error {
	if err := saveJSON(p.Backpack, &p.BackpackJson); err != nil {
		return err
	}
	if err := saveJSON(p.SceneBackpack, &p.SceneBackpackJson); err != nil {
		return err
	}
	if err := saveJSON(p.RuneBackpack, &p.RuneBackpackJson); err != nil {
		return err
	}
	if err := saveJSON(p.HeroBox, &p.HeroBoxJson); err != nil {
		return err
	}
	if err := saveJSON(p.Purse, &p.PurseJson); err != nil {
		return err
	}
	if err := saveJSON(p.DailyReset, &p.DailyResetJson); err != nil {
		return err
	}
	if err := saveJSON(p.WeeklyReset, &p.WeeklyResetJson); err != nil {
		return err
	}
	if err := saveJSON(p.MonthlyReset, &p.MonthlyResetJson); err != nil {
		return err
	}
	if err := saveJSON(p.Mailbox, &p.MailboxJson); err != nil {
		return err
	}
	if err := saveJSON(p.ExtendBox, &p.ExtendBoxJson); err != nil {
		return err
	}
	if err := saveJSON(p.QuestBox, &p.QuestBoxJson); err != nil {
		return err
	}
	if err := saveJSON(p.RechargeBox, &p.RechargeBoxJson); err != nil {
		return err
	}
	if err := saveJSON(p.ArenaBox, &p.ArenaBoxJson); err != nil {
		return err
	}

	return nil
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
	p.BackpackJson = ""
	p.SceneBackpackJson = ""
	p.RuneBackpackJson = ""
	p.HeroBoxJson = ""
	p.PurseJson = ""
	p.DailyResetJson = ""
	p.WeeklyResetJson = ""
	p.MonthlyResetJson = ""
	p.MailboxJson = ""
	p.ExtendBoxJson = ""
	p.QuestBoxJson = ""
	p.RechargeBoxJson = ""
	p.AfterFind(nil)
}


func (p *Player) AfterFind(tx *gorm.DB) error {
	loadJSON(p.BackpackJson, &p.Backpack, func() *Backpack {
		return &Backpack{
			Items:          make(map[string]*Item),
			configProvider: BaseItemConfigProviderInstance,
			Capacity:       9999,
		}
	})

	loadJSON(p.SceneBackpackJson, &p.SceneBackpack, func() *Backpack {
		return &Backpack{
			Items:          make(map[string]*Item),
			configProvider: SceneItemConfigProviderInstance,
		}
	})

	loadJSON(p.RuneBackpackJson, &p.RuneBackpack, func() *Backpack {
		return &Backpack{
			Items:          make(map[string]*Item),
			configProvider: RuneConfigProviderInstance,
		}
	})

	loadJSON(p.HeroBoxJson, &p.HeroBox, func() *HeroBox {
		return &HeroBox{
			Heros: make(map[int32]*Hero),
		}
	})

	loadJSON(p.PurseJson, &p.Purse, func() *Purse {
		return &Purse{
			Diamond: 0,
			Gold:    0,
		}
	})

	p.AttrBox = attribute.NewAttrBox()

	loadJSON(p.DailyResetJson, &p.DailyReset, func() *DailyReset {
		return &DailyReset{
		}
	})

	loadJSON(p.WeeklyResetJson, &p.WeeklyReset, func() *WeeklyReset {
		return &WeeklyReset{
		}
	})

	loadJSON(p.MonthlyResetJson, &p.MonthlyReset, func() *MonthlyResetBox {
		return &MonthlyResetBox{
		}
	})

	loadJSON(p.MailboxJson, &p.Mailbox, func() *Mailbox {
		return &Mailbox{
			Mails: make(map[string]*Mail),
		}
	})

	loadJSON(p.ExtendBoxJson, &p.ExtendBox, func() *ExtendBox {
		return &ExtendBox{}
	})

	for _, hero := range p.HeroBox.Heros {
		hero.AttrBox = attribute.NewAttrBox()
	}

	loadJSON(p.QuestBoxJson, &p.QuestBox, func() *QuestBox {
		return &QuestBox{
			Doing:    make(map[int32]*Quest),
			Finished: make(map[int32]bool),
		}
	})

	loadJSON(p.RechargeBoxJson, &p.RechargeBox, func() *RechargeBox {
		return &RechargeBox{}
	})

	loadJSON(p.ArenaBoxJson, &p.ArenaBox, func() *ArenaBox {
		return &ArenaBox{}
	})

	return nil
}

// PostLoader 接口，用于在加载 JSON 后进行额外的初始化操作
type PostLoader interface {
	AfterLoad()
}

func saveJSON[T any](component *T, jsonTarget *string) error {
	if component == nil {
		*jsonTarget = ""
		return nil
	}
	jsonData, err := json.Marshal(component)
	if err != nil {
		return err
	}
	*jsonTarget = string(jsonData)
	return nil
}

func loadJSON[T any](jsonStr string, target **T, factory func() *T) {
	val := factory()
	if !util.IsEmptyString(jsonStr) {
		_ = json.Unmarshal([]byte(jsonStr), val)
	}
	// 如果实现了 PostLoader 接口，则调用 AfterLoad 进行后处理
	if loader, ok := any(val).(PostLoader); ok {
		loader.AfterLoad()
	}
	*target = val
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

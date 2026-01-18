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
	RechargeRmb    int32              `gorm:"player's recharge rmb"`
	VipLevel       int32              `gorm:"player's vip level"`
	CreateTime     int64              `gorm:"player's create time"`
	Level          int32              `gorm:"player's' level"`
	Stage          int32              `gorm:"player's stage"`
	Backpack       *Backpack          `gorm:"-"`
	BackpackJson   string             `gorm:"backpack"`
	RuneBackpack       *Backpack          `gorm:"-"`
	RuneBackpackJson   string             `gorm:"runeBackpack"`
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
}

func (p *Player) BeforeSave(tx *gorm.DB) error {
	if p.Backpack == nil {
		p.BackpackJson = ""
	} else {
		jsonData, err := json.Marshal(p.Backpack)
		if err != nil {
			return err
		}
		p.BackpackJson = string(jsonData)
	}

	if p.RuneBackpack == nil {
		p.RuneBackpackJson = ""
	} else {
		jsonData, err := json.Marshal(p.RuneBackpack)
		if err != nil {
			return err
		}
		p.RuneBackpackJson = string(jsonData)
	}

	if p.HeroBox == nil {
		p.HeroBoxJson = ""
	} else {
		jsonData, err := json.Marshal(p.HeroBox)
		if err != nil {
			return err
		}
		p.HeroBoxJson = string(jsonData)
	}
	if p.Purse == nil {
		p.PurseJson = ""
	} else {
		jsonData, err := json.Marshal(p.Purse)
		if err != nil {
			return err
		}
		p.PurseJson = string(jsonData)
	}
	if p.DailyReset == nil {
		p.DailyResetJson = ""
	} else {
		jsonData, err := json.Marshal(p.DailyReset)
		if err != nil {
			return err
		}
		p.DailyResetJson = string(jsonData)
	}
	if p.WeeklyReset == nil {
		p.WeeklyResetJson = ""
	} else {
		jsonData, err := json.Marshal(p.WeeklyReset)
		if err != nil {
			return err
		}
		p.WeeklyResetJson = string(jsonData)
	}
	if p.MonthlyReset == nil {
		p.MonthlyResetJson = ""
	} else {
		jsonData, err := json.Marshal(p.MonthlyReset)
		if err != nil {
			return err
		}
		p.MonthlyResetJson = string(jsonData)
	}
	if p.Mailbox == nil {
		p.MailboxJson = ""
	} else {
		jsonData, err := json.Marshal(p.Mailbox)
		if err != nil {
			return err
		}
		p.MailboxJson = string(jsonData)
	}
	if p.ExtendBox == nil {
		p.ExtendBoxJson = ""
	} else {
		jsonData, err := json.Marshal(p.ExtendBox)
		if err != nil {
			return err
		}
		p.ExtendBoxJson = string(jsonData)
	}
	if p.QuestBox == nil {
		p.QuestBoxJson = ""
	} else {
		jsonData, err := json.Marshal(p.QuestBox)
		if err != nil {
			return err
		}
		p.QuestBoxJson = string(jsonData)
	}

	if p.RechargeBox == nil {
		p.RechargeBoxJson = ""
	} else {
		jsonData, err := json.Marshal(p.RechargeBox)
		if err != nil {
			return err
		}
		p.RechargeBoxJson = string(jsonData)
	}


	return nil
}
func (p *Player) AfterFind(tx *gorm.DB) error {
	if util.IsEmptyString(p.BackpackJson) {
		p.Backpack = &Backpack{
			Items: make(map[string]*Item),
			configProvider: BaseItemConfigProviderInstance,
			Capacity: 9999,
		}
	} else {
		backpack := &Backpack{
			Items: make(map[string]*Item),
			configProvider: BaseItemConfigProviderInstance,
			Capacity: 9999,
		}
		json.Unmarshal([]byte(p.BackpackJson), &backpack)
		p.Backpack = backpack
	}

	if util.IsEmptyString(p.RuneBackpackJson) {
		p.RuneBackpack = &Backpack{
			Items: make(map[string]*Item),
			configProvider: RuneConfigProviderInstance,
		}
	} else {
		backpack := &Backpack{
			Items: make(map[string]*Item),
			configProvider: RuneConfigProviderInstance,
		}
		json.Unmarshal([]byte(p.RuneBackpackJson), &backpack)
		p.RuneBackpack = backpack
	}

	if util.IsEmptyString(p.HeroBoxJson) {
		p.HeroBox = &HeroBox{
			Heros: make(map[int32]*Hero),
		}
	} else {
		json.Unmarshal([]byte(p.HeroBoxJson), &p.HeroBox)
	}
	if util.IsEmptyString(p.PurseJson) {
		p.Purse = &Purse{
			Diamond: 0,
			Gold:    0,
		}
	} else {
		json.Unmarshal([]byte(p.PurseJson), &p.Purse)
	}
	p.AttrBox = attribute.NewAttrBox()
	if util.IsEmptyString(p.DailyResetJson) {
		p.DailyReset = &DailyReset{
			LastDailyReset:  0,
			DailyQuestScore: 0,
		}
	} else {
		json.Unmarshal([]byte(p.DailyResetJson), &p.DailyReset)
	}
	if util.IsEmptyString(p.WeeklyResetJson) {
		p.WeeklyReset = &WeeklyReset{
			LastWeeklyReset:  0,
			WeeklyQuestScore: 0,
		}
	} else {
		json.Unmarshal([]byte(p.WeeklyResetJson), &p.WeeklyReset)
	}
	if util.IsEmptyString(p.MonthlyResetJson) {
		p.MonthlyReset = &MonthlyResetBox{
			ResetTime:  0,
			SignInDays: make([]int32, 0),
		}
	} else {
		json.Unmarshal([]byte(p.MonthlyResetJson), &p.MonthlyReset)
	}
	if util.IsEmptyString(p.MailboxJson) {
		p.Mailbox = &Mailbox{
			Mails: make(map[int64]*Mail),
		}
	} else {
		json.Unmarshal([]byte(p.MailboxJson), &p.Mailbox)
	}
	if util.IsEmptyString(p.ExtendBoxJson) {
		p.ExtendBox = &ExtendBox{
			PrivateChats: make(map[string][]ChatMessage),
		}
	} else {
		json.Unmarshal([]byte(p.ExtendBoxJson), &p.ExtendBox)
	}
	for _, hero := range p.HeroBox.Heros {
		hero.AttrBox = attribute.NewAttrBox()
	}
	if util.IsEmptyString(p.QuestBoxJson) {
		p.QuestBox = &QuestBox{
			Doing:    make(map[int32]*Quest),
			Finished: make(map[int32]bool),
		}
	} else {
		json.Unmarshal([]byte(p.QuestBoxJson), &p.QuestBox)
	}

	if util.IsEmptyString(p.RechargeBoxJson) {
		p.RechargeBox = &RechargeBox{
			ActivatedQiRiPay: 0,
			RechargeTimes:    make(map[int32]int32),
			ActivatedPassPay: 0,
		}
	} else {
		json.Unmarshal([]byte(p.RechargeBoxJson), &p.RechargeBox)
	}

	return nil
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

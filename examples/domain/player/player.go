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
	Level          int32              `gorm:"player's' level"`
	Stage          int32              `gorm:"player's stage"`
	Backpack       *Backpack          `gorm:"-"`
	BackpackJson   string             `gorm:"backpack"`
	HeroBox        *HeroBox           `gorm:"-"`
	HeroBoxJson    string             `gorm:"herobox"`
	Purse          *Purse             `gorm:"-"`
	AttrBox        *attribute.AttrBox `gorm:"-"`
	PurseJson      string             `gorm:"purse"`
	DailyReset     *DailyReset        `gorm:"-"`
	DailyResetJson string             `gorm:"dailyreset"`
	Fight          int32              `gorm:"player's fight"`
	Camp           int32              `gorm:"player's camp"`
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
	return nil
}

func (p *Player) AfterFind(tx *gorm.DB) error {
	if util.IsEmptyString(p.BackpackJson) {
		p.Backpack = &Backpack{
			Items: make(map[int32]int32),
		}
	} else {
		json.Unmarshal([]byte(p.BackpackJson), &p.Backpack)
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
	}

	for _, hero := range p.HeroBox.Heros {
		hero.AttrBox = attribute.NewAttrBox()
	}

	return nil
}

func (p *Player) GetId() string {
	return p.Id
}

func (p *Player) NotifyPurseChange() {
	resPurse := &protos.PushPurseInfo{}
	resPurse.Diamond = p.Purse.Diamond
	resPurse.Gold = p.Purse.Gold
	io.NotifyPlayer(p, resPurse)
}

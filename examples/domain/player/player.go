package player

import (
	"encoding/json"
	"io/github/gforgame/db"
	"io/github/gforgame/examples/utils"

	"gorm.io/gorm"
)

type Player struct {
	db.BaseEntity
	ID           string    `gorm:"player's ID"`
	Name         string    `gorm:"player's name"`
	Level        uint      `gorm:"player's' level"`
	Backpack     *Backpack `gorm:"-"`
	BackpackJson string    `gorm:"backpack"`
	HeroBox      *HeroBox  `gorm:"-"`
	HeroBoxJson  string    `gorm:"herobox"`
	Purse        *Purse    `gorm:"-"`
	PurseJson    string    `gorm:"purse"`
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
	return nil
}

func (p *Player) AfterFind(tx *gorm.DB) error {
	if utils.IsEmpty(p.BackpackJson) {
		p.Backpack = &Backpack{
			Items: make(map[int32]int32),
		}
	} else {
		json.Unmarshal([]byte(p.BackpackJson), &p.Backpack)
	}
	if utils.IsEmpty(p.HeroBoxJson) {
		p.HeroBox = &HeroBox{
			Heros: make(map[int32]*Hero),
		}
	} else {
		json.Unmarshal([]byte(p.HeroBoxJson), &p.HeroBox)
	}
	if utils.IsEmpty(p.PurseJson) {
		p.Purse = &Purse{
			Diamond: 0,
			Gold:    0,
		}
	} else {
		json.Unmarshal([]byte(p.PurseJson), &p.Purse)
	}
	return nil
}

func (p *Player) GetID() string {
	return p.ID
}

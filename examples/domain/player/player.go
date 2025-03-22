package player

import (
	"encoding/json"
	"io/github/gforgame/db"
	"io/github/gforgame/examples/utils"

	"gorm.io/gorm"
)

type Player struct {
	db.BaseEntity
	Name         string    `gorm:"player's name"`
	Level        uint      `gorm:"player's' level"`
	Backpack     *Backpack `gorm:"-"`
	BackpackJson string    `gorm:"backpack"`
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
	return nil
}

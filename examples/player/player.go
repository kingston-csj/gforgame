package player

import "io/github/gforgame/db"

type Player struct {
	db.BaseEntity
	Name  string `gorm:"player's name"`
	Level uint   `gorm:"player's' level"`
}

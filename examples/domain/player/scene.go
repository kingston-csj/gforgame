package player

import (
	"io/github/gforgame/db"

	"gorm.io/gorm"
)

type Scene struct {
	db.BaseEntity
	Data string `gorm:"type:text"`
}

func (f *Scene) BeforeSave(tx *gorm.DB) error {
	return nil
}

func (f *Scene) AfterFind(tx *gorm.DB) error {
	return nil
}


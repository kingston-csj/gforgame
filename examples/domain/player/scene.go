package player

import (
	"io/github/gforgame/persist"

	"gorm.io/gorm"
)

type Scene struct {
	persist.BaseEntity
	Data string `gorm:"type:text"`
}

func (f *Scene) BeforeSave(tx *gorm.DB) error {
	return nil
}

func (f *Scene) AfterFind(tx *gorm.DB) error {
	return nil
}


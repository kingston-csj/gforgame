package player

import (
	"github.com/forfun/gforgame/persist"

	"gorm.io/gorm"
)

type Scene struct {
	persist.BaseEntity
	Data string `gorm:"type:text"`
}

func (f *Scene) BeforePersist() error {
	return nil
}

func (f *Scene) AfterLoad() error {
	return nil
}

func (f *Scene) BeforeSave(tx *gorm.DB) error {
	return f.BeforePersist()
}

func (f *Scene) AfterFind(tx *gorm.DB) error {
	return f.AfterLoad()
}

func (f *Scene) SnapshotEntity() (persist.Entity, error) {
	return &Scene{
		BaseEntity: f.BaseEntity,
		Data:       f.Data,
	}, nil
}

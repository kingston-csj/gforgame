package system

import (
	"github.com/forfun/gforgame/persist"

	"gorm.io/gorm"
)

type SystemParameterEnt struct {
	persist.BaseEntity
	Data string `gorm:"column:data"`
}

func (s *SystemParameterEnt) GetData() string {
	return s.Data
}

func (s *SystemParameterEnt) SetData(data string) {
	s.Data = data
}

func (s *SystemParameterEnt) GetId() string {
	return s.Id
}

func (s *SystemParameterEnt) BeforePersist() error {
	return nil
}

func (s *SystemParameterEnt) AfterLoad() error {
	return nil
}

func (s *SystemParameterEnt) BeforeSave(tx *gorm.DB) error {
	return s.BeforePersist()
}

func (s *SystemParameterEnt) AfterFind(tx *gorm.DB) error {
	return s.AfterLoad()
}

func (s *SystemParameterEnt) SnapshotEntity() (persist.Entity, error) {
	return &SystemParameterEnt{
		BaseEntity: s.BaseEntity,
		Data:       s.Data,
	}, nil
}

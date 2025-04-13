package system

import (
	"io/github/gforgame/db"

	"gorm.io/gorm"
)

type SystemParameterEnt struct {
	db.BaseEntity
	Data string `gorm:"column:data"`
}

func (s *SystemParameterEnt) GetData() string {
	return s.Data
}

func (s *SystemParameterEnt) SetData(data string) {
	s.Data = data
}

func (s *SystemParameterEnt) GetID() string {
	return s.Id
}

func (s *SystemParameterEnt) BeforeSave(tx *gorm.DB) error {
	return nil
}

func (s *SystemParameterEnt) AfterFind(tx *gorm.DB) error {
	return nil
}

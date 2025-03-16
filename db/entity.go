package db

import "gorm.io/gorm"

type Entity interface {
	GetId() string
	IsDeleted() bool
	SetDeleted()

	// 在保存之前调用 直接使用 gorm 的钩子
	BeforeSave(tx *gorm.DB) error
	// 在查询之后调用 直接使用 gorm 的钩子
	AfterFind(tx *gorm.DB) error
}

type BaseEntity struct {
	Id     string `json:"id"`
	Delete bool
}

func (s *BaseEntity) GetId() string {
	return s.Id
}

func (s *BaseEntity) IsDeleted() bool {
	return s.Delete
}

func (s *BaseEntity) SetDeleted() {
	s.Delete = true
}

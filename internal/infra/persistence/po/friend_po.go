package po

import (
	"encoding/json"

	frienddomain "github.com/forfun/gforgame/internal/domain/friend"
	"github.com/forfun/gforgame/persist"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type FriendPO struct {
	persist.BaseEntity

	FriendJson string `gorm:"column:friend;type:longtext"`
	ApplyJson  string `gorm:"column:apply;type:longtext"`
}

func (p *FriendPO) TableName() string {
	return schema.NamingStrategy{}.TableName("Friend")
}

func (p *FriendPO) GetKey() string {
	return "Friend_" + p.Id + "_" + p.Id
}

func (p *FriendPO) BeforePersist() error {
	return nil
}

func (p *FriendPO) AfterLoad() error {
	return nil
}

func (p *FriendPO) BeforeSave(tx *gorm.DB) error {
	return p.BeforePersist()
}

func (p *FriendPO) AfterFind(tx *gorm.DB) error {
	return p.AfterLoad()
}

func NewFriendPOFromDomain(f *frienddomain.Friend) (*FriendPO, error) {
	result := &FriendPO{
		BaseEntity: f.BaseEntity,
	}
	var err error
	jsonData, err := json.Marshal(f.Friends)
	if err != nil {
		return nil, err
	}
	result.FriendJson = string(jsonData)
	jsonData, err = json.Marshal(f.Applies)
	if err != nil {
		return nil, err
	}
	result.ApplyJson = string(jsonData)
	return result, nil
}

func (p *FriendPO) ToDomain() (*frienddomain.Friend, error) {
	f := &frienddomain.Friend{
		BaseEntity: p.BaseEntity,
	}
	var err error
	if err = json.Unmarshal([]byte(p.FriendJson), &f.Friends); err != nil {
		return nil, err
	}
	if err = json.Unmarshal([]byte(p.ApplyJson), &f.Applies); err != nil {
		return nil, err
	}
	if err = f.AfterLoad(); err != nil {
		return nil, err
	}
	return f, nil
}

func (p *FriendPO) SnapshotEntity() (persist.Entity, error) {
	cp := *p
	return &cp, nil
}

package player

import (
	"encoding/json"

	"io/github/gforgame/db"

	"gorm.io/gorm"
)

type Friend struct {
	db.BaseEntity
	Id         string                      `gorm:"player_id"`
	FriendJson string                      `gorm:"friend_json"`
	ApplyJson  string                      `gorm:"apply_json"`
	Friends    map[string]bool             `gorm:"-"`
	Applies    map[string]*FriendApplyItem `gorm:"-"`
}

func (f *Friend) BeforeSave(tx *gorm.DB) error {
	jsonData, err := json.Marshal(f.Friends)
	if err != nil {
		return err
	}
	f.FriendJson = string(jsonData)
	jsonData, err = json.Marshal(f.Applies)
	if err != nil {
		return err
	}
	f.ApplyJson = string(jsonData)
	return nil
}

func (f *Friend) AfterFind(tx *gorm.DB) error {
	json.Unmarshal([]byte(f.FriendJson), &f.Friends)
	json.Unmarshal([]byte(f.ApplyJson), &f.Applies)
	return nil
}

func (f *Friend) IsFriend(playerId string) bool {
	_, ok := f.Friends[playerId]
	return ok
}

func (f *Friend) AddFriend(playerId string) {
	f.Friends[playerId] = true
}

func (f *Friend) RemoveFriend(playerId string) {
	delete(f.Friends, playerId)
}

func (f *Friend) AddApply(apply *FriendApplyItem) {
	f.Applies[apply.Id] = apply
}

func (f *Friend) HasApplied(playerId string) bool {
	_, ok := f.Applies[playerId]
	return ok
}

func (f *Friend) GetApply(playerId string) *FriendApplyItem {
	return f.Applies[playerId]
}

func (f *Friend) GetAllApply() []*FriendApplyItem {
	applies := make([]*FriendApplyItem, 0, len(f.Applies))
	for _, apply := range f.Applies {
		applies = append(applies, apply)
	}
	return applies
}

// 清除申请记录
func (f *Friend) ClearApply(id1 string, id2 string) {
	for _, apply := range f.Applies {
		if apply.FromId == id1 && apply.TargetId == id2 {
			delete(f.Applies, apply.Id)
		}
		if apply.FromId == id2 && apply.TargetId == id1 {
			delete(f.Applies, apply.Id)
		}
	}
}

package friend

import (
	"sync"

	"github.com/forfun/gforgame/persist"

	"gorm.io/gorm"
)

type Friend struct {
	persist.BaseEntity

	mu      sync.RWMutex
	Friends map[string]bool             `gorm:"-"`
	Applies map[string]*FriendApplyItem `gorm:"-"`
}

func (f *Friend) AfterLoad() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.Friends == nil {
		f.Friends = make(map[string]bool)
	}
	if f.Applies == nil {
		f.Applies = make(map[string]*FriendApplyItem)
	}

	return nil
}

func (f *Friend) AfterFind(tx *gorm.DB) error {
	return f.AfterLoad()
}

func (f *Friend) IsFriend(playerId string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	_, ok := f.Friends[playerId]
	return ok
}

func (f *Friend) AddFriend(playerId string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.Friends[playerId] = true
}

func (f *Friend) RemoveFriend(playerId string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.Friends, playerId)
}

func (f *Friend) AddApply(apply *FriendApplyItem) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.Applies[apply.Id] = apply
}

func (f *Friend) HasApplied(playerId string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	for _, apply := range f.Applies {
		if apply.FromId == f.Id && apply.TargetId == playerId && apply.Status == 0 {
			return true
		}
	}
	return false
}

func (f *Friend) GetApply(playerId string) *FriendApplyItem {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.Applies[playerId]
}

func (f *Friend) GetAllApply() []*FriendApplyItem {
	f.mu.RLock()
	defer f.mu.RUnlock()
	applies := make([]*FriendApplyItem, 0, len(f.Applies))
	for _, apply := range f.Applies {
		applies = append(applies, apply)
	}
	return applies
}

func (f *Friend) ClearApply(id1 string, id2 string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, apply := range f.Applies {
		if apply.FromId == id1 && apply.TargetId == id2 {
			delete(f.Applies, apply.Id)
		}
		if apply.FromId == id2 && apply.TargetId == id1 {
			delete(f.Applies, apply.Id)
		}
	}
}

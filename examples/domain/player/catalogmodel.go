package player

import (
	"io/github/gforgame/container/set"
)

type CatalogModel struct {
	//已解锁id列表
	UnlockIds *set.Set[int32]
	// 已领取id列表
	ReceivedIds *set.Set[int32]
}

func (c *CatalogModel) AfterLoad() {
	if c.UnlockIds == nil {
		c.UnlockIds = set.NewSet[int32]()
	}
	if c.ReceivedIds == nil {
		c.ReceivedIds = set.NewSet[int32]()
	}
}

func (c *CatalogModel) AddUnlock(id int32) bool {
	if c.ReceivedIds.Contains(id) {
		return false
	}
	if !c.UnlockIds.Contains(id) {
		c.UnlockIds.Add(id)
		return true
	}
	return false
}

func (c *CatalogModel) AddReceived(id int32)  {
	if !c.ReceivedIds.Contains(id) {
		c.ReceivedIds.Add(id)
		c.UnlockIds.Remove(id)
	}
}

// 是否可以领取
func (c *CatalogModel) CanReceived(id int32) bool {
	if c.ReceivedIds.Contains(id) {
		return false
	}
	return c.UnlockIds.Contains(id)
}

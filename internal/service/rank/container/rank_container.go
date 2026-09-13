package container

import (
	"github.com/forfun/gforgame/actor"
)

type AddCmd struct {
	Key   any
	Value any
}
type AddResp struct{}

type RemoveCmd struct {
	Key any
}
type RemoveResp struct{}

type UpdateCmd struct {
	Key   any
	Value any
}
type UpdateResp struct{}

type GetCmd struct {
	Key any
}
type GetResp struct {
	Value any
}

type GetItemsCmd struct{}
type GetItemsResp struct {
	Items []RankEntry
}

type ContainsCmd struct {
	Key any
}
type ContainsResp struct {
	Exists bool
}

type SizeCmd struct{}
type SizeResp struct {
	Size int
}

type ConcurrentRankContainer struct {
	ref *actor.ActorRef
}

func NewConcurrentRankContainer(ref *actor.ActorRef) *ConcurrentRankContainer {
	return &ConcurrentRankContainer{ref: ref}
}

func (c *ConcurrentRankContainer) Add(key, value any) {
	_, _ = c.ref.Ask(&AddCmd{Key: key, Value: value})
}

func (c *ConcurrentRankContainer) Remove(key any) {
	_, _ = c.ref.Ask(&RemoveCmd{Key: key})
}

func (c *ConcurrentRankContainer) Update(key, value any) {
	_, _ = c.ref.Ask(&UpdateCmd{Key: key, Value: value})
}

func (c *ConcurrentRankContainer) Get(key any) any {
	respRaw, err := c.ref.Ask(&GetCmd{Key: key})
	if err != nil {
		return nil
	}
	resp, _ := respRaw.(*GetResp)
	if resp == nil {
		return nil
	}
	return resp.Value
}

func (c *ConcurrentRankContainer) GetItems() []RankEntry {
	respRaw, err := c.ref.Ask(&GetItemsCmd{})
	if err != nil {
		return nil
	}
	resp, _ := respRaw.(*GetItemsResp)
	if resp == nil {
		return nil
	}
	return resp.Items
}

func (c *ConcurrentRankContainer) Contains(key any) bool {
	respRaw, err := c.ref.Ask(&ContainsCmd{Key: key})
	if err != nil {
		return false
	}
	resp, _ := respRaw.(*ContainsResp)
	if resp == nil {
		return false
	}
	return resp.Exists
}

func (c *ConcurrentRankContainer) RankSize() int {
	respRaw, err := c.ref.Ask(&SizeCmd{})
	if err != nil {
		return 0
	}
	resp, _ := respRaw.(*SizeResp)
	if resp == nil {
		return 0
	}
	return resp.Size
}

func (c *ConcurrentRankContainer) Stop() {
	c.ref.Stop()
}

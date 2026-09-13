package rank

import (
	"fmt"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/forfun/gforgame/actor"
	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/internal/service/rank/container"
	"github.com/forfun/gforgame/internal/service/rank/model"
)

type RankActor struct {
	*actor.BaseActor
	rankType RankType
	ranks    *treemap.Map
	capacity int
}

func NewRankActor(rankType RankType, capacity int) *RankActor {
	a := &RankActor{
		rankType: rankType,
		capacity: capacity,
		BaseActor: actor.NewBaseActor(),
	}
	actor.RegisterCmd(a.BaseActor, a.handleAdd)
	actor.RegisterCmd(a.BaseActor, a.handleRemove)
	actor.RegisterCmd(a.BaseActor, a.handleUpdate)
	actor.RegisterCmd(a.BaseActor, a.handleGet)
	actor.RegisterCmd(a.BaseActor, a.handleGetItems)
	actor.RegisterCmd(a.BaseActor, a.handleContains)
	actor.RegisterCmd(a.BaseActor, a.handleSize)
	return a
}

func (a *RankActor) OnStart() {
	a.ranks = treemap.NewWith(model.CompareRank)
}


func (a *RankActor) OnStop() {
	logger.Info(fmt.Sprintf("RankActor stop, rankType=%d", a.rankType))
}

func (a *RankActor) handleAdd(cmd *container.AddCmd) actor.Message {
	a.ranks.Put(cmd.Value, cmd.Key)
	if a.ranks.Size() > a.capacity {
		it := a.ranks.Iterator()
		it.Last()
		a.ranks.Remove(it.Key())
	}
	return &container.AddResp{}
}

func (a *RankActor) handleRemove(cmd *container.RemoveCmd) actor.Message {
	it := a.ranks.Iterator()
	for it.Next() {
		if it.Value() == cmd.Key {
			a.ranks.Remove(it.Key())
			break
		}
	}
	return &container.RemoveResp{}
}

func (a *RankActor) handleUpdate(cmd *container.UpdateCmd) actor.Message {
	it := a.ranks.Iterator()
	for it.Next() {
		if it.Value() == cmd.Key {
			a.ranks.Remove(it.Key())
			break
		}
	}
	a.ranks.Put(cmd.Value, cmd.Key)
	if a.ranks.Size() > a.capacity {
		it := a.ranks.Iterator()
		it.Last()
		a.ranks.Remove(it.Key())
	}
	return &container.UpdateResp{}
}

func (a *RankActor) handleGet(cmd *container.GetCmd) actor.Message {
	it := a.ranks.Iterator()
	for it.Next() {
		if it.Value() == cmd.Key {
			return &container.GetResp{Value: it.Key()}
		}
	}
	return &container.GetResp{Value: nil}
}

func (a *RankActor) handleGetItems(cmd *container.GetItemsCmd) actor.Message {
	items := make([]container.RankEntry, 0, a.ranks.Size())
	it := a.ranks.Iterator()
	for it.Next() {
		items = append(items, container.RankEntry{
			Key:   it.Value(),
			Value: it.Key().(model.BaseRank),
		})
	}
	return &container.GetItemsResp{Items: items}
}

func (a *RankActor) handleContains(cmd *container.ContainsCmd) actor.Message {
	exists := false
	it := a.ranks.Iterator()
	for it.Next() {
		if it.Value() == cmd.Key {
			exists = true
			break
		}
	}
	return &container.ContainsResp{Exists: exists}
}

func (a *RankActor) handleSize(cmd *container.SizeCmd) actor.Message {
	return &container.SizeResp{Size: a.ranks.Size()}
}

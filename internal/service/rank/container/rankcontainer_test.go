package container

import (
	"strings"
	"testing"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/forfun/gforgame/actor"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/rank/model"
)

type TestRank struct {
	Id    string
	Score int64
}

func (r TestRank) GetId() string {
	return r.Id
}

func (r TestRank) AsVo() protos.RankInfo {
	return protos.RankInfo{
		Id:    r.Id,
		Value: r.Score,
	}
}

func (r TestRank) CompareTo(other model.BaseRank) int {
	o := other.(TestRank)
	if r.Score > o.Score {
		return 1
	}
	if r.Score < o.Score {
		return -1
	}

	return strings.Compare(r.Id, o.Id)
}

type testRankActor struct {
	ranks    *treemap.Map
	capacity int
}

func newTestRankActor(capacity int) *testRankActor {
	return &testRankActor{capacity: capacity}
}

func (a *testRankActor) OnStart() {
	a.ranks = treemap.NewWith(model.CompareRank)
}

func (a *testRankActor) OnMessage(raw actor.Message) actor.Message {
	switch cmd := raw.(type) {
	case *AddCmd:
		a.ranks.Put(cmd.Value, cmd.Key)
		if a.ranks.Size() > a.capacity {
			it := a.ranks.Iterator()
			it.Last()
			a.ranks.Remove(it.Key())
		}
		return &AddResp{}
	case *RemoveCmd:
		it := a.ranks.Iterator()
		for it.Next() {
			if it.Value() == cmd.Key {
				a.ranks.Remove(it.Key())
				break
			}
		}
		return &RemoveResp{}
	case *UpdateCmd:
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
		return &UpdateResp{}
	case *GetCmd:
		it := a.ranks.Iterator()
		for it.Next() {
			if it.Value() == cmd.Key {
				return &GetResp{Value: it.Key()}
			}
		}
		return &GetResp{Value: nil}
	case *GetItemsCmd:
		items := make([]RankEntry, 0, a.ranks.Size())
		it := a.ranks.Iterator()
		for it.Next() {
			items = append(items, RankEntry{
				Key:   it.Value(),
				Value: it.Key().(model.BaseRank),
			})
		}
		return &GetItemsResp{Items: items}
	case *ContainsCmd:
		exists := false
		it := a.ranks.Iterator()
		for it.Next() {
			if it.Value() == cmd.Key {
				exists = true
				break
			}
		}
		return &ContainsResp{Exists: exists}
	case *SizeCmd:
		return &SizeResp{Size: a.ranks.Size()}
	}
	return nil
}

func (a *testRankActor) OnStop() {}

func newTestRankContainer(capacity int) *ConcurrentRankContainer {
	sys := actor.NewActorSystem()
	path := actor.NewActorPath("test", "rank", "test_rank")
	ref := sys.Spawn(path, newTestRankActor(capacity))
	return NewConcurrentRankContainer(ref)
}

func TestPlayerLevelRankHandler_UpdateRank(t *testing.T) {

}

func TestRankContainer_RemoveOrdering(t *testing.T) {
	container := newTestRankContainer(7)

	testData := []struct {
		key   string
		score int64
	}{
		{"player1", 100},
		{"player2", 200},
		{"player3", 200},
		{"player4", 400},
		{"player5", 500},
		{"player6", 100},
		{"player7", 700},
	}

	for _, data := range testData {
		container.Add(data.key, TestRank{Id: data.key, Score: data.score})
	}

	items := container.GetItems()
	if len(items) != 7 {
		t.Errorf("Expected 7 items, got %d", len(items))
	}

	expectedInitialOrder := []struct {
		id    string
		score int64
	}{
		{"player7", 700},
		{"player5", 500},
		{"player4", 400},
		{"player3", 200},
		{"player2", 200},
		{"player6", 100},
		{"player1", 100},
	}

	for i, item := range items {
		if item.Value.(TestRank).Score != expectedInitialOrder[i].score {
			t.Errorf("Expected score %d at index %d, got %d",
				expectedInitialOrder[i].score, i, item.Value.(TestRank).Score)
		}
		if item.Value.(TestRank).Id != expectedInitialOrder[i].id {
			t.Errorf("Expected id %s at index %d, got %s",
				expectedInitialOrder[i].id, i, item.Value.(TestRank).Id)
		}
	}

	container.Remove("player2")

	items = container.GetItems()
	if len(items) != 6 {
		t.Errorf("Expected 6 items after removal, got %d", len(items))
	}

	expectedAfterRemoval := []struct {
		id    string
		score int64
	}{
		{"player7", 700},
		{"player5", 500},
		{"player4", 400},
		{"player3", 200},
		{"player6", 100},
		{"player1", 100},
	}

	for i, item := range items {
		if item.Value.(TestRank).Score != expectedAfterRemoval[i].score {
			t.Errorf("Expected score %d at index %d, got %d",
				expectedAfterRemoval[i].score, i, item.Value.(TestRank).Score)
		}
		if item.Value.(TestRank).Id != expectedAfterRemoval[i].id {
			t.Errorf("Expected id %s at index %d, got %s",
				expectedAfterRemoval[i].id, i, item.Value.(TestRank).Id)
		}
	}
}

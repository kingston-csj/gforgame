package container

import (
	"strings"
	"testing"

	"io/github/gforgame/examples/rank/model"
	"io/github/gforgame/protos"
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

func TestPlayerLevelRankHandler_UpdateRank(t *testing.T) {

}

func TestRankContainer_RemoveOrdering(t *testing.T) {
	// 创建一个容量为7的排行榜容器
	container := NewConcurrentRankContainer(7)

	// 添加一些测试数据
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

	// 添加数据
	for _, data := range testData {
		container.Add(data.key, TestRank{Id: data.key, Score: data.score})
	}

	// 验证初始顺序
	items := container.GetItems()
	if len(items) != 7 {
		t.Errorf("Expected 7 items, got %d", len(items))
	}

	// 验证初始顺序是否正确（应该是从大到小，同分按ID排序）
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

	// 删除中间的元素（player2，分数200）
	container.Remove("player2")

	// 获取删除后的数据
	items = container.GetItems()
	if len(items) != 6 {
		t.Errorf("Expected 6 items after removal, got %d", len(items))
	}

	// 验证删除后的顺序是否正确
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

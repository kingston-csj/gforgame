package handler

import (
	"io/github/gforgame/examples/rank/container"
	"io/github/gforgame/examples/rank/model"
	"math"
)

type RankHandler interface {
	Init()
	QueryRanks(start int, end int) []model.BaseRank
	QueryRankOrder(key any) int
	UpdateRank(rank model.BaseRank)
}

// BaseRankHandler 排行榜处理器
type BaseRankHandler struct {
	rankContainer *container.ConcurrentRankContainer
}

func (h *BaseRankHandler) QueryRanks(start int, end int) []model.BaseRank {
	maxSize := math.Min(float64(end), float64(h.rankContainer.RankSize()))

	rankEntries := make([]model.BaseRank, 0, int(maxSize))

	index := 0
	for _, item := range h.rankContainer.GetItems() {
		index++
		if index < start {
			continue
		}
		if index >= int(maxSize) {
			break
		}
		rankEntries = append(rankEntries, item.Value.(model.BaseRank))
	}

	return rankEntries
}

func (h *BaseRankHandler) QueryRankOrder(key any) int {
	index := 1
	for _, item := range h.rankContainer.GetItems() {
		if item.Key == key {
			return index
		}
		index++
	}
	return index
}

func (h *BaseRankHandler) UpdateRank(rank model.BaseRank) {
	h.rankContainer.Update(rank.GetId(), rank)
}

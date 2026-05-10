package container

import "github.com/forfun/gforgame/internal/service/rank/model"

// RankEntry 排行榜条目
type RankEntry struct {
	Key   any
	Value model.BaseRank
}

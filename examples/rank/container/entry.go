package container

import "io/github/gforgame/examples/rank/model"

// RankEntry 排行榜条目
type RankEntry struct {
	Key   any
	Value model.BaseRank
}

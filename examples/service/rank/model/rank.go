package model

import (
	"io/github/gforgame/protos"
	"strings"
)

type BaseRank interface {
	GetId() string
	AsVo() protos.RankInfo
	CompareTo(other BaseRank) int // <0: less, 0: equal, >0: greater
}

// CompareRank 比较两个BaseRank值
// 用于TreeMap的排序
func CompareRank(a, b interface{}) int {
	rankA := a.(BaseRank)
	rankB := b.(BaseRank)
	// 取反比较结果,使得大的值排在前面
	return -rankA.CompareTo(rankB)
}

type PlayerLevelRank struct {
	Id    string
	Level int32
}

func (r *PlayerLevelRank) GetId() string {
	return r.Id
}

func (r *PlayerLevelRank) AsVo() protos.RankInfo {
	return protos.RankInfo{
		Id:          r.Id,
		Order:       0,
		Value:       int64(r.Level),
		SecondValue: 0,
		ExtraInfo:   "",
	}
}

func (r *PlayerLevelRank) CompareTo(other BaseRank) int {
	o, ok := other.(*PlayerLevelRank)
	if !ok {
		return 0
	}
	if r.Id == o.Id {
		return 0
	}
	if r.Level != o.Level {
		if o.Level > r.Level {
			return -1
		}
		return 1
	}

	return strings.Compare(r.Id, o.Id)
}

type PlayerFightingRank struct {
	Id       string
	Fighting int32
}

func (r *PlayerFightingRank) GetId() string {
	return r.Id
}

func (r *PlayerFightingRank) AsVo() protos.RankInfo {
	return protos.RankInfo{
		Id:          r.Id,
		Order:       0,
		Value:       int64(r.Fighting),
		SecondValue: 0,
		ExtraInfo:   "",
	}
}

func (r *PlayerFightingRank) CompareTo(other BaseRank) int {
	o, ok := other.(*PlayerFightingRank)
	if !ok {
		return 0
	}
	if r.Id == o.Id {
		return 0
	}
	if r.Fighting != o.Fighting {
		if o.Fighting > r.Fighting {
			return -1
		}
		return 1
	}

	return strings.Compare(r.Id, o.Id)
}

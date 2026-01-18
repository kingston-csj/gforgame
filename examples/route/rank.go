package route

import (
	"io/github/gforgame/examples/service/rank"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
)

type RankRoute struct {
	network.Base
	service *rank.RankService
}

func NewRankRoute() *RankRoute {
	return &RankRoute{
		service: rank.GetRankService(),
	}
}

func (c *RankRoute) Init() {
}

func (c *RankRoute) ReqRankQuery(s *network.Session, index int32, msg *protos.ReqRankQuery) *protos.ResRankQuery {
	end := int(msg.Start) + int(msg.PageSize)
	records := c.service.QueryRank(rank.RankType(msg.Type), int(msg.Start), end)
	return &protos.ResRankQuery{
		Type:     msg.Type,
		Records:  records,
		MyRecord: protos.RankInfo{},
	}
}

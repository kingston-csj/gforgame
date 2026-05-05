package route

import (
	"github.com/forfun/gforgame/examples/protos"
	"github.com/forfun/gforgame/examples/service/rank"
	"github.com/forfun/gforgame/network"
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
	records := c.service.QueryRanks(rank.RankType(msg.Type), int(msg.Start), end)
	return &protos.ResRankQuery{
		Type:     msg.Type,
		Records:  records,
		MyRecord: protos.RankInfo{},
	}
}

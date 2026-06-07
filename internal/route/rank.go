package route

import (
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/rank"
)

type RankRoute struct {
	service *rank.RankService
}

func NewRankRoute(service *rank.RankService) *RankRoute {
	return &RankRoute{
		service: service,
	}
}

func (c *RankRoute) ReqRankQuery(playerId string, index int32, msg *protos.ReqRankQuery) *protos.ResRankQuery {
	end := int(msg.Start) + int(msg.PageSize)
	records := c.service.QueryRanks(rank.RankType(msg.Type), int(msg.Start), end)
	return &protos.ResRankQuery{
		Type:     msg.Type,
		Records:  records,
		MyRecord: protos.RankInfo{},
	}
}

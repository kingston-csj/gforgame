package rank

import (
	"io/github/gforgame/examples/service/player"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
)

type RankController struct {
	network.Base
}

func NewRankController() *RankController {
	return &RankController{}
}

func (c *RankController) Init() {
	// network.RegisterMessage(protos.CmdRankReqQuery, &protos.ReqRankQuery{})
	// network.RegisterMessage(protos.CmdRankResQuery, &protos.ResRankQuery{})
}

func (c *RankController) ReqRankQuery(s *network.Session, index int, msg *protos.ReqRankQuery) *protos.ResRankQuery {
	rankService := GetRankService()
	end := int(msg.Start) + int(msg.PageSize)
	records := rankService.QueryRank(RankType(msg.Type), int(msg.Start), end)

	rankInfos := make([]protos.RankInfo, 0)
	order := msg.Start
	for _, record := range records {
		vo := record.AsVo()
		vo.Name = player.GetPlayerService().GetPlayerProfileById(record.GetId()).Name
		vo.Order = order
		rankInfos = append(rankInfos, vo)
		order++
	}

	return &protos.ResRankQuery{
		Type:     msg.Type,
		Records:  rankInfos,
		MyRecord: protos.RankInfo{},
	}
}

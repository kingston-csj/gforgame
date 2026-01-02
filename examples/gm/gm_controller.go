package gm

import (
	"strings"

	"io/github/gforgame/examples/constants"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/item"
	"io/github/gforgame/examples/reward"
	playerService "io/github/gforgame/examples/service/player"

	"io/github/gforgame/util"

	"io/github/gforgame/network"
	"io/github/gforgame/protos"
)

type GmController struct {
	network.Base
}

func NewGmController() *GmController {
	return &GmController{}
}

func (ps *GmController) Init() {
}

func (ps *GmController) ReqAction(s *network.Session, index int, msg *protos.ReqGmAction) interface{} {
	topic := msg.Topic
	params := msg.Params
	player := network.GetPlayerBySession(s)
	switch topic {
	case "add_item":
		itemParams := strings.Split(params, "=")
		itemId, err := util.StringToInt32(itemParams[0])
		if err != nil {
			return &protos.ResGmAction{Code: constants.I18N_COMMON_ILLEGAL_PARAMS}
		}
		itemNum, err := util.StringToInt32(itemParams[1])
		if err != nil {
			return &protos.ResGmAction{Code: constants.I18N_COMMON_ILLEGAL_PARAMS}
		}

		item.GetItemService().AddByModelId(player.(*playerdomain.Player), itemId, itemNum)
	case "add_diamond":
		count, _ := util.StringToInt32(params)
		reward := &reward.CurrencyReward{
			Currency:   "diamond",
			Amount: count,
		}
		reward.Reward(player.(*playerdomain.Player))
	case "add_gold":
		count, _ := util.StringToInt32(params)
		reward := &reward.CurrencyReward{
			Currency:   "gold",
			Amount: count,
		}
		reward.Reward(player.(*playerdomain.Player))
	}
	playerService.GetPlayerService().SavePlayer(player.(*playerdomain.Player))

	return &protos.ResGmAction{Code: 0}
}

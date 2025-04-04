package gm

import (
	"strings"

	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/item"
	playerService "io/github/gforgame/examples/player"

	"io/github/gforgame/examples/utils"
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
	network.RegisterMessage(protos.CmdGmReqAction, &protos.ReqGmAction{})
	network.RegisterMessage(protos.CmdGmResAction, &protos.ResGmAction{})
}

func (ps *GmController) ReqAction(s *network.Session, index int, msg *protos.ReqGmAction) interface{} {
	topic := msg.Topic
	params := msg.Params

	switch topic {
	case "add_item":
		itemParams := strings.Split(params, "=")
		itemId, err := utils.StringToInt32(itemParams[0])
		if err != nil {
			return &protos.ResGmAction{Code: constants.COMMON_ILLEGAL_PARAMS}
		}
		itemNum, err := utils.StringToInt32(itemParams[1])
		if err != nil {
			return &protos.ResGmAction{Code: constants.COMMON_ILLEGAL_PARAMS}
		}
		player := context.SessionManager.GetPlayerBySession(s)
		item.GetItemService().AddByModelId(player.(*playerdomain.Player), itemId, itemNum)
	case "add_diamond":
		player := context.SessionManager.GetPlayerBySession(s)
		count, _ := utils.StringToInt32(params)
		player.(*playerdomain.Player).Purse.AddDiamond(count)
		playerService.GetPlayerService().SavePlayer(player.(*playerdomain.Player))
		playerService.GetPlayerService().NotifyPurseChange(player.(*playerdomain.Player))
	case "add_gold":
		player := context.SessionManager.GetPlayerBySession(s)
		count, _ := utils.StringToInt32(params)
		player.(*playerdomain.Player).Purse.AddGold(count)
		playerService.GetPlayerService().SavePlayer(player.(*playerdomain.Player))
		playerService.GetPlayerService().NotifyPurseChange(player.(*playerdomain.Player))
	}

	return &protos.ResGmAction{Code: 0}
}

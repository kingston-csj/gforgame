package gm

import (
	"io/github/gforgame/common/i18n"
	"io/github/gforgame/examples/item"
	"io/github/gforgame/examples/player"

	"io/github/gforgame/examples/utils"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
	"strings"
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
			return &protos.ResGmAction{Code: i18n.ErrorIllegalParams}
		}
		itemNum, err := utils.StringToInt32(itemParams[1])
		if err != nil {
			return &protos.ResGmAction{Code: i18n.ErrorIllegalParams}
		}
		player := player.GetSessionManager().GetPlayerBySession(s)
		item.GetInstance().AddByModelId(player, itemId, itemNum)
	}

	return &protos.ResGmAction{Code: 0}
}

package item

import (
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
)

type ItemController struct {
	network.Base
}

func NewItemController() *ItemController {
	return &ItemController{}
}

func (ps *ItemController) Init() {
	network.RegisterMessage(protos.CmdItemResBackpackInfo, &protos.ResBackpackInfo{})
}

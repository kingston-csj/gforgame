package mixture

import (
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
	"time"
)

type MixtureController struct {
	network.Base
}

func NewMixtureController() *MixtureController {
	return &MixtureController{}
}

func (c *MixtureController) ReqHeartBeat(s *network.Session, index int, msg *protos.ReqHeartBeat) *protos.ResHeartBeat{
	return &protos.ResHeartBeat{
		Index: msg.Index,
		Code:  0,
	}
}

func (c *MixtureController) ReqGetServerTime(s *network.Session, index int, msg *protos.ReqGetServerTime) *protos.ResGetServerTime{
	return &protos.ResGetServerTime{
		ServerTime: time.Now().Unix(),
		Code:  0,
	}
}
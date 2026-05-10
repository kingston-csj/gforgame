package route

import (
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/transfer"
	"github.com/forfun/gforgame/network"
)

type TransferRoute struct {
	network.Base
	service *transfer.TransferService
}

func NewTransferRoute() *TransferRoute {
	return &TransferRoute{}
}

func (r *TransferRoute) Init() {
	r.service = transfer.GetTransferService()
}

func (r *TransferRoute) ReqTransferGateToLogic(playerId string, s *network.Session, index int32, msg *protos.TransferGateToLogic) interface{} {
	r.service.TransferGateToLogic(s, msg)
	return nil
}


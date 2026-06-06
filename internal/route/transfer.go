package route

import (
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/transfer"
	"github.com/forfun/gforgame/network"
)

type TransferRoute struct {
	service *transfer.TransferService
}

func NewTransferRoute(service *transfer.TransferService) *TransferRoute {
	return &TransferRoute{
		service: service,
	}
}

func (r *TransferRoute) ReqTransferGateToLogic(playerId string, s *network.Session, index int32, msg *protos.TransferGateToLogic) interface{} {
	r.service.TransferGateToLogic(s, msg)
	return nil
}


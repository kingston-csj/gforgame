package route

import (
	"fmt"

	"github.com/forfun/gforgame/internal/protos"
	gate "github.com/forfun/gforgame/internal/service/gate"
	"github.com/forfun/gforgame/network"
)

type GateRoute struct {
	service  *gate.GateService
	transfer func(network.Session, *protos.TransferGateToLogic) error
}

// SetTransferHandler 接入游戏服现有的协议解码与玩家 Actor 分发流程。
func (r *GateRoute) SetTransferHandler(handler func(network.Session, *protos.TransferGateToLogic) error) {
	r.transfer = handler
}

func (r *GateRoute) RPC_Transfer(session network.Session, msg *protos.TransferGateToLogic) (*protos.GateRPCResult, error) {
	if r.transfer == nil {
		return nil, fmt.Errorf("gate transfer handler is not configured")
	}
	if err := r.transfer(session, msg); err != nil {
		return nil, err
	}
	return &protos.GateRPCResult{}, nil
}

func (r *GateRoute) RPC_PlayerLogout(msg *protos.NotifyPlayerLogoutToGame) *protos.GateRPCResult {
	r.service.Logout(msg)
	return &protos.GateRPCResult{}
}

func (r *GateRoute) RPC_SyncOnlinePlayers(msg *protos.NotifyOnlinePlayerToGame) *protos.GateRPCResult {
	r.service.RefreshOnlinePlayers(msg.PlayerIds)
	return &protos.GateRPCResult{}
}

func NewGateRoute(service *gate.GateService) *GateRoute {
	return &GateRoute{
		service: service,
	}
}

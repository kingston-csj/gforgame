package cross

import (
	"github.com/forfun/gforgame/internal/infra/net"
	"github.com/forfun/gforgame/internal/protos"
)

// GateService 统一承接跨服消息处理。
type GateService struct {
	onlinePlayerRegistry *net.OnlinePlayerRegistry
}

func NewGateService(onlinePlayerRegistry *net.OnlinePlayerRegistry) *GateService {
	return &GateService{
		onlinePlayerRegistry: onlinePlayerRegistry,
	}
}

func (s *GateService) Init() {
}

// Logout 处理网关通知的玩家下线事件，并移除逻辑服在线标记。
func (s *GateService) Logout(req *protos.NotifyPlayerLogoutToGame) {
	if req == nil || req.PlayerId == "" {
		return
	}
	s.onlinePlayerRegistry.RemoveOnlinePlayer(req.PlayerId)
}

func (s *GateService) RefreshOnlinePlayers(playerIds []string) {
	for _, playerId := range playerIds {
		s.onlinePlayerRegistry.AddOnlinePlayer(playerId)
	}
}

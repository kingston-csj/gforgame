package route

import (
	"github.com/forfun/gforgame/internal/context"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/events"
	"github.com/forfun/gforgame/internal/protos"
	player "github.com/forfun/gforgame/internal/service/player"
	"github.com/forfun/gforgame/internal/service/signin"
	"github.com/forfun/gforgame/network"
)

// SignInRoute 绛惧埌璺敱
type SignInRoute struct {
	network.Base
	service *signin.SignInService
	player  *player.PlayerService
}

func NewSignInRoute(service *signin.SignInService, playerService *player.PlayerService) *SignInRoute {
	return &SignInRoute{
		service: service,
		player:  playerService,
	}
}

func (ps *SignInRoute) Init() {
	context.EventBus.Subscribe(events.PlayerLogin, func(data interface{}) {
		ps.service.OnPlayerLogin(data.(*playerdomain.Player))
	})
}

func (ps *SignInRoute) ReqSignIn(s *network.Session, index int32, msg *protos.ReqSignIn) *protos.ResSignIn{
	player := ps.player.GetPlayerBySession(s)
	err := ps.service.SignIn(player)
	if err != nil {
		return &protos.ResSignIn{
			Code: int32(err.Code()),
		}
	}
	return &protos.ResSignIn{}
}

func (ps *SignInRoute) ReqSignInMakeup(s *network.Session, index int32, msg *protos.ReqSignInMakeup) *protos.ResSignInMakeup{
	player := ps.player.GetPlayerBySession(s)
	err := ps.service.SignInMakeUp(player, msg.Day)
	if err != nil {
		return &protos.ResSignInMakeup{
			Code: int32(err.Code()),
		}
	}
	return &protos.ResSignInMakeup{}
}
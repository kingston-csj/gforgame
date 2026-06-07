package route

import (
	"github.com/forfun/gforgame/internal/protos"
	player "github.com/forfun/gforgame/internal/service/player"
	"github.com/forfun/gforgame/internal/service/signin"
)

// SignInRoute 绛惧埌璺敱
type SignInRoute struct {
	service *signin.SignInService
	player  *player.PlayerService
}

func NewSignInRoute(service *signin.SignInService, playerService *player.PlayerService) *SignInRoute {
	return &SignInRoute{
		service: service,
		player:  playerService,
	}
}

func (ps *SignInRoute) ReqSignIn(playerId string, index int32, msg *protos.ReqSignIn) *protos.ResSignIn{
	player := ps.player.GetPlayer(playerId)
	err := ps.service.SignIn(player)
	if err != nil {
		return &protos.ResSignIn{
			Code: int32(err.Code()),
		}
	}
	return &protos.ResSignIn{}
}

func (ps *SignInRoute) ReqSignInMakeup(playerId string, index int32, msg *protos.ReqSignInMakeup) *protos.ResSignInMakeup{
	player := ps.player.GetPlayer(playerId)
	err := ps.service.SignInMakeUp(player, msg.Day)
	if err != nil {
		return &protos.ResSignInMakeup{
			Code: int32(err.Code()),
		}
	}
	return &protos.ResSignInMakeup{}
}
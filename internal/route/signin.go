package route

import (
	playerrepo "github.com/forfun/gforgame/internal/infra/repository/player"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/signin"
)

// SignInRoute 绛惧埌璺敱
type SignInRoute struct {
	service    *signin.SignInService
	playerrepo *playerrepo.PlayerRepository
}

func NewSignInRoute(service *signin.SignInService, playerRepo *playerrepo.PlayerRepository) *SignInRoute {
	return &SignInRoute{
		service:    service,
		playerrepo: playerRepo,
	}
}

func (ps *SignInRoute) ReqSignIn(playerId string, index int32, msg *protos.ReqSignIn) *protos.ResSignIn {
	player := ps.playerrepo.GetPlayer(playerId)
	err := ps.service.SignIn(player)
	if err != nil {
		return &protos.ResSignIn{
			Code: int32(err.Code()),
		}
	}
	return &protos.ResSignIn{}
}

func (ps *SignInRoute) ReqSignInMakeup(playerId string, index int32, msg *protos.ReqSignInMakeup) *protos.ResSignInMakeup {
	player := ps.playerrepo.GetPlayer(playerId)
	err := ps.service.SignInMakeUp(player, msg.Day)
	if err != nil {
		return &protos.ResSignInMakeup{
			Code: int32(err.Code()),
		}
	}
	return &protos.ResSignInMakeup{}
}

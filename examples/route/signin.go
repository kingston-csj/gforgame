package route

import (
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/service/signin"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
)

type SignInRoute struct {
	network.Base
	service *signin.SignInService
}

func NewSignInRoute() *SignInRoute {
	return &SignInRoute{
		service: signin.GetSignInService(),
	}
}

func (ps *SignInRoute) Init() {
	context.EventBus.Subscribe(events.PlayerLogin, func(data interface{}) {
		ps.service.OnPlayerLogin(data.(*playerdomain.Player))
	})
}

func (ps *SignInRoute) ReqSignIn(s *network.Session, index int32, msg *protos.ReqSignIn) *protos.ResSignIn{
	player := network.GetPlayerBySession(s).(*playerdomain.Player)
	err := ps.service.SignIn(player)
	if err != nil {
		return &protos.ResSignIn{
			Code: int32(err.Code()),
		}
	}
	return &protos.ResSignIn{}
}
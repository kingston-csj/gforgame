package route

import (
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"

	heroService "io/github/gforgame/examples/service/hero"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
)

type HeroRoute struct {
	network.Base
	service *heroService.HeroService
}

func NewHeroRoute() *HeroRoute {
	return &HeroRoute{}
}

func (ps *HeroRoute) Init() {
	context.EventBus.Subscribe(events.PlayerLogin, func(data interface{}) {
		ps.service.OnPlayerLogin(data.(*playerdomain.Player))
	})

	context.EventBus.Subscribe(events.PlayerAfterLoad, func(data interface{}) {
		p := data.(*playerdomain.Player)
		for _, h := range p.HeroBox.Heros {
			ps.service.ReCalculateHeroAttr(p, h, false)
		}
	})
}

func (ps *HeroRoute) OnPlayerLogin(player *playerdomain.Player) {
	ps.service.OnPlayerLogin(player)
}

func (ps *HeroRoute) ReqRecruit(s *network.Session, index int32, msg *protos.ReqHeroRecruit) *protos.ResHeroRecruit {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	err, rewards := ps.service.DoRecruit(p, 1, msg.Counter)
	if err != nil {
		return &protos.ResHeroRecruit{
			Code: int32(err.Code()),
		}
	}
	return &protos.ResHeroRecruit{
		RewardVos: rewards,
	}
}

func (ps *HeroRoute) ReqHeroLevelUp(s *network.Session, index int32, msg *protos.ReqHeroLevelUp) *protos.ResHeroLevelUp {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	return ps.service.DoLevelUp(p, msg.HeroId, msg.ToLevel)
}

func (ps *HeroRoute) ReqHeroUpStage(s *network.Session, index int32, msg *protos.ReqHeroUpStage) *protos.ResHeroUpStage {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	return ps.service.DoStageUp(p, msg.HeroId)
}

func (ps *HeroRoute) ReqHeroCombine(s *network.Session, index int32, msg *protos.ReqHeroCombine) *protos.ResHeroCombine {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	return ps.service.DoCombine(p, msg.HeroId)
}

func (ps *HeroRoute) ReqHeroUpFight(s *network.Session, index int32, msg *protos.ReqHeroUpFight) *protos.ResHeroUpFight {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	return ps.service.DoUpFight(p, msg.HeroId)
}

func (ps *HeroRoute) ReqHeroOffFight(s *network.Session, index int32, msg *protos.ReqHeroOffFight) *protos.ResHeroOffFight {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	return ps.service.DoOffFight(p, msg.HeroId)
}

func (ps *HeroRoute) ReqHeroChangePosition(s *network.Session, index int32, msg *protos.ReqHeroChangePosition) *protos.ResHeroChangePosition {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	return ps.service.DoChangePosition(p, msg.HeroId, msg.Position)
}

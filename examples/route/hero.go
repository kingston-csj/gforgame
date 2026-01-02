package route

import (
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"

	heroService "io/github/gforgame/examples/service/hero"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
)

type HeroController struct {
	network.Base
	service *heroService.HeroService
}

func NewHeroController() *HeroController {
	return &HeroController{}
}

func (ps *HeroController) Init() {
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

func (ps *HeroController) OnPlayerLogin(player *playerdomain.Player) {
	ps.service.OnPlayerLogin(player)
}

func (ps *HeroController) ReqRecruit(s *network.Session, index int, msg *protos.ReqHeroRecruit) *protos.ResHeroRecruit {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	return ps.service.DoRecruit(p, msg.Times)
}
 

func (ps *HeroController) ReqHeroLevelUp(s *network.Session, index int, msg *protos.ReqHeroLevelUp) *protos.ResHeroLevelUp {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	return ps.service.DoLevelUp(p, msg.HeroId, msg.ToLevel)
}

func (ps *HeroController) ReqHeroUpStage(s *network.Session, index int, msg *protos.ReqHeroUpStage) *protos.ResHeroUpStage {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	return ps.service.DoStageUp(p, msg.HeroId)
}

func (ps *HeroController) ReqHeroCombine(s *network.Session, index int, msg *protos.ReqHeroCombine) *protos.ResHeroCombine {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	return ps.service.DoCombine(p, msg.HeroId)
}

func (ps *HeroController) ReqHeroUpFight(s *network.Session, index int, msg *protos.ReqHeroUpFight) *protos.ResHeroUpFight {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	return ps.service.DoUpFight(p, msg.HeroId)
}

func (ps *HeroController) ReqHeroOffFight(s *network.Session, index int, msg *protos.ReqHeroOffFight) *protos.ResHeroOffFight {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	return ps.service.DoOffFight(p, msg.HeroId)
}

func (ps *HeroController) ReqHeroChangePosition(s *network.Session, index int, msg *protos.ReqHeroChangePosition) *protos.ResHeroChangePosition {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	return ps.service.DoChangePosition(p, msg.HeroId, msg.Position)
}

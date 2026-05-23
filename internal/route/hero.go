package route

import (
	"github.com/forfun/gforgame/internal/context"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/events"

	"github.com/forfun/gforgame/internal/protos"
	heroService "github.com/forfun/gforgame/internal/service/hero"
	player "github.com/forfun/gforgame/internal/service/player"
	"github.com/forfun/gforgame/network"
)

type HeroRoute struct {
	network.Base
	service *heroService.HeroService
	player  *player.PlayerService
}

func NewHeroRoute(service *heroService.HeroService, playerService *player.PlayerService) *HeroRoute {
	return &HeroRoute{
		service: service,
		player:  playerService,
	}
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
	p := ps.player.GetPlayerBySession(s)
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
	p := ps.player.GetPlayerBySession(s)
	return ps.service.DoLevelUp(p, msg.HeroId, msg.ToLevel)
}

func (ps *HeroRoute) ReqHeroUpStage(s *network.Session, index int32, msg *protos.ReqHeroUpStage) *protos.ResHeroUpStage {
	p := ps.player.GetPlayerBySession(s)
	return ps.service.DoStageUp(p, msg.HeroId)
}

func (ps *HeroRoute) ReqHeroCombine(s *network.Session, index int32, msg *protos.ReqHeroCombine) *protos.ResHeroCombine {
	p := ps.player.GetPlayerBySession(s)
	return ps.service.DoCombine(p, msg.HeroId)
}

func (ps *HeroRoute) ReqHeroUpFight(s *network.Session, index int32, msg *protos.ReqHeroUpFight) *protos.ResHeroUpFight {
	p := ps.player.GetPlayerBySession(s)
	return ps.service.DoUpFight(p, msg.HeroId)
}

func (ps *HeroRoute) ReqHeroOffFight(s *network.Session, index int32, msg *protos.ReqHeroOffFight) *protos.ResHeroOffFight {
	p := ps.player.GetPlayerBySession(s)
	return ps.service.DoOffFight(p, msg.HeroId)
}

func (ps *HeroRoute) ReqHeroChangePosition(s *network.Session, index int32, msg *protos.ReqHeroChangePosition) *protos.ResHeroChangePosition {
	p := ps.player.GetPlayerBySession(s)
	return ps.service.DoChangePosition(p, msg.HeroId, msg.Position)
}

package route

import (
	playerdomain "github.com/forfun/gforgame/internal/domain/player"

	"github.com/forfun/gforgame/internal/protos"
	heroService "github.com/forfun/gforgame/internal/service/hero"
	player "github.com/forfun/gforgame/internal/service/player"
)

type HeroRoute struct {
	service *heroService.HeroService
	player  *player.PlayerService
}

func NewHeroRoute(service *heroService.HeroService, playerService *player.PlayerService) *HeroRoute {
	return &HeroRoute{
		service: service,
		player:  playerService,
	}
}


func (ps *HeroRoute) OnPlayerLogin(player *playerdomain.Player) {
	ps.service.OnPlayerLogin(player)
}

func (ps *HeroRoute) ReqRecruit(playerId string, index int32, msg *protos.ReqHeroRecruit) *protos.ResHeroRecruit {
	p := ps.player.GetPlayer(playerId)
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

func (ps *HeroRoute) ReqHeroLevelUp(playerId string, index int32, msg *protos.ReqHeroLevelUp) *protos.ResHeroLevelUp {
	p := ps.player.GetPlayer(playerId)
	return ps.service.DoLevelUp(p, msg.HeroId, msg.ToLevel)
}

func (ps *HeroRoute) ReqHeroUpStage(playerId string, index int32, msg *protos.ReqHeroUpStage) *protos.ResHeroUpStage {
	p := ps.player.GetPlayer(playerId)
	return ps.service.DoStageUp(p, msg.HeroId)
}

func (ps *HeroRoute) ReqHeroCombine(playerId string, index int32, msg *protos.ReqHeroCombine) *protos.ResHeroCombine {
	p := ps.player.GetPlayer(playerId)
	return ps.service.DoCombine(p, msg.HeroId)
}

func (ps *HeroRoute) ReqHeroUpFight(playerId string, index int32, msg *protos.ReqHeroUpFight) *protos.ResHeroUpFight {
	p := ps.player.GetPlayer(playerId)
	return ps.service.DoUpFight(p, msg.HeroId)
}

func (ps *HeroRoute) ReqHeroOffFight(playerId string, index int32, msg *protos.ReqHeroOffFight) *protos.ResHeroOffFight {
	p := ps.player.GetPlayer(playerId)
	return ps.service.DoOffFight(p, msg.HeroId)
}

func (ps *HeroRoute) ReqHeroChangePosition(playerId string, index int32, msg *protos.ReqHeroChangePosition) *protos.ResHeroChangePosition {
	p := ps.player.GetPlayer(playerId)
	return ps.service.DoChangePosition(p, msg.HeroId, msg.Position)
}

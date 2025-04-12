package hero

import (
	"fmt"
	"strconv"

	"io/github/gforgame/common"
	"io/github/gforgame/examples/camp"
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/consume"
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"

	"io/github/gforgame/examples/io"
	"io/github/gforgame/examples/item"

	"io/github/gforgame/examples/session"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
)

type HeroController struct {
	network.Base
}

func NewHeroController() *HeroController {
	return &HeroController{}
}

func (ps *HeroController) Init() {
	network.RegisterMessage(protos.CmdHeroReqRecruit, &protos.ReqHeroRecruit{})
	network.RegisterMessage(protos.CmdHeroResRecruit, &protos.ResHeroRecruit{})
	network.RegisterMessage(protos.CmdHeroResAllHero, &protos.ResAllHeroInfo{})

	network.RegisterMessage(protos.CmdHeroReqLevelUp, &protos.ReqHeroLevelUp{})
	network.RegisterMessage(protos.CmdHeroResLevelUp, &protos.ResHeroLevelUp{})
	network.RegisterMessage(protos.CmdHeroPushAdd, &protos.PushHeroAdd{})
	network.RegisterMessage(protos.CmdHeroPushAttrChange, &protos.PushHeroAttrChange{})

	network.RegisterMessage(protos.CmdHeroReqUpStage, &protos.ReqHeroUpStage{})
	network.RegisterMessage(protos.CmdHeroResUpStage, &protos.ResHeroUpStage{})

	context.EventBus.Subscribe(events.PlayerLogin, func(data interface{}) {
		ps.OnPlayerLogin(data.(*playerdomain.Player))
	})
}

func (ps *HeroController) OnPlayerLogin(player *playerdomain.Player) {
	resAllHeroInfo := &protos.ResAllHeroInfo{}

	// 普通英雄
	for _, h := range player.HeroBox.Heros {
		GetHeroService().ReCalculateHeroAttr(player, h, false)
		attrInfos := make([]protos.AttrInfo, 0)
		for _, attr := range h.AttrBox.GetAttrs() {
			attrInfos = append(attrInfos, protos.AttrInfo{
				AttrType: string(attr.AttrType),
				Value:    attr.Value,
			})
		}
		resAllHeroInfo.Heros = append(resAllHeroInfo.Heros, &protos.HeroInfo{
			Id:       h.ModelId,
			Level:    h.Level,
			Position: 0,
			Stage:    h.Stage,
			Attrs:    attrInfos,
			Fight:    h.Fight,
		})
	}

	// 主公
	masterId := camp.GetHeroIdByCamp(player.Camp)
	masterAttrInfos := make([]protos.AttrInfo, 0)
	resAllHeroInfo.Heros = append(resAllHeroInfo.Heros, &protos.HeroInfo{
		Id:       masterId,
		Level:    player.Level,
		Position: 0,
		Stage:    0,
		Attrs:    masterAttrInfos,
		Fight:    0,
	})

	io.NotifyPlayer(player, resAllHeroInfo)
}

func (ps *HeroController) ReqRecruit(s *network.Session, index int, msg *protos.ReqHeroRecruit) *protos.ResHeroRecruit {
	rewardInfos := make([]*protos.RewardInfo, 0)

	p := session.GetPlayerBySession(s).(*playerdomain.Player)
	if p.Backpack.GetItemCount(item.RecruitItemId) < int32(msg.Times) {
		return &protos.ResHeroRecruit{
			Code: constants.I18N_ITEM_NOT_ENOUGH,
		}
	}

	p.Backpack.RemoveItem(item.RecruitItemId, msg.Times)

	for i := 0; i < int(msg.Times); i++ {
		heroData := GetHeroService().GetRandomHero()
		// 如果已经拥有该英雄，则转为碎片
		if p.HeroBox.HasHero(heroData.Id) {
			rewardInfos = append(rewardInfos, &protos.RewardInfo{
				Type:  "item",
				Value: fmt.Sprintf("%d=%d", heroData.Item, heroData.Shard),
			})
			item.GetItemService().AddByModelId(p, heroData.Item, heroData.Shard)
		} else {
			rewardInfos = append(rewardInfos, &protos.RewardInfo{
				Type:  "hero",
				Value: strconv.Itoa(int(heroData.Id)),
			})
			p.HeroBox.AddHero(&playerdomain.Hero{
				ModelId: heroData.Id,
				Level:   1,
			})

			GetHeroService().ReCalculateHeroAttr(p, p.HeroBox.GetHero(heroData.Id), true)
		}

	}

	context.EventBus.Publish(events.PlayerEntityChange, p)

	return &protos.ResHeroRecruit{
		Code:        0,
		RewardInfos: rewardInfos,
	}
}

func (ps *HeroController) ReqHeroLevelUp(s *network.Session, index int, msg *protos.ReqHeroLevelUp) *protos.ResHeroLevelUp {
	p := session.GetPlayerBySession(s).(*playerdomain.Player)

	h := p.HeroBox.GetHero(msg.HeroId)
	if h == nil {
		return &protos.ResHeroLevelUp{
			Code: constants.I18N_COMMON_NOT_FOUND,
		}
	}

	if h.Level >= p.Level {
		return &protos.ResHeroLevelUp{
			Code: constants.I18N_HERO_TIP1,
		}
	}

	stageData := GetHeroService().GetHeroStageData(h.Stage)
	if h.Level >= stageData.MaxLevel {
		return &protos.ResHeroLevelUp{
			Code: constants.I18N_HERO_TIP2,
		}
	}

	costGold := GetHeroService().calcTotalUpLevelConsume(h.Level, msg.ToLevel)
	if !p.Purse.IsEnoughGold(costGold) {
		return &protos.ResHeroLevelUp{
			Code: constants.I18N_GOLD_NOT_ENOUGH,
		}
	}

	consume := consume.CurrencyConsume{
		Kind:   "gold",
		Amount: costGold,
	}
	err := consume.Verify(p)
	if err != nil {
		return &protos.ResHeroLevelUp{
			Code: int32(err.(*common.BusinessRequestException).Code()),
		}
	}
	consume.Consume(p)

	h.Level = msg.ToLevel
	GetHeroService().ReCalculateHeroAttr(p, h, true)
	context.EventBus.Publish(events.PlayerEntityChange, p)

	return &protos.ResHeroLevelUp{
		Code: 0,
	}
}

func (ps *HeroController) ReqHeroUpStage(s *network.Session, index int, msg *protos.ReqHeroUpStage) *protos.ResHeroUpStage {
	p := session.GetPlayerBySession(s).(*playerdomain.Player)

	h := p.HeroBox.GetHero(msg.HeroId)
	if h == nil {
		return &protos.ResHeroUpStage{
			Code: constants.I18N_COMMON_NOT_FOUND,
		}
	}

	stageData := GetHeroService().GetHeroStageData(h.Stage)
	if h.Level < stageData.MaxLevel {
		return &protos.ResHeroUpStage{
			Code: constants.I18N_HERO_TIP3,
		}
	}

	costItem := consume.ItemConsume{
		ItemId: constants.GAME_UPSTAGE_ITEM_ID,
		Amount: stageData.Cost,
	}
	err := costItem.Verify(p)
	if err != nil {
		return &protos.ResHeroUpStage{
			Code: int32(err.(*common.BusinessRequestException).Code()),
		}
	}
	costItem.Consume(p)

	h.Stage = h.Stage + 1

	GetHeroService().ReCalculateHeroAttr(p, h, true)
	context.EventBus.Publish(events.PlayerEntityChange, p)

	return &protos.ResHeroUpStage{
		Code: 0,
	}
}

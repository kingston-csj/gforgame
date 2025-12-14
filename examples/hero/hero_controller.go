package hero

import (
	"fmt"
	"slices"
	"strconv"

	"io/github/gforgame/common"
	"io/github/gforgame/examples/camp"
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/config/container"
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/consume"
	"io/github/gforgame/examples/context"
	configdomain "io/github/gforgame/examples/domain/config"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/item"

	"io/github/gforgame/examples/io"

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
	context.EventBus.Subscribe(events.PlayerLogin, func(data interface{}) {
		ps.OnPlayerLogin(data.(*playerdomain.Player))
	})

	context.EventBus.Subscribe(events.PlayerAfterLoad, func(data interface{}) {
		p := data.(*playerdomain.Player)
		for _, h := range p.HeroBox.Heros {
			GetHeroService().ReCalculateHeroAttr(p, h, false)
		}
	})
}

func (ps *HeroController) OnPlayerLogin(player *playerdomain.Player) {
	resAllHeroInfo := &protos.PushAllHeroInfo{}

	// 普通英雄
	for _, h := range player.HeroBox.Heros {
		GetHeroService().ReCalculateHeroAttr(player, h, false)
		attrInfos := make([]protos.AttrInfo, 0)
		for _, attr := range h.AttrBox.GetAttrs() {
			attrInfos = append(attrInfos, protos.AttrInfo{
				AttrType: string(attr.AttrType),
				Value:    int32(attr.Value),
			})
		}
		resAllHeroInfo.Heros = append(resAllHeroInfo.Heros, &protos.HeroInfo{
			Id:       h.ModelId,
			Level:    h.Level,
			Position: h.Position,
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
		Stage:    player.Stage,
		Attrs:    masterAttrInfos,
		Fight:    0,
	})

	io.NotifyPlayer(player, resAllHeroInfo)
}

func (ps *HeroController) ReqRecruit(s *network.Session, index int, msg *protos.ReqHeroRecruit) *protos.ResHeroRecruit {
	rewardInfos := make([]*protos.RewardVo, 0)

	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	if p.Backpack.GetItemCount(item.RecruitItemId) < int32(msg.Times) {
		return &protos.ResHeroRecruit{
			Code: constants.I18N_ITEM_NOT_ENOUGH,
		}
	}

	p.Backpack.ReduceByModelId(item.RecruitItemId, msg.Times)

	for i := 0; i < int(msg.Times); i++ {
		heroData := GetHeroService().GetRandomHero()
		// 如果已经拥有该英雄，则转为碎片
		if p.HeroBox.HasHero(heroData.Id) {
			rewardInfos = append(rewardInfos, &protos.RewardVo{
				Type:  "item",
				Value: fmt.Sprintf("%d=%d", heroData.Item, heroData.Shard),
			})
			item.GetItemService().AddByModelId(p, heroData.Item, heroData.Shard)
		} else {
			rewardInfos = append(rewardInfos, &protos.RewardVo{
				Type:  "hero",
				Value: strconv.Itoa(int(heroData.Id)),
			})
			ps.NewHero(p, heroData.Id)

			GetHeroService().ReCalculateHeroAttr(p, p.HeroBox.GetHero(heroData.Id), true)
		}

	}

	context.EventBus.Publish(events.PlayerEntityChange, p)

	return &protos.ResHeroRecruit{
		Code:        0,
		RewardInfos: rewardInfos,
	}
}

func (ps *HeroController) NewHero(p *playerdomain.Player,heroId int32) {
	p.HeroBox.AddHero(&playerdomain.Hero{
				ModelId: heroId,
				Level:   1,
			})
	context.EventBus.Publish(events.HeroGain, &events.HeroGainEvent{
		Player: p,
		HeroId: heroId,
	})
}

func (ps *HeroController) ReqHeroLevelUp(s *network.Session, index int, msg *protos.ReqHeroLevelUp) *protos.ResHeroLevelUp {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	toLevel := msg.ToLevel
	h := p.HeroBox.GetHero(msg.HeroId)
	if h == nil {
		return &protos.ResHeroLevelUp{
			Code: constants.I18N_COMMON_NOT_FOUND,
		}
	}

	if toLevel > p.Level {
		return &protos.ResHeroLevelUp{
			Code: constants.I18N_HERO_TIP1,
		}
	}
	if toLevel <= h.Level {
		return &protos.ResHeroLevelUp{
			Code: constants.I18N_COMMON_ILLEGAL_PARAMS,
		}
	}

	stageContainer := config.GetSpecificContainer[ container.HeroStageContainer]("herostage")

	stageData := stageContainer.GetRecordByStage(h.Stage)
	if stageData == nil {
		return &protos.ResHeroLevelUp{
			Code: constants.I18N_COMMON_NOT_FOUND,
		}
	}
	if h.Level >= stageData.MaxLevel {
		return &protos.ResHeroLevelUp{
			Code: constants.I18N_HERO_TIP2,
		}
	}

	costGold := GetHeroService().CalcTotalUpLevelConsume(h.Level, toLevel)
	if !p.Purse.IsEnoughGold(costGold) {
		return &protos.ResHeroLevelUp{
			Code: constants.I18N_GOLD_NOT_ENOUGH,
		}
	}

	consume := consume.CurrencyConsume{
		Currency:   "gold",
		Amount: costGold,
	}
	err := consume.Verify(p)
	if err != nil {
		return &protos.ResHeroLevelUp{
			Code: int32(err.(*common.BusinessRequestException).Code()),
		}
	}
	consume.Consume(p)

	h.Level = toLevel
	GetHeroService().ReCalculateHeroAttr(p, h, true)
	context.EventBus.Publish(events.PlayerEntityChange, p)

	return &protos.ResHeroLevelUp{
		Code: 0,
	}
}

func (ps *HeroController) ReqHeroUpStage(s *network.Session, index int, msg *protos.ReqHeroUpStage) *protos.ResHeroUpStage {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)

	h := p.HeroBox.GetHero(msg.HeroId)
	if h == nil {
		return &protos.ResHeroUpStage{
			Code: constants.I18N_COMMON_NOT_FOUND,
		}
	}

	stageContainer := config.GetSpecificContainer[container.HeroStageContainer]("herostage")
	stageData := stageContainer.GetRecordByStage(h.Stage)
	if stageData == nil {
		return &protos.ResHeroUpStage{
			Code: constants.I18N_COMMON_NOT_FOUND,
		}
	}
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

func (ps *HeroController) ReqHeroCombine(s *network.Session, index int, msg *protos.ReqHeroCombine) *protos.ResHeroCombine {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)

	h := p.HeroBox.GetHero(msg.HeroId)
	if h != nil {
		return &protos.ResHeroCombine{
			Code: constants.I18N_HERO_TIP5,
		}
	}
	heroData := config.QueryById[configdomain.HeroData](msg.HeroId)
	itemConsume := consume.ItemConsume{
		ItemId: heroData.Item,
		Amount: heroData.Shard,
	}
	err := itemConsume.Verify(p)
	if err != nil {
		return &protos.ResHeroCombine{
			Code: int32(err.(*common.BusinessRequestException).Code()),
		}
	}
	itemConsume.Consume(p)

	p.HeroBox.AddHero(&playerdomain.Hero{
		ModelId: heroData.Id,
		Level:   1,
	})

	GetHeroService().ReCalculateHeroAttr(p, p.HeroBox.GetHero(heroData.Id), true)
	context.EventBus.Publish(events.PlayerEntityChange, p)

	return &protos.ResHeroCombine{
		Code: 0,
	}
}

func (ps *HeroController) ReqHeroUpFight(s *network.Session, index int, msg *protos.ReqHeroUpFight) *protos.ResHeroUpFight {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)

	h := p.HeroBox.GetHero(msg.HeroId)
	if h == nil {
		return &protos.ResHeroUpFight{
			Code: constants.I18N_COMMON_NOT_FOUND,
		}
	}

	pos := p.HeroBox.GetEmpostPos()
	if len(pos) == 0 {
		return &protos.ResHeroUpFight{
			Code: constants.I18N_HERO_TIP6,
		}
	}
	// 判断位置是否为空闲位置
	if !slices.Contains(pos, msg.Position) {
		return &protos.ResHeroUpFight{
			Code: constants.I18N_COMMON_ILLEGAL_PARAMS,
		}
	}

	h.Position = msg.Position
	GetHeroService().ReCalculateHeroAttr(p, h, true)
	context.EventBus.Publish(events.PlayerEntityChange, p)

	return &protos.ResHeroUpFight{
		Code: 0,
	}
}

func (ps *HeroController) ReqHeroOffFight(s *network.Session, index int, msg *protos.ReqHeroOffFight) *protos.ResHeroOffFight {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)

	h := p.HeroBox.GetHero(msg.HeroId)
	if h == nil {
		return &protos.ResHeroOffFight{
			Code: constants.I18N_COMMON_NOT_FOUND,
		}
	}

	if h.Position == 0 {
		return &protos.ResHeroOffFight{
			Code: constants.I18N_HERO_TIP7,
		}
	}

	if len(p.HeroBox.GetUpFightHeros()) == 1 {
		return &protos.ResHeroOffFight{
			Code: constants.I18N_HERO_TIP8,
		}
	}

	h.Position = 0
	GetHeroService().ReCalculateHeroAttr(p, h, true)
	context.EventBus.Publish(events.PlayerEntityChange, p)

	return &protos.ResHeroOffFight{
		Code: 0,
	}
}

func (ps *HeroController) ReqHeroChangePosition(s *network.Session, index int, msg *protos.ReqHeroChangePosition) *protos.ResHeroChangePosition {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)

	h := p.HeroBox.GetHero(msg.HeroId)
	if h == nil {
		return &protos.ResHeroChangePosition{
			Code: constants.I18N_COMMON_NOT_FOUND,
		}
	}

	if h.Position == msg.Position || h.Position == 0 {
		return &protos.ResHeroChangePosition{
			Code: constants.I18N_COMMON_ILLEGAL_PARAMS,
		}
	}

	// 如果目标位置有英雄，则表示交换
	prevHero := p.HeroBox.GetHeroByPosition(msg.Position)
	if prevHero != nil {
		prevPos := prevHero.Position
		prevHero.Position = h.Position
		h.Position = prevPos
		context.EventBus.Publish(events.PlayerEntityChange, p)
		return &protos.ResHeroChangePosition{
			Code:  0,
			PosA:  h.Position,
			HeroA: h.ModelId,
			PosB:  prevHero.Position,
			HeroB: prevHero.ModelId,
		}
	} else {
		h.Position = msg.Position
		context.EventBus.Publish(events.PlayerEntityChange, p)
		return &protos.ResHeroChangePosition{
			Code:  0,
			PosA:  h.Position,
			HeroA: h.ModelId,
			PosB:  0,
			HeroB: 0,
		}
	}
}

package hero

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"

	"io/github/gforgame/common"
	"io/github/gforgame/data"
	"io/github/gforgame/examples/camp"
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/config/container"
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/consume"
	"io/github/gforgame/examples/context"
	configdomain "io/github/gforgame/examples/domain/config"
	"io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/fight/attribute"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/examples/reward"
	"io/github/gforgame/examples/service/item"
	"io/github/gforgame/protos"
)

type HeroService struct{}

var (
	instance *HeroService
	once     sync.Once
)

func GetHeroService() *HeroService {
	once.Do(func() {
		instance = &HeroService{}
	})
	return instance
}

func (ps *HeroService) OnPlayerLogin(player *player.Player) {
	resAllHeroInfo := &protos.PushAllHeroInfo{}

	// 普通英雄
	for _, h := range player.HeroBox.Heros {
		ps.ReCalculateHeroAttr(player, h, false)
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

func (ps *HeroService) DoRecruit(p *player.Player, typ int32, times int32) (*common.BusinessRequestException, []*protos.RewardVo) {
	itemId := constants.ITEM_RECRUIT_ID
	if typ == 2 {
		itemId = constants.ITEM_RECRUIT_ID2
	}

	maxTimes := int32(51)
	// 检测次数
	if typ == 1 {
		if p.DailyReset.NormalRecruitTimes + times > maxTimes {
			return common.NewBusinessRequestException( constants.I18N_COMMON_ILLEGAL_PARAMS), nil
			}
	} else {
		if p.DailyReset.HighRecruitTimes + times > maxTimes {
			return common.NewBusinessRequestException( constants.I18N_COMMON_ILLEGAL_PARAMS), nil
			}
		}

	free := false
	// heroId := 0

	// 每天首次免费
	if times == 1 {
		if typ == 1{
			if !p.DailyReset.NormalRecruitFreeUsed {
				p.DailyReset.NormalRecruitFreeUsed = true
				free = true
			}
		}else {
			if !p.DailyReset.HighRecruitFreeUsed {
				p.DailyReset.HighRecruitFreeUsed = true
				free = true
			}
		}
	}

	if !free {
		// 优先消耗招募令
		if p.Backpack.IsEnough(itemId, times) {
				return common.NewBusinessRequestException( constants.I18N_ITEM_NOT_ENOUGH), nil
		} else {
			// 不足扣钻石
			itemCount := p.Backpack.GetItemCount(constants.ITEM_DIAMOND_ID)
			if itemCount < times {
				return common.NewBusinessRequestException( constants.I18N_ITEM_NOT_ENOUGH), nil
			}
			commonContainer := config.QueryContainer[configdomain.CommonData, *container.CommonContainer]()
			// 招募消耗钻石
			exchangeMoney := commonContainer.GetInt32Value("heroRecruitDiamond")
			owed := times - itemCount
			needMoney := owed * exchangeMoney
			if needMoney > 0 {
				itemConsume := consume.ItemConsume{
					ItemId: itemId,
					Amount: itemCount,
				}
				itemConsume.Consume(p, constants.ActionType_HeroRecruit)
			}
			if itemCount > 0 {
				item.GetItemService().UseByModelId(p, itemId, itemCount)
			}
		}
	
		if typ == 1{
			p.DailyReset.NormalRecruitTimes += times
		}else {
			p.DailyReset.HighRecruitTimes += times
		}
	}

	gachaContainer := config.QueryContainer[configdomain.GachaData, *container.GachaContainer]()

	rewardVos := make([]*protos.RewardVo, 0)
	for i := 0; i < int(times); i++ {
		gachaData := gachaContainer.RandItem(typ)
		rewards := reward.ParseReward(gachaData.Rewards)
		realReward := reward.GetSingleReward(rewards)
		if realReward.(*reward.HeroReward) != nil {
			heroReward := realReward.(*reward.HeroReward)
			heroData := config.QueryById[configdomain.HeroData](heroReward.HeroId)
			// 如果已经拥有该英雄，则转为碎片
			if p.HeroBox.HasHero(heroData.Id) {
				rewardVos = append(rewardVos, &protos.RewardVo{
					Type:  "item",
					Value: fmt.Sprintf("%d=%d", heroData.ShardItem, heroData.ShardAmount),
				})
				item.GetItemService().AddByModelId(p, heroData.ShardItem, heroData.ShardAmount)
			} else {	
				rewardVos = append(rewardVos, &protos.RewardVo{
					Type:  "hero",
					Value: strconv.Itoa(int(heroData.Id)),
				})
				ps.NewHero(p, heroData.Id)

				// ps.ReCalculateHeroAttr(p, p.HeroBox.GetHero(heroData.Id), true)
			}
		}
	}

	context.EventBus.Publish(events.PlayerEntityChange, p)

	return nil, rewardVos
}

func (ps *HeroService) NewHero(p *player.Player, heroId int32) {
	p.HeroBox.AddHero(&player.Hero{
				ModelId: heroId,
				Level:   1,
			})
	context.EventBus.Publish(events.HeroGain, &events.HeroGainEvent{
		Player: p,
		HeroId: heroId,
	})
}

func (ps *HeroService) GetRandomHero() configdomain.HeroData {
	heroDatas := ps.filterNormalHeros()
	// 根据HeroData的Prob进行抽奖
	var totalProb int32 = 0
	for _, heroData := range heroDatas {
		totalProb += heroData.Prob
	}

	randProb := rand.Int31n(totalProb)
	var currentProb int32 = 0
	var selectedHero configdomain.HeroData

	for _, heroData := range heroDatas {
		currentProb += heroData.Prob
		if randProb < currentProb {
			selectedHero = *heroData
			break
		}
	}

	return selectedHero
}

// 过滤掉主公
func (ps *HeroService) filterNormalHeros() []*configdomain.HeroData {
	container := config.QueryContainer[configdomain.HeroData, *data.Container[int32, configdomain.HeroData]]()

	var result []*configdomain.HeroData
	for _, heroData := range container.GetAllRecords() {
		// 主公概率为0
		if heroData.Prob > 0 {
			result = append(result, heroData)
		}
	}
	return result
}

func (ps *HeroService) DoLevelUp(p *player.Player, heroId int32, toLevel int32) *protos.ResHeroLevelUp {
	h := p.HeroBox.GetHero(heroId)
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

	costGold := ps.CalcTotalUpLevelConsume(h.Level, toLevel)
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
	consume.Consume(p, constants.ActionType_HeroUpLevel)

	h.Level = toLevel
	ps.ReCalculateHeroAttr(p, h, true)
	context.EventBus.Publish(events.PlayerEntityChange, p)

	return &protos.ResHeroLevelUp{
		Code: 0,
	}
}

func (ps *HeroService) DoStageUp(p *player.Player, heroId int32) *protos.ResHeroUpStage {
	h := p.HeroBox.GetHero(heroId)
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
	costItem.Consume(p, constants.ActionType_HeroUpStage)

	h.Stage = h.Stage + 1

	ps.ReCalculateHeroAttr(p, h, true)
	context.EventBus.Publish(events.PlayerEntityChange, p)

	return &protos.ResHeroUpStage{
		Code: 0,
	}
}

func (ps *HeroService) DoCombine(p *player.Player, heroId int32) *protos.ResHeroCombine {
	h := p.HeroBox.GetHero(heroId)
	if h != nil {
		return &protos.ResHeroCombine{
			Code: constants.I18N_HERO_TIP5,
		}
	}
	heroData := config.QueryById[configdomain.HeroData](heroId)
	itemConsume := consume.ItemConsume{
		ItemId: heroData.ShardItem,
		Amount: heroData.ShardAmount,
	}
	err := itemConsume.Verify(p)
	if err != nil {
		return &protos.ResHeroCombine{
			Code: int32(err.(*common.BusinessRequestException).Code()),
		}
	}
	itemConsume.Consume(p, constants.ActionType_HeroCombine)

	p.HeroBox.AddHero(&player.Hero{
		ModelId: heroData.Id,
		Level:   1,
	})

	ps.ReCalculateHeroAttr(p, p.HeroBox.GetHero(heroData.Id), true)
	context.EventBus.Publish(events.PlayerEntityChange, p)

	return &protos.ResHeroCombine{
		Code: 0,
	}
}

func (ps *HeroService) DoUpFight(p *player.Player, heroId int32) *protos.ResHeroUpFight {
	h := p.HeroBox.GetHero(heroId)
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

	h.Stage = h.Stage + 1

	ps.ReCalculateHeroAttr(p, h, true)
	context.EventBus.Publish(events.PlayerEntityChange, p)

	return &protos.ResHeroUpFight{
		Code: 0,
	}
}

func (ps *HeroService) DoOffFight(p *player.Player, heroId int32) *protos.ResHeroOffFight {
	h := p.HeroBox.GetHero(heroId)
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

	h.Position = 0

	ps.ReCalculateHeroAttr(p, h, true)
	context.EventBus.Publish(events.PlayerEntityChange, p)

	return &protos.ResHeroOffFight{
		Code: 0,
	}
}

func (ps *HeroService) DoChangePosition(p *player.Player, heroId int32, position int32) *protos.ResHeroChangePosition {
	h := p.HeroBox.GetHero(heroId)
	if h == nil {
		return &protos.ResHeroChangePosition{
			Code: constants.I18N_COMMON_NOT_FOUND,
		}
	}

	if h.Position == position {
		return &protos.ResHeroChangePosition{
			Code: constants.I18N_HERO_TIP8,
		}
	}

	h.Position = position

	ps.ReCalculateHeroAttr(p, h, true)
	context.EventBus.Publish(events.PlayerEntityChange, p)

	return &protos.ResHeroChangePosition{
		Code: 0,
	}
}

// 重新计算武将属性
func (ps *HeroService) ReCalculateHeroAttr(p *player.Player, hero *player.Hero, notify bool) {
	// 英雄本身属性
	heroData := config.QueryById[configdomain.HeroData](hero.ModelId)
	attrContainer := attribute.NewAttrBox()
	attrContainer.AddAttrs(heroData.GetHeroAttrs())

	// 英雄等级属性
	levelContainer := config.QueryContainer[configdomain.HeroLevelData, *container.HeroLevelContainer]()
	levelData := levelContainer.GetLevelData(hero.ModelId, hero.Level)
	if levelData != nil {
		attrContainer.AddAttrs(levelData.GetHeroLevelAttrs())
	}

	// 英雄突破属性
	stageContainer := config.QueryContainer[configdomain.HeroStageData, *container.HeroStageContainer]()
	stageData := stageContainer.GetRecordByStage(hero.Stage)
	if stageData != nil {
		attrContainer.AddAttrs(stageData.Attrs)
	}

	hero.AttrBox = attrContainer
	hero.Fight = attribute.CalculateFightingPower(attrContainer)

	if notify {
		attrs := make([]protos.AttrInfo, 0, len(attrContainer.Attrs))
		for attrType, value := range attrContainer.Attrs {
			attrs = append(attrs, protos.AttrInfo{
				AttrType: string(attrType),
				Value:    value,
			})
		}

		io.NotifyPlayer(p, &protos.PushHeroAttrChange{
			HeroId: int32(hero.ModelId),
			Attrs:  attrs,
			Fight:  attribute.CalculateFightingPower(attrContainer),
		})
	}

	context.EventBus.Publish(events.PlayerAttrChange, p)
}

func (ps *HeroService) CalcTotalUpLevelConsume(fromLevel int32, toLevel int32) int32 {
	levelContainer := config.QueryContainer[configdomain.HeroLevelData, *container.HeroLevelContainer]()
	total := int32(0)
	for i := fromLevel; i < toLevel; i++ {
		levelData := levelContainer.GetLevelData(i, i)
		if levelData != nil {
			total += levelData.Cost
		}
	}
	return total
}

package hero

import (
	"math/rand"
	"sync"

	"io/github/gforgame/examples/context"
	"io/github/gforgame/examples/domain/config"
	"io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/fight/attribute"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/protos"
)

type HeroService struct{}

var (
	instance        *HeroService
	once            sync.Once
	stageDataMapper map[int32]config.HeroStageData = make(map[int32]config.HeroStageData)
)

func GetHeroService() *HeroService {
	once.Do(func() {
		instance = &HeroService{}
		stageDatas := context.GetDataManager().GetRecords("herostage")
		for _, stageData := range stageDatas {
			stageData := stageData.(config.HeroStageData)
			stageDataMapper[stageData.Stage] = stageData
		}
	})
	return instance
}

func (ps *HeroService) GetRandomHero() config.HeroData {
	heroDatas := ps.filterNormalHeros()
	// 根据HeroData的Prob进行抽奖
	var totalProb int32 = 0
	for _, heroData := range heroDatas {
		totalProb += heroData.Prob
	}

	randProb := rand.Int31n(totalProb)
	var currentProb int32 = 0
	var selectedHero config.HeroData

	for _, heroData := range heroDatas {
		currentProb += heroData.Prob
		if randProb < currentProb {
			selectedHero = heroData
			break
		}
	}

	return selectedHero
}

// 过滤掉主公
func (ps *HeroService) filterNormalHeros() []config.HeroData {
	heroDatas := context.GetDataManager().GetRecords("hero")
	var result []config.HeroData
	for _, heroDataRecord := range heroDatas {
		heroData := heroDataRecord.(config.HeroData)
		// 主公概率为0
		if heroData.Prob > 0 {
			result = append(result, heroData)
		}
	}
	return result
}

// 重新计算武将属性
func (ps *HeroService) ReCalculateHeroAttr(p *player.Player, hero *player.Hero, notify bool) {
	// 英雄本身属性
	heroData := context.GetConfigRecordAs[config.HeroData]("hero", int64(hero.ModelId))
	attrContainer := attribute.NewAttrBox()
	attrContainer.AddAttrs(heroData.GetHeroAttrs())

	// 英雄等级属性
	heroLevelData := context.GetConfigRecordAs[config.HeroLevelData]("herolevel", int64(hero.Level))
	if heroLevelData != nil {
		attrContainer.AddAttrs(heroLevelData.GetHeroLevelAttrs())
	}

	// 英雄突破属性
	heroStageData := context.GetConfigRecordAs[config.HeroStageData]("herostage", int64(hero.Stage))
	if heroStageData != nil {
		attrContainer.AddAttrs(heroStageData.Attrs)
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
	totle := int32(0)
	for i := fromLevel; i < toLevel; i++ {
		record := context.GetDataManager().GetRecord("herolevel", int64(i))

		heroLevelData := record.(config.HeroLevelData)
		totle += heroLevelData.Cost
	}
	return totle
}

func (ps *HeroService) GetHeroStageData(stage int32) (config.HeroStageData, bool) {
	stageData, ok := stageDataMapper[stage]
	return stageData, ok
}

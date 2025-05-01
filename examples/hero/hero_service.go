package hero

import (
	"math/rand"
	"sync"

	"io/github/gforgame/data"
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/config/container"
	"io/github/gforgame/examples/context"
	configdomain "io/github/gforgame/examples/domain/config"
	"io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/fight/attribute"
	"io/github/gforgame/examples/io"
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

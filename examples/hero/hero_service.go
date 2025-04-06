package hero

import (
	"math/rand"
	"sync"

	"io/github/gforgame/examples/attribute"
	"io/github/gforgame/examples/context"
	"io/github/gforgame/examples/domain/config"
	"io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
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

func (ps *HeroService) GetRandomHero() config.HeroData {
	heroDatas := context.GetDataManager().GetRecords("hero")
	// 根据HeroData的Prob进行抽奖
	var totalProb int32 = 0
	for _, data := range heroDatas {
		heroData := data.(config.HeroData)
		totalProb += heroData.Prob
	}

	// Generate random number between 0 and total probability
	randProb := rand.Int31n(totalProb)
	var currentProb int32 = 0
	var selectedHero config.HeroData

	// Find the hero based on probability ranges
	for _, data := range heroDatas {
		heroData := data.(config.HeroData)
		currentProb += heroData.Prob
		if randProb < currentProb {
			selectedHero = heroData
			break
		}
	}

	return selectedHero
}

func (ps *HeroService) RecalculateHeroAttr(p *player.Player, hero *player.Hero, notify bool) {
	// 英雄本身属性
	heroDataRecord := context.GetDataManager().GetRecord("hero", int64(hero.ModelId))
	heroData := heroDataRecord.(config.HeroData)
	attrContainer := attribute.NewAttrBox()
	attrContainer.AddAttrs(heroData.GetHeroAttrs())

	// 英雄等级属性
	heroLevelDataRecord := context.GetDataManager().GetRecord("herolevel", int64(hero.Level))
	heroLevelData := heroLevelDataRecord.(config.HeroLevelData)
	attrContainer.AddAttrs(heroLevelData.GetHeroLevelAttrs())

	context.EventBus.Publish(events.PlayerAttrChange, p)

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
}

func (ps *HeroService) calcTotalUpLevelConsume(fromLevel int32, toLevel int32) int32 {
	totle := int32(0)
	for i := fromLevel; i < toLevel; i++ {
		record := context.GetDataManager().GetRecord("herolevel", int64(i))

		heroLevelData := record.(config.HeroLevelData)
		totle += heroLevelData.Cost
	}
	return totle
}

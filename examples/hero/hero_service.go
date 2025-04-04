package hero

import (
	"math/rand"
	"sync"

	"io/github/gforgame/examples/context"
	"io/github/gforgame/examples/domain/config"
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

func (ps *HeroService) calcTotalUpLevelConsume(fromLevel int32, toLevel int32) int32 {
	totle := int32(0)
	for i := fromLevel; i < toLevel; i++ {
		record := context.GetDataManager().GetRecord("herolevel", int64(i))

		heroLevelData := record.(config.HeroLevelData)
		totle += heroLevelData.Cost
	}
	return totle
}

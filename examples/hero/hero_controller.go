package hero

import (
	"fmt"
	"math/rand"
	"strconv"

	"io/github/gforgame/examples/context"
	"io/github/gforgame/examples/domain/config"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/item"
	"io/github/gforgame/examples/player"
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
}

func (ps *HeroController) ReqRecruit(s *network.Session, index int, msg *protos.ReqHeroRecruit) *protos.ResHeroRecruit {
	rewardInfos := make([]*protos.RewardInfo, 0)

	p := context.SessionManager.GetPlayerBySession(s).(*playerdomain.Player)

	for i := 0; i < int(msg.Times); i++ {
		heroData := ps.GetRandomHero()
		// 如果已经拥有该英雄，则转为碎片
		if p.HeroBox.HasHero(heroData.Id) {
			rewardInfos = append(rewardInfos, &protos.RewardInfo{
				Type:  "item",
				Value: fmt.Sprintf("%d=%d", heroData.Item, heroData.Shard),
			})
			item.GetInstance().AddByModelId(p, heroData.Item, heroData.Shard)
		} else {
			rewardInfos = append(rewardInfos, &protos.RewardInfo{
				Type:  "hero",
				Value: strconv.Itoa(int(heroData.Id)),
			})
			p.HeroBox.AddHero(&playerdomain.Hero{
				ModelId: heroData.Id,
			})
		}

	}

	player.GetPlayerService().SavePlayer(p)
	// Return the recruited hero info

	return &protos.ResHeroRecruit{
		Code:        0,
		RewardInfos: rewardInfos,
	}
}

func (ps *HeroController) GetRandomHero() config.HeroData {
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

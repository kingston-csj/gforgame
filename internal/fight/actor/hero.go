package actor

import (
	"github.com/forfun/gforgame/common/util"
	"github.com/forfun/gforgame/internal/config"
	configdomain "github.com/forfun/gforgame/internal/domain/config"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/fight/attribute"
	"github.com/forfun/gforgame/internal/fight/buff"
	"github.com/forfun/gforgame/internal/fight/state"
)

type Hero struct {
	baseActor
}

func NewHero(modelId int32, camp int32, attrBox *attribute.AttrBox, skills []int32) *Hero {
	return &Hero{
		baseActor: baseActor{
			id:       util.GetNextID(),
			modelId:  modelId,
			attrBox:  attrBox,
			hp:       int32(attrBox.GetAttr(attribute.Hp).Value),
			camp:     camp,
			skills:   skills,
			buffBox:  buff.NewBuffBox(),
			stateBox: state.NewStateBox(),
		},
	}
}

func NewHero2(hero *playerdomain.Hero) *Hero {
	heroData := config.QueryById[configdomain.HeroData](hero.ModelId)
	return &Hero{
		baseActor: baseActor{
			id:       util.GetNextID(),
			modelId:  hero.ModelId,
			attrBox:  hero.AttrBox,
			hp:       hero.AttrBox.GetAttr(attribute.Hp).Value,
			skills:   heroData.Skills,
			buffBox:  buff.NewBuffBox(),
			stateBox: state.NewStateBox(),
		},
	}
}

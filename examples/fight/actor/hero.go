package actor

import (
	"io/github/gforgame/examples/context"
	"io/github/gforgame/examples/domain/config"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/fight/attribute"
	"io/github/gforgame/examples/fight/buff"
	"io/github/gforgame/examples/fight/state"
	"io/github/gforgame/util"
)

type Hero struct {
	baseActor
}

func NewHero(modelId int32, camp int32, attrBox *attribute.AttrBox, skills []int32) *Hero {
	return &Hero{
		baseActor: baseActor{
			id:       util.GetNextId(),
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
	heroData := context.GetConfigRecordAs[config.HeroData]("hero", int64(hero.ModelId))

	return &Hero{
		baseActor: baseActor{
			id:       util.GetNextId(),
			modelId:  hero.ModelId,
			attrBox:  hero.AttrBox,
			hp:       hero.AttrBox.GetAttr(attribute.Hp).Value,
			skills:   heroData.Skills,
			buffBox:  buff.NewBuffBox(),
			stateBox: state.NewStateBox(),
		},
	}
}

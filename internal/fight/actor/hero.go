package actor

import (
	"github.com/forfun/gforgame/internal/config"
	"github.com/forfun/gforgame/internal/domain/attr"
	configdomain "github.com/forfun/gforgame/internal/domain/config"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/fight/buff"
	"github.com/forfun/gforgame/internal/fight/state"
	"github.com/forfun/gforgame/internal/idgen"
)

type Hero struct {
	baseActor
}

func NewHero(modelId int32, camp int32, attrBox *attr.AttrBox, skills []int32) *Hero {
	return &Hero{
		baseActor: baseActor{
			id:       idgen.GetNextID(),
			modelId:  modelId,
			attrBox:  attrBox,
			hp:       int32(attrBox.GetAttr(attr.Hp).Value),
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
			id:       idgen.GetNextID(),
			modelId:  hero.ModelId,
			attrBox:  hero.AttrBox,
			hp:       hero.AttrBox.GetAttr(attr.Hp).Value,
			skills:   heroData.Skills,
			buffBox:  buff.NewBuffBox(),
			stateBox: state.NewStateBox(),
		},
	}
}

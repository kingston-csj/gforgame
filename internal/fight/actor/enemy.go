package actor

import (
	"github.com/forfun/gforgame/internal/fight/attribute"

	"github.com/forfun/gforgame/common/util"
	"github.com/forfun/gforgame/internal/fight/state"
)

type Enemy struct {
	baseActor
}

func NewEnemy(modelId int32, camp int32, attrBox *attribute.AttrBox, skills []int32) *Enemy {
	return &Enemy{
		baseActor: baseActor{
			id:       util.GetNextID(),
			modelId:  modelId,
			attrBox:  attrBox,
			hp:       int32(attrBox.GetAttr(attribute.Hp).Value),
			camp:     camp,
			skills:   skills,
			stateBox: state.NewStateBox(),
		},
	}
}

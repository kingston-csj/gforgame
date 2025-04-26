package actor

import (
	"io/github/gforgame/examples/fight/attribute"

	"io/github/gforgame/examples/fight/state"
	"io/github/gforgame/util"
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

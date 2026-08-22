package actor

import (
	"github.com/forfun/gforgame/internal/domain/attr"
	"github.com/forfun/gforgame/internal/fight/state"
	"github.com/forfun/gforgame/internal/idgen"
)

type Enemy struct {
	baseActor
}

func NewEnemy(modelId int32, camp int32, attrBox *attr.AttrBox, skills []int32) *Enemy {
	return &Enemy{
		baseActor: baseActor{
			id:       idgen.GetNextID(),
			modelId:  modelId,
			attrBox:  attrBox,
			hp:       int32(attrBox.GetAttr(attr.Hp).Value),
			camp:     camp,
			skills:   skills,
			stateBox: state.NewStateBox(),
		},
	}
}

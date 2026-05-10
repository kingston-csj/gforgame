package actor

import (
	"github.com/forfun/gforgame/internal/fight/attribute"
	"github.com/forfun/gforgame/internal/fight/buff"
	"github.com/forfun/gforgame/internal/fight/state"
)

type Actor interface {
	GetId() string

	GetAttackSpeed() int32

	GetModelId() int32

	GetHp() int32

	IsDead() bool

	ChangeHp(delta int32)

	GetCamp() int32

	GetSkills() []int32

	// 获取下一个技能
	NextSkill() int32

	GetAttrValue(attrType attribute.AttrType) int32

	GetBuffBox() *buff.BuffBox

	GetStateBox() *state.StateBox
}

package actor

import (
	"io/github/gforgame/examples/context"
	"io/github/gforgame/examples/domain/config"
	"io/github/gforgame/examples/fight/attribute"
	"io/github/gforgame/examples/fight/buff"
	"io/github/gforgame/examples/fight/state"
	"math/rand"
	"time"
)

type baseActor struct {
	attrBox  *attribute.AttrBox
	id       string
	modelId  int32
	hp       int32
	camp     int32
	skills   []int32
	buffBox  *buff.BuffBox
	stateBox *state.StateBox
}

func (a *baseActor) GetAttrContainer() *attribute.AttrBox {
	return a.attrBox
}

func (a *baseActor) GetId() string {
	return a.id
}

func (a *baseActor) GetModelId() int32 {
	return a.modelId
}

func (a *baseActor) GetHp() int32 {
	return a.hp
}

func (a *baseActor) GetCamp() int32 {
	return a.camp
}

func (a *baseActor) IsDead() bool {
	return a.hp <= 0
}

func (a *baseActor) ChangeHp(delta int32) {
	a.hp += delta
}

func (a *baseActor) GetAttackSpeed() int32 {
	return int32(a.attrBox.GetAttr(attribute.Speed).Value)
}

func (a *baseActor) GetSkills() []int32 {
	return a.skills
}

func (a *baseActor) GetStateBox() *state.StateBox {
	return a.stateBox
}

// 获取下一个技能
// 如果技能列表为空，则返回0
// 如果技能列表只有一个技能，则返回该技能
// 如果技能列表有多个技能，则随机返回一个技能
func (a *baseActor) NextSkill() int32 {
	skillIds := a.skills
	if len(skillIds) == 0 {
		return 0
	}
	tmpSkillIds := make([]int32, 0)

	for _, skillId := range skillIds {
		skillDataRecord := context.GetDataManager().GetRecord("skill", int64(skillId))
		skillData := skillDataRecord.(config.SkillData)
		if skillData.Type == 1 {
			tmpSkillIds = append(tmpSkillIds, skillId)
		}
	}
	if len(tmpSkillIds) == 0 {
		return 0
	}
	if len(tmpSkillIds) == 1 {
		return tmpSkillIds[0]
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	skillId := tmpSkillIds[r.Intn(len(tmpSkillIds))]
	return skillId
}

// 属性值由两部分组成： 本身养成属性+战斗触发buff
func (a *baseActor) GetAttrValue(attrType attribute.AttrType) int32 {
	return a.attrBox.GetAttrValue(attrType) + a.buffBox.GetAttrValue(attrType)
}

func (a *baseActor) GetBuffBox() *buff.BuffBox {
	return a.buffBox
}

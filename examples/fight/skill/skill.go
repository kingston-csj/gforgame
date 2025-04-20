package skill

import (
	"io/github/gforgame/examples/context"
	"io/github/gforgame/examples/domain/config"
	"io/github/gforgame/examples/fight/actor"
)

type Skill interface {
	GetSkillId() int32

	GetEffectType() int32

	// 解析技能参数
	AnalyzeParams(params string)
	// 施法者
	GetAttacker() actor.Actor

	// 计算伤害或治疗值
	CalculateDamage(attacker actor.Actor, target actor.Actor) int32

	// 技能原型
	Prototype() config.SkillData
}

type baseSkill struct {
	Skill
	SkillId int32
	value   string
}

func (s *baseSkill) AnalyzeParams(params string) {
	s.value = params
}

func (s *baseSkill) GetSkillId() int32 {
	return s.SkillId
}

func (s *baseSkill) Prototype() config.SkillData {
	return *context.GetConfigRecordAs[config.SkillData]("skill", int64(s.GetSkillId()))
}

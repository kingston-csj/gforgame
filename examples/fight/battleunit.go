package fight

import (
	"io/github/gforgame/examples/fight/actor"
	"io/github/gforgame/examples/fight/attribute"
	"io/github/gforgame/examples/fight/skill"
)

type BattleUnit struct {
	skill    skill.Skill
	attacker actor.Actor
	defender actor.Actor
	hurt     int32
}

func NewBattleUnit(skill skill.Skill, attacker actor.Actor, defender actor.Actor) *BattleUnit {
	return &BattleUnit{skill: skill, attacker: attacker, defender: defender}
}

func (b *BattleUnit) GetAttacker() actor.Actor {
	return b.attacker
}

func (b *BattleUnit) GetDefender() actor.Actor {
	return b.defender
}

func (b *BattleUnit) GetHurt() int32 {
	return b.hurt
}

func (b *BattleUnit) GetSkill() skill.Skill {
	return b.skill
}

/*
* 计算伤害
* 伤害 = 攻方攻击力*伤害倍率/10000 - 守方防御力
* 伤害最小为1
 */
func (b *BattleUnit) CalculateHurt() int32 {
	skillData := b.skill.Prototype()

	damageRate := skillData.DamageRate

	attacker := b.GetAttacker()
	defender := b.GetDefender()
	// 攻方攻击力*伤害倍率-对方防御力
	damage := int32(float32(attacker.GetAttrValue(attribute.Attack)) * float32(damageRate) / 10000)
	damage = max(damage-int32(defender.GetAttrValue(attribute.Defense)), 1)

	return damage
}

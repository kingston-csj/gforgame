package skill

import (
	"io/github/gforgame/examples/context"
	"io/github/gforgame/examples/domain/config"
)

func NewSkill(skillId int32) Skill {
	skillData := context.GetConfigRecordAs[config.SkillData]("skill", int64(skillId))
	switch skillData.EffectType {
	case 1:
		return &CommonSkill{
			baseSkill: baseSkill{
				SkillId: skillId,
			},
		}
	default:
		return nil
	}
}

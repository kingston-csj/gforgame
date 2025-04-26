package skill

import (
	"io/github/gforgame/examples/config"
	configdomain "io/github/gforgame/examples/domain/config"
)

func NewSkill(skillId int32) Skill {
	skillData := config.QueryById[configdomain.SkillData](skillId)
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

package skill

import (
	"github.com/forfun/gforgame/examples/config"
	configdomain "github.com/forfun/gforgame/examples/domain/config"
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

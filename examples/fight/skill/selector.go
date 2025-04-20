package skill

import (
	"io/github/gforgame/examples/domain/config"
	"io/github/gforgame/examples/fight/actor"
	"io/github/gforgame/examples/fight/match"
	"math/rand/v2"
)

var (
	selector1 = &SelfRangeSelector{}
	selector2 = &FriendRangeSelector{}
	selector3 = &EnemyRangeSelector{}
)

// 技能选择器 根据技能类型选择目标
type Selector interface {
	// 选择目标
	Select(match *match.Match, skill *config.SkillData, actor actor.Actor) []actor.Actor
	// 技能类型
	Type() int32
}

func GetSelector(typ int32) Selector {
	switch typ {
	case 1:
		return selector1
	case 2:
		return selector2
	case 3:
		return selector3
	default:
		return nil
	}
}

// 仅对施法者自己生效
type SelfRangeSelector struct {
	Selector
}

func (s *SelfRangeSelector) Select(match *match.Match, skill *config.SkillData, a actor.Actor) []actor.Actor {
	return []actor.Actor{a}
}

func (s *SelfRangeSelector) Type() int32 {
	return 1
}

// 对施法者自己以及友方生效
type FriendRangeSelector struct {
	Selector
}

func (s *FriendRangeSelector) Select(match *match.Match, skill *config.SkillData, a actor.Actor) []actor.Actor {
	target := match.GetMyTeam(a).GetLivingActors()
	if int32(len(target)) > skill.AoeRange {
		// 随机取
		rand.Shuffle(len(target), func(i, j int) {
			target[i], target[j] = target[j], target[i]
		})
		target = target[:skill.AoeRange]
	}
	return target
}

func (s *FriendRangeSelector) Type() int32 {
	return 2
}

// 对施法者敌方单位生效
type EnemyRangeSelector struct {
	Selector
}

func (s *EnemyRangeSelector) Select(match *match.Match, skill *config.SkillData, a actor.Actor) []actor.Actor {
	target := match.GetEnemyTeam(a).GetLivingActors()
	if int32(len(target)) > skill.AoeRange {
		// 随机取
		rand.Shuffle(len(target), func(i, j int) {
			target[i], target[j] = target[j], target[i]
		})
		target = target[:skill.AoeRange]
	}
	return target
}

func (s *EnemyRangeSelector) Type() int32 {
	return 3
}

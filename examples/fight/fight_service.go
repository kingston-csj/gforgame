package fight

import (
	"io/github/gforgame/examples/config"
	configdomain "io/github/gforgame/examples/domain/config"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/fight/actor"
	"io/github/gforgame/examples/fight/match"
	"io/github/gforgame/examples/fight/report"
	skillservice "io/github/gforgame/examples/fight/skill"
	"io/github/gforgame/examples/player"
	"sort"
	"sync"
)

var (
	instance *FightService
	once     sync.Once
)

func GetFightService() *FightService {
	once.Do(func() {
		instance = &FightService{}
	})
	return instance
}

type FightService struct {
}

func (s *FightService) StartFight(match *match.Match) {
	// 战斗开始，给所有角色添加永久buff
	allActors := match.GetAllLiveActors()
	for _, actor := range allActors {
		skillIds := actor.GetSkills()
		for _, skillId := range skillIds {
			skillData := config.QueryById[configdomain.SkillData](skillId)
			if skillData.BuffId > 0 {
				actor.GetBuffBox().AddBuff(skillData.BuffId)
			}
		}
		// 重算战斗属性
		actor.GetBuffBox().RefreshAttrs()
	}

	round := int32(1)
	for {
		if match.CheckWin() != 0 {
			break
		}
		s.RoundBegin(match, round)
		round++
	}
	match.Report.Winner = match.CheckWin()
	match.Report.Display()
}

func (s *FightService) RoundBegin(match *match.Match, round int32) {
	roundReport := report.NewRoundReport(round)
	allActors := match.GetAllLiveActors()
	// 根据 attackSpeed 排序
	sort.Slice(allActors, func(i, j int) bool {
		return allActors[i].GetAttackSpeed() < allActors[j].GetAttackSpeed()
	})

	// 从后往前遍历
	for i := len(allActors) - 1; i >= 0; i-- {
		actor := allActors[i]
		if actor.IsDead() {
			continue
		}
		if match.CheckWin() != 0 {
			break
		}
		if !actor.GetStateBox().CanAttack() {
			continue
		}

		skillId := actor.NextSkill()

		skillReport := report.NewSkillReport(actor.GetId(), skillId)

		skillData := config.QueryById[configdomain.SkillData](skillId)
		// 选择目标
		selector := skillservice.GetSelector(skillData.Selector)
		targets := selector.Select(match, skillData, actor)

		skill := skillservice.NewSkill(skillData.Id)

		// 多人技能
		for _, target := range targets {
			unit := NewBattleUnit(skill, actor, target)
			hurt := unit.CalculateHurt()
			target.ChangeHp(-hurt)
			// 添加伤害单位
			damageUnit := report.NewDamageUnit(actor.GetId(), target.GetId(), hurt)
			// 添加到单次技能报告
			skillReport.AddDamageUnit(damageUnit)
		}

		// 添加到回合报告
		roundReport.AddSkillReport(skillReport)
	}

	// 检查buff生命周期
	for _, actor := range allActors {
		actor.GetBuffBox().CheckBuffLife()
	}
	match.Report.AddRoundReport(roundReport)
}

func (s *FightService) Test() {
	p1 := player.GetPlayerService().GetPlayer("111")
	p2 := player.GetPlayerService().GetPlayer("aaa")
	team1 := match.NewTeam(match.BlueCamp, s.getFightActors(p1))
	team2 := match.NewTeam(match.RedCamp, s.getFightActors(p2))
	match := match.NewMatch(team1, team2)
	s.StartFight(match)
}

func (s *FightService) getFightActors(p *playerdomain.Player) []actor.Actor {
	fighters := make([]actor.Actor, 0)
	for _, hero := range p.HeroBox.GetAllHeros() {
		fighters = append(fighters, actor.NewHero2(hero))
	}
	return fighters
}

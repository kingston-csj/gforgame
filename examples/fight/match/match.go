package match

import (
	"io/github/gforgame/examples/fight/actor"
	"io/github/gforgame/examples/fight/report"
	"io/github/gforgame/util"
)

type Match struct {
	Id string
	// 红方队伍
	RedTeam *Team
	// 蓝方队伍
	BlueTeam *Team
	// 当前回合
	Round int32
	// 战斗报告
	Report *report.BattleReport
}

func NewMatch(redTeam *Team, blueTeam *Team) *Match {
	return &Match{Id: util.GetNextId(), RedTeam: redTeam, BlueTeam: blueTeam, Report: report.NewBattleReport()}
}

func (m *Match) GetMyTeam(actor actor.Actor) *Team {
	if actor.GetCamp() == int32(m.RedTeam.Camp) {
		return m.RedTeam
	}
	return m.BlueTeam
}

func (m *Match) GetEnemyTeam(actor actor.Actor) *Team {
	if actor.GetCamp() == int32(m.RedTeam.Camp) {
		return m.BlueTeam
	}
	return m.RedTeam
}

// 获取所有存活的角色
func (m *Match) GetAllLiveActors() []actor.Actor {
	allActors := append(m.RedTeam.Actors, m.BlueTeam.Actors...)
	liveActors := make([]actor.Actor, 0)
	for _, actor := range allActors {
		if !actor.IsDead() {
			liveActors = append(liveActors, actor)
		}
	}
	return liveActors
}

func (m *Match) CheckWin() int32 {
	if m.RedTeam.IsDead() {
		return BlueCamp
	}
	if m.BlueTeam.IsDead() {
		return RedCamp
	}
	return 0
}

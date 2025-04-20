package match

import "io/github/gforgame/examples/fight/actor"

type Team struct {
	Camp   int32
	Actors []actor.Actor // 战斗单位
}

func NewTeam(camp int32, actors []actor.Actor) *Team {
	return &Team{Camp: camp, Actors: actors}
}

func (t *Team) GetLivingActors() []actor.Actor {
	livingActors := make([]actor.Actor, 0)
	for _, actor := range t.Actors {
		if !actor.IsDead() {
			livingActors = append(livingActors, actor)
		}
	}
	return livingActors
}

func (t *Team) IsDead() bool {
	return len(t.GetLivingActors()) == 0
}

package player

import (
	"io/github/gforgame/examples/config"
	configdomain "io/github/gforgame/examples/domain/config"
	"io/github/gforgame/protos"
)


type Quest struct {
	Id int32

	Type int32

	Progress int32

	Target int32

	Status int8
}

func (q *Quest) IsComplete() bool {
	return q.Progress >= q.Target
}

func (q *Quest) AddProgress(progress int32) {
	q.Progress += progress
}

func (q *Quest) SetProgress(progress int32) {
	q.Progress = progress
}

func (q *Quest) ToVo() *protos.QuestVo {
	return &protos.QuestVo{
		Id:       int32(q.Id),
		Progress: int32(q.Progress),
		Target:   int32(q.Target),
		Status:   q.Status,
	}
}

func (q *Quest) Prototype() *configdomain.QuestData {
	return config.QueryById[configdomain.QuestData](q.Id)
}

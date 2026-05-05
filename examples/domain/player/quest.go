package player

import (
	"github.com/forfun/gforgame/common/util/timeutil"
	"github.com/forfun/gforgame/examples/config"
	"github.com/forfun/gforgame/examples/constants"
	configdomain "github.com/forfun/gforgame/examples/domain/config"
	"github.com/forfun/gforgame/examples/protos"
)

type Quest struct {
	Id int32

	Type int32
	// 进度
	Progress int32
	// 目标
	Target int32
	// 状态
	Status int8
	// 接受时间（毫秒）
	AcceptTime int64
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

func (q *Quest) Reset() {
	q.Progress = 0
	q.Status = int8(constants.QuestStatusInit)
	q.AcceptTime = 0
}

func (q *Quest) ToVo() *protos.QuestVo {
	return &protos.QuestVo{
		Id:         int32(q.Id),
		Progress:   int32(q.Progress),
		Target:     int32(q.Target),
		Status:     q.Status,
		AcceptTime: q.AcceptTime * timeutil.MILLIS_PER_SECOND,
	}
}

func (q *Quest) Prototype() *configdomain.QuestData {
	return config.QueryById[configdomain.QuestData](q.Id)
}

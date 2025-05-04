package player

import (
	"time"

	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/domain/config/item"
)

type Mail struct {
	Id      int64            `json:"id"`
	Title   string           `json:"title"`
	Content string           `json:"content"`
	Rewards []item.RewardDef `json:"rewards"`
	Status  int32            `json:"status"`
	Time    int64            `json:"time"`
}

func (m *Mail) IsExpired() bool {
	return time.Now().Unix() > m.Time+constants.MAIL_EXPIRE_TIME
}

package player

import (
	"time"

	"io/github/gforgame/domain"
)

type Mail struct {
	Id      string            `json:"id"`
	Title   string           `json:"title"`
	Content string           `json:"content"`
	Rewards []domain.RewardDefLite `json:"rewards"`
	Status  int32            `json:"status"`
	// 接收时间（单位：秒）
	Time    int64            `json:"time"`
	// 过期时间（单位：秒）
	ExpiredTime int64            `json:"expiredTime"`
	Params []string `json:"params"`
}

func (m *Mail) IsExpired() bool {
	return time.Now().Unix() > m.ExpiredTime
}

package config

import (
	domainconsume "github.com/forfun/gforgame/internal/domain/consume"
	domainreward "github.com/forfun/gforgame/internal/domain/reward"
)

type MallData struct {
	Id       int32                 `json:"id" excel:"id"`
	Type     int32                 `json:"type" excel:"type"`
	Rewards  domainreward.Reward   `json:"rewards" excel:"rewards"`
	Consumes domainconsume.Consume `json:"consume" excel:"consumes"`
	// 每日限购
	DailyBuy int32 `json:"dailyBuy" excel:"dailyBuy"`
	// 终身限购
	LifeTimeBuy int32 `json:"lifeTimeBuy" excel:"lifeTimeBuy"`
}

func (m *MallData) GetId() int32 {
	return m.Id
}

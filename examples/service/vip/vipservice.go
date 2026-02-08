package vip

import (
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/config/container"
	"io/github/gforgame/examples/context"
	configdomain "io/github/gforgame/examples/domain/config"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/protos"
	"sync"
	"time"
) 


type VipService struct{}


var (
	instance *VipService
	once     sync.Once
)

func GetVipService() *VipService {
	once.Do(func() {
		instance = &VipService{}	
	})
	return instance
}

func (v *VipService) CheckRecharge(p *playerdomain.Player, rechargeId int32 ) {
	// rechargeData := config.QueryById[configdomain.RechargeData](int64(rechargeId))
	vipContainer := config.GetSpecificContainer[*container.VipContainer]()
	var result []*configdomain.VipData
	 // 从最大开始找
	for _, vipData := range vipContainer.GetAllRecords() {
		if p.RechargeRmb >= vipData.Money {
			result = append(result, vipData)
		}
	}
	// 取最后一个
	if len(result) > 0 {
		topVip := result[len(result)-1].Id
		if p.VipLevel != topVip {
			p.VipLevel = topVip
		}
	}

	commonContainer := config.GetSpecificContainer[*container.CommonContainer]()
	// vip每个周期的充值金额
	periodMoney := commonContainer.GetFloat32Value("vipPeriodMoney")
	newMoney :=p.ExtendBox.AddVipPeriodMoney(periodMoney)
	if newMoney >= periodMoney {
		p.ExtendBox.VipExpiredTime = time.Now().Unix() + int64(commonContainer.GetInt32Value("vipPeriod"))
		p.ExtendBox.VipPeriodMoney = newMoney - periodMoney
	}

	context.EventBus.Publish(events.PlayerEntityChange, p)
}

func (v *VipService) RefreshVipInfo(p *playerdomain.Player) {
	push := protos.PushVipQueryInfo{
		VipLevel: p.VipLevel,
		ExpiredTime: p.ExtendBox.VipExpiredTime,
		RechargeRmb: float32(p.RechargeRmb),
		PeriodRechargeRmb: p.ExtendBox.VipPeriodMoney,
	}
	io.NotifyPlayer(p, push)
}

// getExtraArenaTimes 获取额外的竞技场次数
func (v *VipService) GetExtraArenaTimes(p *playerdomain.Player) int32 {
	vipData := v.getEffectiveVipData(p)
	if vipData == nil {
		return 0
	}
	return vipData.ArenaTimes
}

func (v *VipService) getEffectiveVipData(p *playerdomain.Player) *configdomain.VipData {
	if p.VipLevel <= 0 {
		return nil
	}
	// 检查是否过期
	if p.ExtendBox.VipExpiredTime < time.Now().Unix() {
		return nil
	}
	vipContainer := config.GetSpecificContainer[*container.VipContainer]()
	vipData := vipContainer.GetRecord(p.VipLevel)
	return vipData
}

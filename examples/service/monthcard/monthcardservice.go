package monthcard

import (
	"fmt"
	"io/github/gforgame/common"
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/constants"
	configdomain "io/github/gforgame/examples/domain/config"
	"io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/examples/reward"
	"io/github/gforgame/protos"
	"sync"
	"time"
)

// 月卡服务
type MonthCardService struct{}

var (
	instance *MonthCardService
	once     sync.Once
)

func GetMonthCardService() *MonthCardService {
	once.Do(func() {
		instance = &MonthCardService{}
	})
	return instance
}

func (ps *MonthCardService) OnPlayerLogin(player *player.Player) {
	ps.sendInfo(player)
}

func (ps *MonthCardService) OnDailyReset(player *player.Player) {

}

func (ps *MonthCardService) sendInfo(player *player.Player) {
    push := &protos.PushMonthCardInfo{}
    silverCard := player.RechargeBox.GetOrCreateMonthlyCardVo(constants.MonthCardTypeSilver)
	goldCard := player.RechargeBox.GetOrCreateMonthlyCardVo(constants.MonthCardTypeGold)

	silverCardVo := &protos.MonthlyCardVo{
		ExpiredTime: silverCard.ExpiredTime,
	}
	goldCardVo := &protos.MonthlyCardVo{
		ExpiredTime: goldCard.ExpiredTime,
	}
	push.SilverCard = silverCardVo
	push.GoldCard = goldCardVo
	
    io.NotifyPlayer(player, push)
}

func (ps *MonthCardService) OnRecharge(player *player.Player, rechargeId int32) {
	rechargeData := config.QueryById[configdomain.RechargeData](rechargeId)
	if rechargeData == nil {
		return
	}
	if rechargeData.Type != constants.RechargeTypeMonthCard {
		return
	}
	cardType := constants.MonthCardTypeSilver
	monthCard := player.RechargeBox.GetOrCreateMonthlyCardVo(constants.MonthCardTypeSilver)
	if rechargeData.Id != 20 {
		monthCard = player.RechargeBox.GetOrCreateMonthlyCardVo(constants.MonthCardTypeGold)
		cardType = constants.MonthCardTypeGold
	}
	monthCardData := config.QueryById[configdomain.MonthlyCardData](cardType)
	now := time.Now().Unix()
	//  过期，从今天开始算
	if monthCard.ExpiredTime < now{
		monthCard.ExpiredTime, _ = calcExpiredTime(int(monthCardData.ValidDays))
		ps.sendInfo(player)
	}
}


// calcExpiredTime
// 功能：当前时间截断到当天0点 → 加lastDays天 → 设为23:59:59 → 转东八区毫秒时间戳
func calcExpiredTime(lastDays int) (int64, error) {
	// 1. 加载东八区时区（UTC+8，对应Java的ZoneOffset.ofHours(8)）
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return 0, fmt.Errorf("加载东八区时区失败：%w", err)
	}

	// 2. 获取当前东八区时间，并截断到当天00:00:00
	now := time.Now().In(loc)
	truncatedNow := time.Date(
		now.Year(),   // 年
		now.Month(),  // 月
		now.Day(),    // 日
		0,            // 小时置0
		0,            // 分钟置0
		0,            // 秒置0
		0,            // 纳秒置0
		loc,          // 时区
	)

	// 3. 加上lastDays天，得到目标日期
	targetDate := truncatedNow.AddDate(0, 0, lastDays)

	// 4. 构造目标日期的23:59:59
	// 用time.Date重新构造时间，指定时分秒为23:59:59，纳秒为0
	targetDateTime := time.Date(
		targetDate.Year(),   // 目标年
		targetDate.Month(),  // 目标月
		targetDate.Day(),    // 目标日
		23,                  // 小时设为23
		59,                  // 分钟设为59
		59,                  // 秒设为59
		0,                   // 纳秒置0
		loc,                 // 东八区时区
	)

	// 5. 转换为毫秒级时间戳
	epochMilli := targetDateTime.Unix()*1000 + int64(targetDateTime.Nanosecond()/1e6)

	return epochMilli, nil
}


func (ps *MonthCardService) TakeReward( player *player.Player, typ int32) *common.BusinessRequestException {
    monthCard := player.RechargeBox.GetOrCreateMonthlyCardVo(constants.MonthCardTypeSilver) 
	if typ == 1 {
		monthCard = player.RechargeBox.GetOrCreateMonthlyCardVo(constants.MonthCardTypeGold)
	}
	if monthCard.IsActivated() {
		return common.NewBusinessRequestException(constants.I18N_MONTH_CARD_TIPS1)
	}
	monthCardData := config.QueryById[configdomain.MonthlyCardData](typ)

	if typ == 1 {
		if player.DailyReset.SilverMonthCardReward {
			return common.NewBusinessRequestException(constants.I18N_MONTH_CARD_TIPS2)
		}
		player.DailyReset.SilverMonthCardReward = true
	} else {
		if player.DailyReset.GoldMonthCardReward {
			return common.NewBusinessRequestException(constants.I18N_MONTH_CARD_TIPS2)
		}
		player.DailyReset.GoldMonthCardReward = true
	}
	rewards := reward.ParseReward(monthCardData.Rewards)
	rewards.Reward(player, constants.ActionType_MonthCardGetReward)
	return nil
}

func (ps *MonthCardService) GetEffectiveMonthCardDatas(player *player.Player) []*configdomain.MonthlyCardData {
	monthCardDatas := make([]*configdomain.MonthlyCardData, 0)
	silverCard := player.RechargeBox.GetOrCreateMonthlyCardVo(constants.MonthCardTypeSilver)
	if silverCard.IsActivated() {
		monthCardDatas = append(monthCardDatas, config.QueryById[configdomain.MonthlyCardData](constants.MonthCardTypeSilver))
	}
	goldCard := player.RechargeBox.GetOrCreateMonthlyCardVo(constants.MonthCardTypeGold)
	if goldCard.IsActivated() {
		monthCardDatas = append(monthCardDatas, config.QueryById[configdomain.MonthlyCardData](constants.MonthCardTypeGold))
	}
	return monthCardDatas
}

// GetExtraArenaTimes 获取额外的竞技场次数
func (ps *MonthCardService) GetExtraArenaTimes(player *player.Player) int32 {
	monthCardDatas := ps.GetEffectiveMonthCardDatas(player)
	if len(monthCardDatas) == 0 {
		return 0
	}
	extraTimes := int32(0)
	for _, monthCardData := range monthCardDatas {
		extraTimes += monthCardData.ArenaTimes
	}
	return extraTimes
}
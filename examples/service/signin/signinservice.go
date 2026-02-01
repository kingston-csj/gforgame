package signin

import (
	"io/github/gforgame/common"
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/context"
	configdomain "io/github/gforgame/examples/domain/config"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/examples/reward"
	"io/github/gforgame/protos"
	"sync"
	"time"
)

type SignInService struct {
}

var (
	instance *SignInService
	once     sync.Once
)

func GetSignInService() *SignInService {
	once.Do(func() {
		instance = &SignInService{}
	})
	return instance
}

func (s *SignInService) OnPlayerLogin(player *playerdomain.Player) {
    push := &protos.PushSigninInfo{};
    push.SigninDays = player.MonthlyReset.SignInDays
    // 本月总天数
    push.DaysInMonth = getDaysOfCurrMonth();
    push.NthDay = int32(time.Now().Day());
	push.SignInMakeUp = player.MonthlyReset.SignInMakeUp
   
    io.NotifyPlayer(player, push)
}

func getDaysOfCurrMonth() int32 {
	now := time.Now()
	nextMonthFirstDay := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
	lastDayOfCurrentMonth := nextMonthFirstDay.Add(-24 * time.Hour)
	daysInMonth := lastDayOfCurrentMonth.Day()
	return int32(daysInMonth)
}

func IntSliceContains(slice []int32, target int32) bool {
	for _, num := range slice {
		if num == target {
			return true
		}
	}
	return false
}

func (s *SignInService) SignIn(player *playerdomain.Player) *common.BusinessRequestException {
    nthDay := int32(time.Now().Day())
    monthlyResetBox := player.MonthlyReset
    if IntSliceContains(monthlyResetBox.SignInDays, nthDay) {
        return common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
    }
    monthlyResetBox.SignInDays = append(monthlyResetBox.SignInDays, nthDay);

    signinData := config.QueryById[configdomain.SigninData](nthDay)
    if (signinData == nil) {
        return common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
    }
    rewards := reward.ParseReward(signinData.Rewards)
    rewards.Reward(player, constants.ActionType_Signin)
    context.EventBus.Publish(events.PlayerEntityChange, player)
    return nil
}

func (s *SignInService) SignInMakeUp(player *playerdomain.Player, day int32) *common.BusinessRequestException {
    monthlyResetBox := player.MonthlyReset
    if _, ok := monthlyResetBox.SignInMakeUp[day]; ok {
        return common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
    }
	nthDay := int32(time.Now().Day())
	if day < 1 || day > nthDay {
        return common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
    }
    monthlyResetBox.SignInMakeUp[day] = day
	context.EventBus.Publish(events.PlayerEntityChange, player)
    return nil
}

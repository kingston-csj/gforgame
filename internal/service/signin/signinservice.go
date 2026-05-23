package signin

import (
	"time"

	"github.com/forfun/gforgame/common/errors"
	commonerrors "github.com/forfun/gforgame/common/errors"
	"github.com/forfun/gforgame/internal/config"
	"github.com/forfun/gforgame/internal/constants"
	"github.com/forfun/gforgame/internal/context"
	configdomain "github.com/forfun/gforgame/internal/domain/config"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/events"
	"github.com/forfun/gforgame/internal/io"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/reward"
)

// 签到模块
type SignInService struct {
}

func NewSignInService() *SignInService {
	return &SignInService{}
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

func (s *SignInService) SignIn(player *playerdomain.Player) *commonerrors.BusinessError {
    nthDay := int32(time.Now().Day())
    monthlyResetBox := player.MonthlyReset
    if IntSliceContains(monthlyResetBox.SignInDays, nthDay) {
        return errors.NewBusinessError(constants.I18N_COMMON_ILLEGAL_PARAMS)
    }
    monthlyResetBox.SignInDays = append(monthlyResetBox.SignInDays, nthDay);

    signinData := config.QueryById[configdomain.SigninData](nthDay)
    if (signinData == nil) {
        return errors.NewBusinessError(constants.I18N_COMMON_ILLEGAL_PARAMS)
    }
    rewards := reward.ParseReward(signinData.Rewards)
    rewards.Reward(player, constants.ActionType_Signin)
    context.EventBus.Publish(events.PlayerEntityChange, player)
    return nil
}

func (s *SignInService) SignInMakeUp(player *playerdomain.Player, day int32) *commonerrors.BusinessError {
    monthlyResetBox := player.MonthlyReset
    if _, ok := monthlyResetBox.SignInMakeUp[day]; ok {
        return errors.NewBusinessError(constants.I18N_COMMON_ILLEGAL_PARAMS)
    }
	nthDay := int32(time.Now().Day())
	if day < 1 || day > nthDay {
        return errors.NewBusinessError(constants.I18N_COMMON_ILLEGAL_PARAMS)
    }
    monthlyResetBox.SignInMakeUp[day] = day
	context.EventBus.Publish(events.PlayerEntityChange, player)
    return nil
}

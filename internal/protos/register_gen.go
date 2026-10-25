// 该文件为程序自动生成，请勿手动修改

package protos

import (
	"github.com/forfun/gforgame/network/protocol"
)

// RegisterMessages 将生成的协议注册到指定实例；每个工厂调用一次。
func RegisterMessages(factory *protocol.MessageFactory) {
	// ----from activity.go----
	factory.RegisterMessage(1651, &PushActivityLoadAll{})

	// ----from catalog.go----
	factory.RegisterMessage(3101, &ReqCatalogReward{})
	factory.RegisterMessage(3151, &ResCatalogReward{})
	factory.RegisterMessage(3198, &PushCatalogAdd{})
	factory.RegisterMessage(3199, &PushCatalogInfo{})

	// ----from chat.go----
	factory.RegisterMessage(1801, &ReqChat{})
	factory.RegisterMessage(1851, &ResChat{})
	factory.RegisterMessage(1899, &PushChatNewMessage{})

	// ----from friend.go----
	factory.RegisterMessage(1902, &ReqFriendSearchPlayers{})
	factory.RegisterMessage(1903, &ReqFriendQueryMyFriends{})
	factory.RegisterMessage(1904, &ReqFriendApply{})
	factory.RegisterMessage(1905, &ReqFriendDealApplyRecord{})
	factory.RegisterMessage(1906, &ReqFriendDelete{})
	factory.RegisterMessage(1952, &ResFriendSearchPlayers{})
	factory.RegisterMessage(1953, &ResFriendQueryMyFriends{})
	factory.RegisterMessage(1954, &ResFriendApply{})
	factory.RegisterMessage(1955, &ResFriendDealApplyRecord{})
	factory.RegisterMessage(1956, &ResFriendDelete{})
	factory.RegisterMessage(1997, &PushFriendInfo{})

	// ----from gm.go----
	factory.RegisterMessage(-251, &ResGmCommand{})
	factory.RegisterMessage(-201, &ReqGmCommand{})

	// ----from hero.go----
	factory.RegisterMessage(801, &ReqHeroRecruit{})
	factory.RegisterMessage(802, &ReqHeroUpFight{})
	factory.RegisterMessage(803, &ReqHeroOffFight{})
	factory.RegisterMessage(804, &ReqHeroLevelUp{})
	factory.RegisterMessage(805, &ReqHeroUpStage{})
	factory.RegisterMessage(807, &ReqHeroCombine{})
	factory.RegisterMessage(808, &ReqHeroChangePosition{})
	factory.RegisterMessage(851, &ResHeroRecruit{})
	factory.RegisterMessage(852, &ResHeroUpFight{})
	factory.RegisterMessage(853, &ResHeroOffFight{})
	factory.RegisterMessage(854, &ResHeroLevelUp{})
	factory.RegisterMessage(855, &ResHeroUpStage{})
	factory.RegisterMessage(857, &PushAllHeroInfo{})
	factory.RegisterMessage(858, &ResHeroCombine{})
	factory.RegisterMessage(859, &ResHeroChangePosition{})
	factory.RegisterMessage(5006, &PushHeroAdd{})
	factory.RegisterMessage(5007, &PushHeroAttrChange{})

	// ----from item.go----
	factory.RegisterMessage(250, &PushBackpackInfo{})
	factory.RegisterMessage(253, &PushItemChanged{})
	factory.RegisterMessage(4002, &PushPurseInfo{})

	// ----from mail.go----
	factory.RegisterMessage(501, &ReqMailRead{})
	factory.RegisterMessage(502, &ReqMailGetReward{})
	factory.RegisterMessage(504, &ReqMailGetAllRewards{})
	factory.RegisterMessage(505, &ReqMailDeleteAll{})
	factory.RegisterMessage(551, &ResMailRead{})
	factory.RegisterMessage(552, &ResMailGetReward{})
	factory.RegisterMessage(554, &ResMailGetAllRewards{})
	factory.RegisterMessage(555, &ResMailDeleteAll{})
	factory.RegisterMessage(599, &PushMailAll{})
	factory.RegisterMessage(6009, &ReqMailReadAll{})
	factory.RegisterMessage(6010, &ResMailReadAll{})

	// ----from mall.go----
	factory.RegisterMessage(1101, &ReqMallBuy{})
	factory.RegisterMessage(1151, &ResMallBuy{})

	// ----from mixture.go----
	factory.RegisterMessage(9906, &ReqClientUploadEvent{})
	factory.RegisterMessage(9956, &ResClientUploadEvent{})

	// ----from monthcard.go----
	factory.RegisterMessage(2102, &ReqMonthCardGetReward{})
	factory.RegisterMessage(2152, &ResMonthCardGetReward{})
	factory.RegisterMessage(2198, &PushMonthCardInfo{})

	// ----from player.go----
	factory.RegisterMessage(103, &ReqPlayerLogin{})
	factory.RegisterMessage(105, &ReqPlayerUpLevel{})
	factory.RegisterMessage(106, &ReqPlayerUpStage{})
	factory.RegisterMessage(109, &ReqEditClientData{})
	factory.RegisterMessage(110, &ReqPlayerRefreshScore{})
	factory.RegisterMessage(152, &ResPlayerUpLevel{})
	factory.RegisterMessage(153, &ResPlayerUpStage{})
	factory.RegisterMessage(154, &ResPlayerLogin{})
	factory.RegisterMessage(155, &PushLoadComplete{})
	factory.RegisterMessage(156, &PushDailyResetInfo{})
	factory.RegisterMessage(171, &ResEditClientData{})
	factory.RegisterMessage(172, &ResPlayerRefreshScore{})
	factory.RegisterMessage(2005, &ReqPlayerLoadingFinish{})
	factory.RegisterMessage(2008, &PushPlayerFightChange{})

	// ----from quest.go----
	factory.RegisterMessage(701, &ReqQuestTakeReward{})
	factory.RegisterMessage(702, &ReqQuestTakeProgressReward{})
	factory.RegisterMessage(706, &ReqQuestTakeAllRewards{})
	factory.RegisterMessage(753, &ResQuestTakeProgressReward{})
	factory.RegisterMessage(754, &ResQuestTakeReward{})
	factory.RegisterMessage(760, &PushQuestReplace{})
	factory.RegisterMessage(762, &ResQuestTakeAllRewards{})
	factory.RegisterMessage(791, &PushAchievementInfo{})
	factory.RegisterMessage(795, &PushQuestRefreshVo{})
	factory.RegisterMessage(797, &PushQuestAutoTakeReward{})
	factory.RegisterMessage(798, &PushQuestDailyInfo{})
	factory.RegisterMessage(799, &PushQuestWeeklyInfo{})

	// ----from rank.go----
	factory.RegisterMessage(7001, &ReqRankQuery{})
	factory.RegisterMessage(7002, &ResRankQuery{})

	// ----from recharge.go----
	factory.RegisterMessage(2298, &PushRechargeInfo{})
	factory.RegisterMessage(2299, &PushRechargePay{})

	// ----from sigin.go----
	factory.RegisterMessage(3001, &ReqSignIn{})
	factory.RegisterMessage(3002, &ReqSignInMakeup{})
	factory.RegisterMessage(3051, &ResSignIn{})
	factory.RegisterMessage(3052, &ResSignInMakeup{})
	factory.RegisterMessage(3099, &PushSigninInfo{})

	// ----from system.go----
	factory.RegisterMessage(-152, &ResGetServerTime{})
	factory.RegisterMessage(-151, &ResHeartBeat{})
	factory.RegisterMessage(-102, &ReqGetServerTime{})
	factory.RegisterMessage(-101, &ReqHeartBeat{})

}

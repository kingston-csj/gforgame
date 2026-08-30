// 该文件为程序自动生成，请勿手动修改

package protos

import "github.com/forfun/gforgame/network/protocol"

func init() {
	// ----from activity.go----
	protocol.RegisterMessage(1651, &PushActivityLoadAll{})

	// ----from catalog.go----
	protocol.RegisterMessage(3101, &ReqCatalogReward{})
	protocol.RegisterMessage(3151, &ResCatalogReward{})
	protocol.RegisterMessage(3198, &PushCatalogAdd{})
	protocol.RegisterMessage(3199, &PushCatalogInfo{})

	// ----from chat.go----
	protocol.RegisterMessage(1801, &ReqChat{})
	protocol.RegisterMessage(1851, &ResChat{})
	protocol.RegisterMessage(1899, &PushChatNewMessage{})

	// ----from friend.go----
	protocol.RegisterMessage(1902, &ReqFriendSearchPlayers{})
	protocol.RegisterMessage(1903, &ReqFriendQueryMyFriends{})
	protocol.RegisterMessage(1904, &ReqFriendApply{})
	protocol.RegisterMessage(1905, &ReqFriendDealApplyRecord{})
	protocol.RegisterMessage(1906, &ReqFriendDelete{})
	protocol.RegisterMessage(1952, &ResFriendSearchPlayers{})
	protocol.RegisterMessage(1953, &ResFriendQueryMyFriends{})
	protocol.RegisterMessage(1954, &ResFriendApply{})
	protocol.RegisterMessage(1955, &ResFriendDealApplyRecord{})
	protocol.RegisterMessage(1956, &ResFriendDelete{})
	protocol.RegisterMessage(1997, &PushFriendInfo{})

	// ----from gm.go----
	protocol.RegisterMessage(-251, &ResGmCommand{})
	protocol.RegisterMessage(-201, &ReqGmCommand{})

	// ----from hero.go----
	protocol.RegisterMessage(801, &ReqHeroRecruit{})
	protocol.RegisterMessage(802, &ReqHeroUpFight{})
	protocol.RegisterMessage(803, &ReqHeroOffFight{})
	protocol.RegisterMessage(804, &ReqHeroLevelUp{})
	protocol.RegisterMessage(805, &ReqHeroUpStage{})
	protocol.RegisterMessage(807, &ReqHeroCombine{})
	protocol.RegisterMessage(808, &ReqHeroChangePosition{})
	protocol.RegisterMessage(851, &ResHeroRecruit{})
	protocol.RegisterMessage(852, &ResHeroUpFight{})
	protocol.RegisterMessage(853, &ResHeroOffFight{})
	protocol.RegisterMessage(854, &ResHeroLevelUp{})
	protocol.RegisterMessage(855, &ResHeroUpStage{})
	protocol.RegisterMessage(857, &PushAllHeroInfo{})
	protocol.RegisterMessage(858, &ResHeroCombine{})
	protocol.RegisterMessage(859, &ResHeroChangePosition{})
	protocol.RegisterMessage(5006, &PushHeroAdd{})
	protocol.RegisterMessage(5007, &PushHeroAttrChange{})

	// ----from item.go----
	protocol.RegisterMessage(250, &PushBackpackInfo{})
	protocol.RegisterMessage(253, &PushItemChanged{})
	protocol.RegisterMessage(4002, &PushPurseInfo{})

	// ----from mail.go----
	protocol.RegisterMessage(501, &ReqMailRead{})
	protocol.RegisterMessage(502, &ReqMailGetReward{})
	protocol.RegisterMessage(504, &ReqMailGetAllRewards{})
	protocol.RegisterMessage(505, &ReqMailDeleteAll{})
	protocol.RegisterMessage(551, &ResMailRead{})
	protocol.RegisterMessage(552, &ResMailGetReward{})
	protocol.RegisterMessage(554, &ResMailGetAllRewards{})
	protocol.RegisterMessage(555, &ResMailDeleteAll{})
	protocol.RegisterMessage(599, &PushMailAll{})
	protocol.RegisterMessage(6009, &ReqMailReadAll{})
	protocol.RegisterMessage(6010, &ResMailReadAll{})

	// ----from mall.go----
	protocol.RegisterMessage(1101, &ReqMallBuy{})
	protocol.RegisterMessage(1151, &ResMallBuy{})

	// ----from mixture.go----
	protocol.RegisterMessage(9906, &ReqClientUploadEvent{})
	protocol.RegisterMessage(9956, &ResClientUploadEvent{})

	// ----from monthcard.go----
	protocol.RegisterMessage(2102, &ReqMonthCardGetReward{})
	protocol.RegisterMessage(2152, &ResMonthCardGetReward{})
	protocol.RegisterMessage(2198, &PushMonthCardInfo{})

	// ----from player.go----
	protocol.RegisterMessage(103, &ReqPlayerLogin{})
	protocol.RegisterMessage(105, &ReqPlayerUpLevel{})
	protocol.RegisterMessage(106, &ReqPlayerUpStage{})
	protocol.RegisterMessage(109, &ReqEditClientData{})
	protocol.RegisterMessage(110, &ReqPlayerRefreshScore{})
	protocol.RegisterMessage(152, &ResPlayerUpLevel{})
	protocol.RegisterMessage(153, &ResPlayerUpStage{})
	protocol.RegisterMessage(154, &ResPlayerLogin{})
	protocol.RegisterMessage(155, &PushLoadComplete{})
	protocol.RegisterMessage(156, &PushDailyResetInfo{})
	protocol.RegisterMessage(171, &ResEditClientData{})
	protocol.RegisterMessage(172, &ResPlayerRefreshScore{})
	protocol.RegisterMessage(2005, &ReqPlayerLoadingFinish{})

	// ----from quest.go----
	protocol.RegisterMessage(701, &ReqQuestTakeReward{})
	protocol.RegisterMessage(702, &ReqQuestTakeProgressReward{})
	protocol.RegisterMessage(706, &ReqQuestTakeAllRewards{})
	protocol.RegisterMessage(707, &ReqQuestEntrust{})
	protocol.RegisterMessage(753, &ResQuestTakeProgressReward{})
	protocol.RegisterMessage(754, &ResQuestTakeReward{})
	protocol.RegisterMessage(760, &PushQuestReplace{})
	protocol.RegisterMessage(762, &ResQuestTakeAllRewards{})
	protocol.RegisterMessage(763, &ResQuestEntrust{})
	protocol.RegisterMessage(791, &PushAchievementInfo{})
	protocol.RegisterMessage(795, &PushQuestRefreshVo{})
	protocol.RegisterMessage(797, &PushQuestAutoTakeReward{})
	protocol.RegisterMessage(798, &PushQuestDailyInfo{})
	protocol.RegisterMessage(799, &PushQuestWeeklyInfo{})

	// ----from rank.go----
	protocol.RegisterMessage(7001, &ReqRankQuery{})
	protocol.RegisterMessage(7002, &ResRankQuery{})

	// ----from recharge.go----
	protocol.RegisterMessage(2298, &PushRechargeInfo{})
	protocol.RegisterMessage(2299, &PushRechargePay{})

	// ----from sigin.go----
	protocol.RegisterMessage(3001, &ReqSignIn{})
	protocol.RegisterMessage(3002, &ReqSignInMakeup{})
	protocol.RegisterMessage(3051, &ResSignIn{})
	protocol.RegisterMessage(3052, &ResSignInMakeup{})
	protocol.RegisterMessage(3099, &PushSigninInfo{})

	// ----from system.go----
	protocol.RegisterMessage(-152, &ResGetServerTime{})
	protocol.RegisterMessage(-151, &ResHeartBeat{})
	protocol.RegisterMessage(-102, &ReqGetServerTime{})
	protocol.RegisterMessage(-101, &ReqHeartBeat{})

	// ----from transfer.go----
	protocol.RegisterMessage(-352, &ResServerLogin{})
	protocol.RegisterMessage(-302, &ReqServerLogin{})
	protocol.RegisterMessage(-300, &TransferGateToLogic{})

}

// 该文件为程序自动生成，请勿手动修改

package protos

import (
	"io/github/gforgame/network"
)

func init() {
	// ----from activity.go----
	network.RegisterMessage(1651, &PushActivityLoadAll{})

	// ----from catalog.go----
	network.RegisterMessage(3101, &ReqCatalogReward{})
	network.RegisterMessage(3151, &ResCatalogReward{})
	network.RegisterMessage(3198, &PushCatalogAdd{})
	network.RegisterMessage(3199, &PushCatalogInfo{})

	// ----from chat.go----
	network.RegisterMessage(1801, &ReqChat{})
	network.RegisterMessage(1851, &ResChat{})
	network.RegisterMessage(1899, &PushChatNewMessage{})

	// ----from friend.go----
	network.RegisterMessage(1902, &ReqFriendSearchPlayers{})
	network.RegisterMessage(1903, &ReqFriendQueryMyFriends{})
	network.RegisterMessage(1904, &ReqFriendApply{})
	network.RegisterMessage(1905, &ReqFriendDealApplyRecord{})
	network.RegisterMessage(1906, &ReqFriendDelete{})
	network.RegisterMessage(1952, &ResFriendSearchPlayers{})
	network.RegisterMessage(1953, &ResFriendQueryMyFriends{})
	network.RegisterMessage(1954, &ResFriendApply{})
	network.RegisterMessage(1955, &ResFriendDealApplyRecord{})
	network.RegisterMessage(1956, &ResFriendDelete{})
	network.RegisterMessage(1997, &PushFriendInfo{})

	// ----from gm.go----
	network.RegisterMessage(-251, &ResGmCommand{})
	network.RegisterMessage(-201, &ReqGmCommand{})

	// ----from hero.go----
	network.RegisterMessage(801, &ReqHeroRecruit{})
	network.RegisterMessage(802, &ReqHeroUpFight{})
	network.RegisterMessage(803, &ReqHeroOffFight{})
	network.RegisterMessage(804, &ReqHeroLevelUp{})
	network.RegisterMessage(805, &ReqHeroUpStage{})
	network.RegisterMessage(807, &ReqHeroCombine{})
	network.RegisterMessage(808, &ReqHeroChangePosition{})
	network.RegisterMessage(851, &ResHeroRecruit{})
	network.RegisterMessage(852, &ResHeroUpFight{})
	network.RegisterMessage(853, &ResHeroOffFight{})
	network.RegisterMessage(854, &ResHeroLevelUp{})
	network.RegisterMessage(855, &ResHeroUpStage{})
	network.RegisterMessage(857, &PushAllHeroInfo{})
	network.RegisterMessage(858, &ResHeroCombine{})
	network.RegisterMessage(859, &ResHeroChangePosition{})
	network.RegisterMessage(5006, &PushHeroAdd{})
	network.RegisterMessage(5007, &PushHeroAttrChange{})

	// ----from item.go----
	network.RegisterMessage(250, &PushBackpackInfo{})
	network.RegisterMessage(253, &PushItemChanged{})
	network.RegisterMessage(4002, &PushPurseInfo{})

	// ----from mail.go----
	network.RegisterMessage(501, &ReqMailRead{})
	network.RegisterMessage(502, &ReqMailGetReward{})
	network.RegisterMessage(504, &ReqMailGetAllRewards{})
	network.RegisterMessage(505, &ReqMailDeleteAll{})
	network.RegisterMessage(551, &ResMailRead{})
	network.RegisterMessage(552, &ResMailGetReward{})
	network.RegisterMessage(554, &ResMailGetAllRewards{})
	network.RegisterMessage(555, &ResMailDeleteAll{})
	network.RegisterMessage(599, &PushMailAll{})
	network.RegisterMessage(6009, &ReqMailReadAll{})
	network.RegisterMessage(6010, &ResMailReadAll{})

	// ----from mall.go----
	network.RegisterMessage(1101, &ReqMallBuy{})
	network.RegisterMessage(1151, &ResMallBuy{})

	// ----from mixture.go----
	network.RegisterMessage(9902, &ReqIdleViewReward{})
	network.RegisterMessage(9903, &ReqIdleGetReward{})
	network.RegisterMessage(9906, &ReqClientUploadEvent{})
	network.RegisterMessage(9953, &ResIdleGetReward{})
	network.RegisterMessage(9956, &ResClientUploadEvent{})
	network.RegisterMessage(9999, &PushIdleInfo{})

	// ----from monthcard.go----
	network.RegisterMessage(2102, &ReqMonthCardGetReward{})
	network.RegisterMessage(2152, &ResMonthCardGetReward{})
	network.RegisterMessage(2198, &PushMonthCardInfo{})

	// ----from player.go----
	network.RegisterMessage(103, &ReqPlayerLogin{})
	network.RegisterMessage(105, &ReqPlayerUpLevel{})
	network.RegisterMessage(106, &ReqPlayerUpStage{})
	network.RegisterMessage(109, &ReqEditClientData{})
	network.RegisterMessage(110, &ReqPlayerRefreshScore{})
	network.RegisterMessage(152, &ResPlayerUpLevel{})
	network.RegisterMessage(153, &ResPlayerUpStage{})
	network.RegisterMessage(154, &ResPlayerLogin{})
	network.RegisterMessage(155, &PushLoadComplete{})
	network.RegisterMessage(156, &PushDailyResetInfo{})
	network.RegisterMessage(171, &ResEditClientData{})
	network.RegisterMessage(172, &ResPlayerRefreshScore{})
	network.RegisterMessage(1800, &ReqJoinRoom{})
	network.RegisterMessage(2003, &ReqPlayerCreate{})
	network.RegisterMessage(2004, &ResPlayerCreate{})
	network.RegisterMessage(2005, &ReqPlayerLoadingFinish{})
	network.RegisterMessage(2008, &PushPlayerFightChange{})

	// ----from quest.go----
	network.RegisterMessage(701, &ReqQuestTakeReward{})
	network.RegisterMessage(702, &ReqQuestTakeProgressReward{})
	network.RegisterMessage(706, &ReqQuestTakeAllRewards{})
	network.RegisterMessage(707, &ReqQuestEntrust{})
	network.RegisterMessage(753, &ResQuestTakeProgressReward{})
	network.RegisterMessage(754, &ResQuestTakeReward{})
	network.RegisterMessage(760, &PushQuestReplace{})
	network.RegisterMessage(762, &ResQuestTakeAllRewards{})
	network.RegisterMessage(763, &ResQuestEntrust{})
	network.RegisterMessage(791, &PushAchievementInfo{})
	network.RegisterMessage(795, &PushQuestRefreshVo{})
	network.RegisterMessage(797, &PushQuestAutoTakeReward{})
	network.RegisterMessage(798, &PushQuestDailyInfo{})
	network.RegisterMessage(799, &PushQuestWeeklyInfo{})

	// ----from rank.go----
	network.RegisterMessage(7001, &ReqRankQuery{})
	network.RegisterMessage(7002, &ResRankQuery{})

	// ----from recharge.go----
	network.RegisterMessage(2298, &PushRechargeInfo{})
	network.RegisterMessage(2299, &PushRechargePay{})

	// ----from scene.go----
	network.RegisterMessage(2801, &ResSceneGetData{})
	network.RegisterMessage(2802, &ReqSceneSetData{})
	network.RegisterMessage(2851, &ReqSceneGetData{})
	network.RegisterMessage(2852, &ResSceneSetData{})

	// ----from sigin.go----
	network.RegisterMessage(3001, &ReqSignIn{})
	network.RegisterMessage(3002, &ReqSignInMakeup{})
	network.RegisterMessage(3051, &ResSignIn{})
	network.RegisterMessage(3052, &ResSignInMakeup{})
	network.RegisterMessage(3099, &PushSigninInfo{})

	// ----from system.go----
	network.RegisterMessage(-152, &ResGetServerTime{})
	network.RegisterMessage(-151, &ResHeartBeat{})
	network.RegisterMessage(-102, &ReqGetServerTime{})
	network.RegisterMessage(-101, &ReqHeartBeat{})

}

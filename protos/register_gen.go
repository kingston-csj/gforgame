package protos

import (
	"io/github/gforgame/network"
)

func init() {
	// ----from chat.go----
	network.RegisterMessage(1899, &PushChatNewMessage{})
	network.RegisterMessage(1801, &ReqChat{})
	network.RegisterMessage(1851, &ResChat{})

	// ----from player.go----
	network.RegisterMessage(103, &ReqPlayerLogin{})
	network.RegisterMessage(2005, &ReqPlayerLoadingFinish{})
	network.RegisterMessage(154, &ResPlayerLogin{})
	network.RegisterMessage(2003, &ReqPlayerCreate{})
	network.RegisterMessage(2004, &ResPlayerCreate{})
	network.RegisterMessage(1800, &ReqJoinRoom{})
	network.RegisterMessage(105, &ReqPlayerUpLevel{})
	network.RegisterMessage(152, &ResPlayerUpLevel{})
	network.RegisterMessage(2008, &PushPlayerFightChange{})
	network.RegisterMessage(106, &ReqPlayerUpStage{})
	network.RegisterMessage(153, &ResPlayerUpStage{})
	network.RegisterMessage(155, &PushLoadComplete{})
	network.RegisterMessage(156, &PushDailyResetInfo{})

	// ----from scene.go----
	network.RegisterMessage(2851, &ReqSceneGetData{})
	network.RegisterMessage(2802, &ReqSceneSetData{})
	network.RegisterMessage(2801, &ResSceneGetData{})
	network.RegisterMessage(2852, &ResSceneSetData{})

	// ----from friend.go----
	network.RegisterMessage(1997, &PushFriendInfo{})
	network.RegisterMessage(1902, &ReqFriendSearchPlayers{})
	network.RegisterMessage(1904, &ReqFriendApply{})
	network.RegisterMessage(1905, &ReqFriendDealApplyRecord{})
	network.RegisterMessage(1906, &ReqFriendDelete{})
	network.RegisterMessage(1903, &ReqFriendQueryMyFriends{})
	network.RegisterMessage(1954, &ResFriendApply{})
	network.RegisterMessage(1955, &ResFriendDealApplyRecord{})
	network.RegisterMessage(1956, &ResFriendDelete{})
	network.RegisterMessage(1953, &ResFriendQueryMyFriends{})
	network.RegisterMessage(1952, &ResFriendSearchPlayers{})

	// ----from hero.go----
	network.RegisterMessage(5001, &ReqHeroRecruit{})
	network.RegisterMessage(5002, &ResHeroRecruit{})
	network.RegisterMessage(857, &PushAllHeroInfo{})
	network.RegisterMessage(5004, &ReqHeroLevelUp{})
	network.RegisterMessage(5005, &ResHeroLevelUp{})
	network.RegisterMessage(5008, &ReqHeroUpStage{})
	network.RegisterMessage(5009, &ResHeroUpStage{})
	network.RegisterMessage(5006, &PushHeroAdd{})
	network.RegisterMessage(5007, &PushHeroAttrChange{})
	network.RegisterMessage(5010, &ReqHeroCombine{})
	network.RegisterMessage(5011, &ResHeroCombine{})
	network.RegisterMessage(5012, &ReqHeroUpFight{})
	network.RegisterMessage(5013, &ResHeroUpFight{})
	network.RegisterMessage(5014, &ReqHeroOffFight{})
	network.RegisterMessage(5015, &ResHeroOffFight{})
	network.RegisterMessage(5016, &ReqHeroChangePosition{})
	network.RegisterMessage(5017, &ResHeroChangePosition{})

	// ----from item.go----
	network.RegisterMessage(250, &PushBackpackInfo{})
	network.RegisterMessage(4002, &PushPurseInfo{})
	network.RegisterMessage(253, &PushItemChanged{})

	// ----from recharge.go----
	network.RegisterMessage(2299, &PushRechargePay{})
	network.RegisterMessage(2298, &PushRechargeInfo{})

	// ----from system.go----
	network.RegisterMessage(-101, &ReqHeartBeat{})
	network.RegisterMessage(-151, &ResHeartBeat{})
	network.RegisterMessage(-102, &ReqGetServerTime{})
	network.RegisterMessage(-152, &ResGetServerTime{})

	// ----from gm.go----
	network.RegisterMessage(-201, &ReqGmCommand{})
	network.RegisterMessage(-251, &ResGmCommand{})

	// ----from mail.go----
	network.RegisterMessage(599, &PushMailAll{})
	network.RegisterMessage(504, &ReqMailGetAllRewards{})
	network.RegisterMessage(554, &ResMailGetAllRewards{})
	network.RegisterMessage(501, &ReqMailRead{})
	network.RegisterMessage(551, &ResMailRead{})
	network.RegisterMessage(502, &ReqMailGetReward{})
	network.RegisterMessage(552, &ResMailGetReward{})
	network.RegisterMessage(505, &ReqMailDeleteAll{})
	network.RegisterMessage(555, &ResMailDeleteAll{})
	network.RegisterMessage(6009, &ReqMailReadAll{})
	network.RegisterMessage(6010, &ResMailReadAll{})

	// ----from mall.go----
	network.RegisterMessage(1101, &ReqMallBuy{})
	network.RegisterMessage(1151, &ResMallBuy{})

	// ----from quest.go----
	network.RegisterMessage(797, &PushQuestAutoTakeReward{})
	network.RegisterMessage(798, &PushQuestDailyInfo{})
	network.RegisterMessage(795, &PushQuestRefreshVo{})
	network.RegisterMessage(760, &PushQuestReplace{})
	network.RegisterMessage(799, &PushQuestWeeklyInfo{})
	network.RegisterMessage(791, &PushAchievementInfo{})
	network.RegisterMessage(706, &ReqQuestTakeAllRewards{})
	network.RegisterMessage(702, &ReqQuestTakeProgressReward{})
	network.RegisterMessage(701, &ReqQuestTakeReward{})
	network.RegisterMessage(762, &ResQuestTakeAllRewards{})
	network.RegisterMessage(753, &ResQuestTakeProgressReward{})
	network.RegisterMessage(754, &ResQuestTakeReward{})

	// ----from rank.go----
	network.RegisterMessage(7001, &ReqRankQuery{})
	network.RegisterMessage(7002, &ResRankQuery{})

	// ----from sigin.go----
	network.RegisterMessage(3099, &PushSigninInfo{})
	network.RegisterMessage(3001, &ReqSignIn{})
	network.RegisterMessage(3002, &ResSignIn{})

}

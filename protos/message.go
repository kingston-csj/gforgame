package protos

const (
	CmdHeartBeatReq = -101
	CmdHeartBeatRes = -151

	CmdGetServerTimeReq = -102
	CmdGetServerTimeRes = -152

	CmdGmReqCommand = -201
	CmdGmResCommand = -251
		
	CmdPlayerReqLogin         = 103
	CmdPlayerReqUpLevel       = 105
	CmdPlayerReqUpStage       = 106
	CmdPlayerResLogin         = 154
	CmdPlayerPushLoadComplete = 155
	CmdPlayerResUpLevel       = 152
	CmdPlayerResUpStage       = 153
	CmdPlayerPushDailyResetInfo = 156


	CmdPlayerReqCreate        = 2003
	CmdPlayerResCreate        = 2004
	CmdPlayerReqLoadingFinish = 2005
	CmdPlayerPushFightChange  = 2008

	

	CmdGmReqAction = 3001
	CmdGmResAction = 3002

	CmdItemPushBackpackInfo = 250
	CmdItemResPurseInfo    = 4002
	CmdItemPushChanged     = 253

	CmdHeroReqRecruit        = 5001
	CmdHeroResRecruit        = 5002
	CmdHeroPushAllHero        = 857
	CmdHeroReqLevelUp        = 5004
	CmdHeroResLevelUp        = 5005
	CmdHeroPushAdd           = 5006
	CmdHeroPushAttrChange    = 5007
	CmdHeroReqUpStage        = 5008
	CmdHeroResUpStage        = 5009
	CmdHeroReqCombine        = 5010
	CmdHeroResCombine        = 5011
	CmdHeroReqUpFight        = 5012
	CmdHeroResUpFight        = 5013
	CmdHeroReqOffFight       = 5014
	CmdHeroResOffFight       = 5015
	CmdHeroReqChangePosition = 5016
	CmdHeroResChangePosition = 5017
	
	CmdMailReqRead         = 501
	CmdMailReqGetReward    = 502
	CmdMailReqGetAllReward = 504
	CmdMailReqDeleteAll    = 505
	
	CmdMailResRead         = 551	
	CmdMailResGetReward    = 552
	CmdMailResGetAllReward = 554
	CmdMailResDeleteAll    = 555
	CmdMailReqReadAll      = 6009
	CmdMailResReadAll      = 6010
	CmdMailPushAll         = 599


	
	CMD_RES_REPLACE_QUEST      = 760
	CMD_REQ_QUEST_ALL_REWARD      = 706
	CMD_REQ_QUEST_PROGRESS_REWARD      = 702
	CMD_REQ_QUEST_REWARD      = 701

	CMD_RES_QUEST_ALL_REWARD = 762
	CMD_RES_QUEST_PROGRESS_REWARD = 753
	CMD_RES_QUEST_REWARD = 754
	CMD_PUSH_ACHIEVEMENT      = 791
	CMD_PUSH_UPDATE_QUEST      = 795
	CMD_PUSH_QUEST_AUTO_REWARD = 797
	CMD_PUSH_DAILY_QUEST      = 798
	CMD_PUSH_WEEKLY_QUEST      = 799


	CmdMallPushInfo = 1199
	CmdReqMallBuy = 1101
	CmdResMallBuy = 1151

	CmdChaJoinRoom = 1800
	CmdChatReqChat = 1801
	CmdChatResChat = 1851
	CmdChatPushNew = 1899

	CmdFriendReqSearchPlayers = 1902
	CmdFriendReqQueryFriends  = 1903
	CmdFriendReqApply         = 1904
	CmdFriendReqDealApply     = 1905
	CmdFriendReqDelete        = 1906
	CmdFriendResSearchPlayers = 1952
	CmdFriendResQueryFriends  = 1953
	CmdFriendResApply         = 1954
	CmdFriendResDealApply     = 1955
	CmdFriendResDelete        = 1956
	CmdFriendPushApplyList    = 1997

	CmdPushRechargePayInfo = 2298
	CmdPushRechargePay = 2299

	CmdRankReqQuery = 7001
	CmdRankResQuery = 7002


	CmdSceneReqGetData = 2851
	CmdSceneReqSetData = 2802

	CmdSceneResGetData = 2801
	CmdSceneResSetData = 2852

	CmdSignInReqSignIn = 3001
	CmdSignInResSignIn = 3002
	CmdSignInPush = 3099

)

package protos

const (
	CmdPlayerReqLogin         = 2001
	CmdPlayerResLogin         = 2002
	CmdPlayerReqCreate        = 2003
	CmdPlayerResCreate        = 2004
	CmdPlayerReqLoadingFinish = 2005
	CmdPlayerReqUpLevel       = 2006
	CmdPlayerResUpLevel       = 2007
	CmdPlayerPushFightChange  = 2008
	CmdPlayerReqUpStage       = 2009
	CmdPlayerResUpStage       = 2010

	CmdGmReqAction = 3001
	CmdGmResAction = 3002

	CmdItemResBackpackInfo = 4001
	CmdItemResPurseInfo    = 4002
	CmdItemPushChanged     = 4003

	CmdHeroReqRecruit        = 5001
	CmdHeroResRecruit        = 5002
	CmdHeroResAllHero        = 5003
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

	CmdMailReqGetAllReward = 6001
	CmdMailResGetAllReward = 6002
	CmdMailReqRead         = 6003
	CmdMailResRead         = 6004
	CmdMailReqGetReward    = 6005
	CmdMailResGetReward    = 6006
	CmdMailReqDeleteAll    = 6007
	CmdMailResDeleteAll    = 6008
	CmdMailReqReadAll      = 6009
	CmdMailResReadAll      = 6010
	CmdMailPushAll         = 6011

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

	CmdRankReqQuery = 7001
	CmdRankResQuery = 7002
)

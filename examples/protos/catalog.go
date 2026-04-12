package protos

type CatalogModel struct { //图鉴模型
	UnlockIds   []int32 `json:"unlockIds"`   //已解锁id列表
	ReceivedIds []int32 `json:"receivedIds"` //已领取id列表
}

type ReqCatalogReward struct { //图鉴——领取奖励
	_    struct{} `cmd_ref:"CmdReqCatalogReward"`
	Id   int32    `json:"id"`   //图鉴id
	Type int32    `json:"type"` //类型
}

type ResCatalogReward struct { //图鉴——领取奖励
	_         struct{}    `cmd_ref:"CmdResCatalogReward"`
	RewardVos []*RewardVo `json:"rewardVos"` //奖励列表
	Code      int32       `json:"code"`      //状态码
}

type PushCatalogInfo struct { //图鉴——主界面
	_            struct{}     `cmd_ref:"CmdPushCatalogInfo"`
	SitemCatalog CatalogModel `json:"sitemCatalog"` //场景道具
	ItemCatalog  CatalogModel `json:"itemCatalog"`  //物品模型
	MenuCatalog  CatalogModel `json:"menuCatalog"`  //菜单模型
}

type PushCatalogAdd struct { //图鉴——新增
	_      struct{} `cmd_ref:"CmdPushCatalogAdd"`
	Typ    int32    `json:"type"`   //类型
	ItemId int32    `json:"itemId"` //id
}
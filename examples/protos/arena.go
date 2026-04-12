package protos

type ChallengeRecord struct {
	OpponentId       string `json:"opponentId"`       // 对手id
	OpponentName     string `json:"opponentName"`     // 对手名称
	OpponentHead     int32  `json:"opponentHead"`     // 对手头像
	OpponentFighting int32  `json:"opponentFighting"` // 对手战力
	ChallengeTime    int64  `json:"challengeTime"`    // 挑战时间
	Winner           string `json:"winner"`           // 获胜方id
	Score            int32  `json:"score"`            // 得分，获胜为正数， 失败为负数
}

type MatchMember struct {
	PlayerId   string `json:"playerId"`   // 玩家id
	PlayerName string `json:"playerName"` // 玩家名称
	Score      int32  `json:"score"`      // 玩家得分
	Head       int32  `json:"head"`       // 玩家头像
	Fighting   int32  `json:"fighting"`   // 玩家战力
}

type ReqArenaApplyFight struct { //竞技场——申请挑战
	_        struct{} `cmd_ref:"CmdReqArenaApplyFight"`
	TargetId string   `json:"targetId"` // 目标玩家id
}

type ReqArenaEditDefenseTeam struct { //竞技场——编辑防守队伍
	_    struct{} `cmd_ref:"CmdReqArenaEditDefenseTeam"`
	Team []int32  `json:"team"` // 防守队伍，英雄id列表
}

type ReqArenaFightEnd struct { //竞技场——战斗结束
	_        struct{} `cmd_ref:"CmdReqArenaFightEnd"`
	TargetId string   `json:"targetId"` // 目标玩家id
	IsWinner bool     `json:"isWinner"` // 是否获胜
}

type ReqArenaQueryChallengeRecord struct { //竞技场——查询挑战记录
	_ struct{} `cmd_ref:"CmdReqArenaQueryChallengeRecord"`
}

type ReqArenaQueryDefenseTeam struct { //竞技场——查询防守队伍
	_ struct{} `cmd_ref:"CmdReqArenaQueryDefenseTeam"`
}

type ReqArenaQueryMatchList struct { //竞技场——查询匹配列表
	_ struct{} `cmd_ref:"CmdReqArenaQueryMatchList"`
}

type ReqArenaRefreshMatchList struct { //竞技场——刷新匹配列表
	_ struct{} `cmd_ref:"CmdReqArenaRefreshMatchList"`
}

type ResArenaApplyFight struct { //竞技场——申请挑战结果
	_    struct{} `cmd_ref:"CmdResArenaApplyFight"`
	Code int32    `json:"code"` // 结果码
}

type ResArenaEditDefenseTeam struct { //竞技场——编辑防守队伍结果
	_    struct{} `cmd_ref:"CmdResArenaEditDefenseTeam"`
	Code int32    `json:"code"` // 结果码
}

type ResArenaFightEnd struct { //竞技场——战斗结束结果
	_                  struct{} `cmd_ref:"CmdResArenaFightEnd"`
	MyInitScore        int32    `json:"myInitScore"`        // 我的初始分
	MyChangedScore     int32    `json:"myChangedScore"`     // 我的得分变化
	TargetInitScore    int32    `json:"targetInitScore"`    // 目标初始分
	TargetChangedScore int32    `json:"targetChangedScore"` // 目标得分变化
	Code               int32    `json:"code"`               // 结果码
}

type ResArenaQueryChallengeRecord struct { //竞技场——查询挑战记录结果
	_       struct{}          `cmd_ref:"CmdResArenaQueryChallengeRecord"`
	Records []ChallengeRecord `json:"records"` // 挑战记录列表
	Code    int32             `json:"code"`    // 结果码
}

type ResArenaQueryDefenseTeam struct { //竞技场——查询防守队伍结果
	_    struct{} `cmd_ref:"CmdResArenaQueryDefenseTeam"`
	Team []int32  `json:"team"` // 防守队伍，英雄id列表
	Code int32    `json:"code"` // 结果码
}

type ResArenaQueryMatchList struct { //竞技场——查询匹配列表结果
	_       struct{}      `cmd_ref:"CmdResArenaQueryMatchList"`
	Matches []MatchMember `json:"matches"` // 匹配列表
	Ticket  string        `json:"ticket"`  // 门票数量
	Code    int32         `json:"code"`    // 结果码
}

type ResArenaRefreshMatchList struct { //竞技场——刷新匹配列表结果
	_       struct{}      `cmd_ref:"CmdResArenaRefreshMatchList"`
	Members []MatchMember `json:"members"` // 匹配列表
	Code    int32         `json:"code"`    // 结果码
}
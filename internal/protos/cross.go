package protos

type TransferGateToLogic struct {
	_        struct{} `cmd_ref:"CmdTransferMsgGateToLogic"`
	PlayerId string   `json:"playerId"`
	Cmd      int32    `json:"cmd"`
	Index    int32    `json:"index"`
	Body     []byte   `json:"body"`
}

func (t *TransferGateToLogic) GetPlayerID() string {
	if t == nil {
		return ""
	}
	return t.PlayerId
}

func (t *TransferGateToLogic) GetTransferCmd() int32 {
	if t == nil {
		return 0
	}
	return t.Cmd
}

func (t *TransferGateToLogic) GetTransferBody() []byte {
	if t == nil {
		return nil
	}
	return t.Body
}

type ReqServerLogin struct { // 服务节点登录请求
	_        struct{} `cmd_ref:"CmdReqServerLogin"`
	ServerId int32    `json:"serverId"`
}

type ResServerLogin struct { // 服务节点登录响应
	_    struct{} `cmd_ref:"CmdResServerLogin"`
	Code int32    `json:"code"`
}

type NotifyPlayerLogoutToGame struct { // 网关-游戏服：玩家退出游戏请求
	_        struct{} `cmd_ref:"CmdNotifyPlayerLogoutToGame"`
	PlayerId string   `json:"playerId"`
}

type NotifyOnlinePlayerToGame struct { // 网关-游戏服：游戏服重启时，批量同步在线玩家到游戏服
	_         struct{} `cmd_ref:"CmdNotifyOnlinePlayerToGame"`
	PlayerIds []string `json:"playerIds"`
}

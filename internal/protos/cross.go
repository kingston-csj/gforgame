package protos

type TransferGateToLogic struct {
	_        struct{} `cmd_ref:"CmdTransferMsgGateToLogic" scope:"rpc"`
	PlayerId string   `json:"playerId"`
	Cmd      int32    `json:"cmd"`
	Index    int32    `json:"index"`
	Body     []byte   `json:"body"`
}

func (*TransferGateToLogic) RPCCommand() int32 { return CmdTransferMsgGateToLogic }

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
	_        struct{} `cmd_ref:"CmdReqServerLogin" scope:"rpc"`
	ServerId int32    `json:"serverId"`
}

type ResServerLogin struct { // 服务节点登录响应
	_    struct{} `cmd_ref:"CmdResServerLogin" scope:"rpc"`
	Code int32    `json:"code"`
}

type NotifyPlayerLogoutToGame struct { // 网关-游戏服：玩家退出游戏请求
	_        struct{} `cmd_ref:"CmdNotifyPlayerLogoutToGame" scope:"rpc"`
	PlayerId string   `json:"playerId"`
}

func (*NotifyPlayerLogoutToGame) RPCCommand() int32 { return CmdNotifyPlayerLogoutToGame }

type NotifyOnlinePlayerToGame struct { // 网关-游戏服：游戏服重启时，批量同步在线玩家到游戏服
	_         struct{} `cmd_ref:"CmdNotifyOnlinePlayerToGame" scope:"rpc"`
	PlayerIds []string `json:"playerIds"`
}

func (*NotifyOnlinePlayerToGame) RPCCommand() int32 { return CmdNotifyOnlinePlayerToGame }

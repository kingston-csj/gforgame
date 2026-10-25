package protos

import (
	"github.com/forfun/gforgame/network/protocol"
	"github.com/forfun/gforgame/rpc"
)

const (
	CmdGateRPCResult = -353
	CmdGateRPCError  = -354
)

// GateRPCResult 确认网关请求已被接收；业务回包仍通过玩家推送发送。
type GateRPCResult struct{}

// RegisterGateRPCMessages 只注册网关 RPC 的公共应答和错误协议。
func RegisterGateRPCMessages(factory *protocol.MessageFactory) {
	factory.RegisterMessage(CmdGateRPCResult, &GateRPCResult{})
	factory.RegisterMessage(CmdGateRPCError, &rpc.ErrorResponse{})
}

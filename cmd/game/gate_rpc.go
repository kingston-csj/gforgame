package main

import (
	"fmt"

	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/route"
	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/protocol"
)

// newGateRPCService adapts gateway requests to the existing player Actor pipeline.
// RPC listener setup and lifecycle are managed independently by rpc.TCPServer.
func newGateRPCService(service *route.GateRoute, tasks *GameTaskHandler, factory *protocol.MessageFactory) *route.GateRoute {
	transform := NewGateTransformHandler(factory)
	service.SetTransferHandler(func(session network.Session, msg *protos.TransferGateToLogic) error {
		frame := &protocol.RequestDataFrame{Msg: msg}
		frame.Header.Cmd = protos.CmdTransferMsgGateToLogic
		if !transform.MessageReceived(session, frame) {
			return fmt.Errorf("invalid gate transfer, cmd=%d", msg.Cmd)
		}
		if !tasks.MessageReceived(session, frame) {
			return fmt.Errorf("gate transfer rejected, cmd=%d", msg.Cmd)
		}
		return nil
	})
	return service
}

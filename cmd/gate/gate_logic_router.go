package main

import (
	"fmt"
	"reflect"

	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/common/util/jsonutil"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/protocol"
)

// 作为逻辑层，接收logic层的推送
type LogicRouter struct {
	router *network.MessageRoute
}

func (g *LogicRouter) MessageReceived(session *network.Session, frame *protocol.RequestDataFrame) bool {
	defer func() {
		if r := recover(); r != nil {
			logger.ErrorNoStack(fmt.Errorf("panic recovered: %v", r))
		}
	}()
	if frame.Header.Cmd == protos.CmdTransferMsgGateToLogic {
		transferResp := &protos.TransferGateToLogic{}
		switch body := frame.Msg.(type) {
		case *protos.TransferGateToLogic:
			transferResp = body
		case []byte:
			if err := jsonutil.JsonBytesToStruct(body, transferResp); err != nil {
				logger.ErrorNoStack(fmt.Errorf("decode transfer response failed: %v", err))
				return false
			}
		default:
			logger.ErrorNoStack(fmt.Errorf("transfer response payload type invalid: %T", frame.Msg))
			return false
		}
		if err := forwardTransferToClient(transferResp); err != nil {
			logger.ErrorNoStack(err)
			return false
		}
		return true
	}
	return true
}

func forwardTransferToClient(transfer *protos.TransferGateToLogic) error {
	if transfer.PlayerId == "" {
		return fmt.Errorf("transfer response playerId is empty, cmd=%d", transfer.Cmd)
	}
	clientSession := network.GetSessionByPlayerId(transfer.PlayerId)
	if clientSession == nil {
		return fmt.Errorf("client session not found, playerId=%s cmd=%d", transfer.PlayerId, transfer.Cmd)
	}
	typ, _ := network.GetMessageType(transfer.Cmd)
	if typ == nil {
		return fmt.Errorf("transfer response message type not found, cmd=%d", transfer.Cmd)
	}
	resp := reflect.New(typ.Elem()).Interface()
	if err := gateMsgCodec.Decode(transfer.Body, resp); err != nil {
		return fmt.Errorf("decode transfer response body failed, cmd=%d err=%v", transfer.Cmd, err)
	}
	if err := clientSession.Send(resp, transfer.Index); err != nil {
		return fmt.Errorf("send transfer response to client failed, playerId=%s cmd=%d err=%v", transfer.PlayerId, transfer.Cmd, err)
	}
	jsonStr, err := jsonutil.StructToJSON(resp)
	if err != nil {
		return fmt.Errorf("encode transfer response body failed, cmd=%d err=%v", transfer.Cmd, err)
	}
	logger.Info(fmt.Sprintf("send transfer response to client: cmd %d, 内容: %s", transfer.Cmd, jsonStr))
	return nil
}

func newLogicIoDispatcher() network.IoDispatch {
	router := network.NewMessageRoute()
	ioDispatcher := &network.BaseIoDispatch{}
	ioDispatcher.AddHandler(&LogicRouter{router: router})
	return ioDispatcher
}

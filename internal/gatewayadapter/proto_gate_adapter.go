package gatewayadapter

import (
	"fmt"

	"github.com/forfun/gforgame/common/util/jsonutil"
	"github.com/forfun/gforgame/gateway/contract"
	"github.com/forfun/gforgame/internal/protos"
)

type ProtoGateProtocolAdapter struct{}

type protoLoginRequest struct {
	req *protos.ReqPlayerLogin
}

func (r protoLoginRequest) GetPlayerID() string {
	if r.req == nil {
		return ""
	}
	return r.req.PlayerId
}

func (r protoLoginRequest) GetServerID() int32 {
	if r.req == nil {
		return 0
	}
	return r.req.ServerId
}

type protoTransferMessage struct {
	transfer *protos.TransferGateToLogic
}

func (m protoTransferMessage) GetPlayerID() string {
	if m.transfer == nil {
		return ""
	}
	return m.transfer.PlayerId
}

func (m protoTransferMessage) GetTransferCmd() int32 {
	if m.transfer == nil {
		return 0
	}
	return m.transfer.Cmd
}

func (m protoTransferMessage) GetTransferIndex() int32 {
	if m.transfer == nil {
		return 0
	}
	return m.transfer.Index
}

func (m protoTransferMessage) GetTransferBody() []byte {
	if m.transfer == nil {
		return nil
	}
	return m.transfer.Body
}

func NewProtoClientLoginAdapter() contract.ClientLoginAdapter {
	return &ProtoGateProtocolAdapter{}
}

func NewProtoTransferCodec() contract.TransferCodec {
	return &ProtoGateProtocolAdapter{}
}

func (a *ProtoGateProtocolAdapter) LoginCmd() int32 {
	return protos.CmdReqPlayerLogin
}

func (a *ProtoGateProtocolAdapter) TransferCmd() int32 {
	return protos.CmdTransferMsgGateToLogic
}

func (a *ProtoGateProtocolAdapter) DecodeLoginRequest(body []byte) (contract.GateLoginRequest, error) {
	req := &protos.ReqPlayerLogin{}
	if err := jsonutil.JsonBytesToStruct(body, req); err != nil {
		return nil, err
	}
	return protoLoginRequest{req: req}, nil
}

func (a *ProtoGateProtocolAdapter) NewReplacingLoginPush() any {
	return &protos.PushReplacingLogin{}
}

func (a *ProtoGateProtocolAdapter) NewTransferMessage(playerID string, cmd int32, index int32, body []byte) any {
	return &protos.TransferGateToLogic{
		PlayerId: playerID,
		Cmd:      cmd,
		Index:    index,
		Body:     body,
	}
}

func (a *ProtoGateProtocolAdapter) ParseTransferMessage(msg any) (contract.GateTransferMessage, error) {
	switch body := msg.(type) {
	case *protos.TransferGateToLogic:
		return protoTransferMessage{transfer: body}, nil
	case []byte:
		transfer := &protos.TransferGateToLogic{}
		if err := jsonutil.JsonBytesToStruct(body, transfer); err != nil {
			return nil, err
		}
		return protoTransferMessage{transfer: transfer}, nil
	default:
		return nil, fmt.Errorf("transfer response payload type invalid: %T", msg)
	}
}

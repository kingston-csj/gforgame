package transfer

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/forfun/gforgame/codec/json"
	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/common/util/jsonutil"
	"github.com/forfun/gforgame/internal/context"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/protocol"
)

type TransferService struct{}

var (
	instance *TransferService
	once     sync.Once
	msgCodec = json.NewSerializer()
)

func GetTransferService() *TransferService {
	once.Do(func() {
		instance = &TransferService{}
	})
	return instance
}

func (s *TransferService) TransferGateToLogic(session *network.Session, transfer *protos.TransferGateToLogic) {
	typ, _ := network.GetMessageType(transfer.Cmd)
	if typ == nil {
		logger.ErrorNoStack(fmt.Errorf("message type not found: %d", transfer.Cmd))
		return
	}
	msg := reflect.New(typ.Elem()).Interface()
	if err := msgCodec.Decode(transfer.Body, msg); err != nil {
		logger.ErrorNoStack(fmt.Errorf("decode transfer body failed, cmd=%d err=%v", transfer.Cmd, err))
		return
	}

	// 复用当前逻辑服的真实路由分发表进行二次分发
	router := getLogicRouter()
	if router == nil {
		logger.ErrorNoStack(fmt.Errorf("logic router not ready, cmd=%d", transfer.Cmd))
		return
	}
	frame := &protocol.RequestDataFrame{
		Header: protocol.MessageHeader{
			Cmd:      transfer.Cmd,
			Index:    transfer.Index,
			Payload: transfer.PlayerId,
		},
		Msg: msg,
	}
	msgName, _ := network.GetMsgName(frame.Header.Cmd)
	jsonStr, err := jsonutil.StructToJSON(frame.Msg)
	if err == nil {
		if strings.Index(msgName, "HeartBeat") == -1 {
			logger.Info(fmt.Sprintf("接收消息: cmd:%d, name:%s, 内容:%s",  frame.Header.Cmd, msgName, jsonStr))
		}
	}
	msgHandler, _ := router.GetHandler(frame.Header.Cmd)
	if msgHandler == nil {
		logger.ErrorNoStack(fmt.Errorf("transfer target handler is nil: %d", frame.Header.Cmd))
		return
	}
	args := network.BuildHandlerArgs(msgHandler, session, frame.Header.Index, frame.Msg, frame.Header.Payload)
	values := msgHandler.Method.Func.Call(args)
	if len(values) > 0 {
		resp := values[0].Interface()
		respCmd, err := network.GetMessageCmd(resp)
		if err != nil {
			logger.ErrorNoStack(fmt.Errorf("resolve transfer response cmd failed, reqCmd=%d err=%v", frame.Header.Cmd, err))
			return
		}
		respBody, err := msgCodec.Encode(resp)
		if err != nil {
			logger.ErrorNoStack(fmt.Errorf("encode transfer response body failed, respCmd=%d err=%v", respCmd, err))
			return
		}
		respTransfer := &protos.TransferGateToLogic{
			PlayerId: transfer.PlayerId,
			Cmd:      respCmd,
			Index:    transfer.Index,
			Body:     respBody,
		}
		if err := session.Send(respTransfer, frame.Header.Index); err != nil {
			logger.ErrorNoStack(fmt.Errorf("send transfer response failed, reqCmd=%d respCmd=%d err=%v", frame.Header.Cmd, respCmd, err))
			return
		}
	}
}

func getLogicRouter() *network.MessageRoute {
	if context.TcpServer != nil && context.TcpServer.Router != nil {
		return context.TcpServer.Router
	}
	if context.WsServer != nil && context.WsServer.Router != nil {
		return context.WsServer.Router
	}
	return nil
}

package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/forfun/gforgame/codec"
	"github.com/forfun/gforgame/codec/json"
	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/common/util/jsonutil"
	serverconfig "github.com/forfun/gforgame/config"
	"github.com/forfun/gforgame/internal/constants"
	"github.com/forfun/gforgame/internal/io"
	"github.com/forfun/gforgame/internal/protos"
	protocolValidator "github.com/forfun/gforgame/internal/validator"
	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/protocol"
)

var logicResponseCodec = json.NewSerializer()

// GateTransformHandler 将网关消息TransferGateToLogic转换为实际的消息
type GateTransformHandler struct {
	codec codec.MessageCodec
}

func NewGateTransformHandler() *GateTransformHandler {
	return &GateTransformHandler{
		codec: json.NewSerializer(),
	}
}

func (g *GateTransformHandler) MessageReceived(session *network.Session, frame *protocol.RequestDataFrame) bool {
	if frame.Header.Cmd != protos.CmdTransferMsgGateToLogic {
		return true
	}
	transfer, ok := frame.Msg.(*protos.TransferGateToLogic)
	if !ok || transfer == nil {
		logger.ErrorNoStack(fmt.Errorf("invalid transfer message: %T", frame.Msg))
		return false
	}
	typ, _ := network.GetMessageType(transfer.Cmd)
	if typ == nil {
		logger.ErrorNoStack(fmt.Errorf("message type not found: %d", transfer.Cmd))
		return false
	}
	msg := reflect.New(typ.Elem()).Interface()
	if err := g.codec.Decode(transfer.Body, msg); err != nil {
		logger.ErrorNoStack(fmt.Errorf("decode transfer body failed, cmd=%d err=%v", transfer.Cmd, err))
		return false
	}
	frame.Header.Cmd = transfer.Cmd
	frame.Header.Index = transfer.Index
	frame.Header.Payload = transfer.PlayerId
	frame.Msg = msg
	return true
}

type GameTaskHandler struct {
	router *network.MessageRoute
}

func (g *GameTaskHandler) MessageReceived(session *network.Session, frame *protocol.RequestDataFrame) bool {
	defer func() {
		if r := recover(); r != nil {
			logger.ErrorNoStack(fmt.Errorf("panic recovered: %v", r))
		}
	}()
	
	// 补齐 playerId，后续路由和回包都依赖它。
	fillFramePayloadFromSession(session, frame)
	// 先定位消息处理器，找不到就直接终止当前消息。
	msgHandler := g.getMessageHandler(frame)
	if msgHandler == nil {
		return false
	}
	// 验证协议的参数，如果校验失败，直接返回错误码
	if msgHandler.NeedValidate {
		validationErrors := protocolValidator.ValidateStruct(frame.Msg)
		if len(validationErrors) > 0 {
			errMsg := protocolValidator.FormatValidationErrors(validationErrors)
			if errMsg != "" {
				// logger.Info(fmt.Sprintf("validation failed for cmd=%d: %s", frame.Header.Cmd, errMsg))
				if resp, ok := buildErrorResponse(msgHandler, constants.I18N_COMMON_PROTOCOL_VALIDATION_FAILED); ok {
					session.Send(resp, frame.Header.Index)
				}
				return false
			}
		}
	}
	// 直连模式下打印入站消息，便于本地排查。
	logInboundMessage(session, frame)
	// 优先走代码生成分发，失败时再回退到反射调用。
	return g.dispatchMessage(session, frame, msgHandler)
}

func sendResponse(session *network.Session, frame *protocol.RequestDataFrame, resp any) error {
	if resp == nil {
		return nil
	}
	if !serverconfig.ServerConfig.UseGateMode || frame.Header.Payload == "" {
		return session.Send(resp, frame.Header.Index)
	}
	
	io.NotifyByPlayerId(frame.Header.Payload, frame.Header.Index, resp)
	return nil
}

func fillFramePayloadFromSession(session *network.Session, frame *protocol.RequestDataFrame) {
	if frame.Header.Payload != "" {
		return
	}
	if id, ok := session.GetAttr("id"); ok {
		if sid, ok := id.(string); ok {
			frame.Header.Payload = sid
		}
	}
}

func (g *GameTaskHandler) getMessageHandler(frame *protocol.RequestDataFrame) *network.Handler {
	msgHandler, _ := g.router.GetHandler(frame.Header.Cmd)
	if msgHandler == nil {
		logger.ErrorNoStack(fmt.Errorf("msgHandler is nil: %v", frame.Header.Cmd))
	}
	return msgHandler
}

func logInboundMessage(session *network.Session, frame *protocol.RequestDataFrame) {
	msgName, _ := network.GetMsgName(frame.Header.Cmd)
	jsonStr, err := jsonutil.StructToJSON(frame.Msg)
	if err != nil || strings.Index(msgName, "HeartBeat") != -1 {
		return
	}
	id := "anonymous"
	if serverconfig.ServerConfig.UseGateMode {
		id = frame.Header.Payload
	} else {
		id, ok := session.GetAttr("id")
		if ok {
			id = id.(string)
		}
	}
	
	logger.Info(fmt.Sprintf("[%s] 接收消息: cmd:%d, name:%s, 内容:%s", id, frame.Header.Cmd, msgName, jsonStr))
}

func (g *GameTaskHandler) dispatchMessage(session *network.Session, frame *protocol.RequestDataFrame, msgHandler *network.Handler) bool {
	resp, handled, dispatchErr, panicErr := callGeneratedRouteHandlerSafely(frame.Header.Cmd, msgHandler, session, frame.Header.Index, frame.Msg, frame.Header.Payload)
	if handled {
		if panicErr != nil {
			return handleRoutePanic(session, frame, msgHandler, "generated route handler panic", panicErr)
		}
		if dispatchErr == nil {
			return sendHandlerResponse(session, frame, resp)
		}
		// 静态分发失败时回退反射调用，保证兼容性。
		logger.ErrorNoStack(fmt.Errorf("generated dispatch failed, fallback to reflect: cmd=%d err=%v", frame.Header.Cmd, dispatchErr))
	}

	// 反射分发只作为兜底路径，避免生成代码缺失时消息直接丢失。
	args := buildHandlerArgs(msgHandler, session, frame.Header.Index, frame.Msg, frame.Header.Payload)
	values, panicErr := callRouteHandlerSafely(msgHandler, args)
	if panicErr != nil {
		return handleRoutePanic(session, frame, msgHandler, "route handler panic", panicErr)
	}
	if len(values) == 0 {
		return true
	}
	return sendHandlerResponse(session, frame, values[0].Interface())
}

func buildHandlerArgs(msgHandler *network.Handler, session *network.Session, index int32, msg any, playerID string) []reflect.Value {
	args := make([]reflect.Value, 0, 5)
	args = append(args, msgHandler.Receiver)
	if msgHandler.HasPlayer {
		args = append(args, reflect.ValueOf(playerID))
	}
	if msgHandler.HasSession {
		args = append(args, reflect.ValueOf(session))
	}
	if msgHandler.Indindexed {
		args = append(args, reflect.ValueOf(index))
	}
	args = append(args, reflect.ValueOf(msg))
	return args
}

func handleRoutePanic(session *network.Session, frame *protocol.RequestDataFrame, msgHandler *network.Handler, title string, panicErr error) bool {
	logger.Error(fmt.Sprintf("%s: cmd=%d method=%s", title, frame.Header.Cmd, msgHandler.Method.Name), panicErr)
	if errorResp, ok := buildErrorResponse(msgHandler, constants.I18N_COMMON_INTERNAL_ERROR); ok {
		if err := sendResponse(session, frame, errorResp); err != nil {
			logger.Error("send error response failed: %v", err)
		}
	}
	return false
}

func sendHandlerResponse(session *network.Session, frame *protocol.RequestDataFrame, resp any) bool {
	if err := sendResponse(session, frame, resp); err != nil {
		logger.Error("send response failed: %v", err)
		return false
	}
	return true
}

func callGeneratedRouteHandlerSafely(cmd int32, msgHandler *network.Handler, session *network.Session, index int32, msg any, playerID string) (resp any, handled bool, dispatchErr error, panicErr error) {
	invoker, ok := getGeneratedRouteInvoker(cmd)
	if !ok {
		return nil, false, nil, nil
	}
	handled = true
	defer func() {
		if r := recover(); r != nil {
			panicErr = logger.PanicToError(r)
		}
	}()
	resp, dispatchErr = invoker(msgHandler, playerID, session, index, msg)
	return resp, handled, dispatchErr, nil
}

func getGeneratedRouteInvoker(cmd int32) (generatedRouteInvoker, bool) {
	if generatedRouteDispatchers == nil {
		return nil, false
	}
	invoker, ok := generatedRouteDispatchers[cmd]
	return invoker, ok
}

func callRouteHandlerSafely(msgHandler *network.Handler, args []reflect.Value) (values []reflect.Value, panicErr error) {
	defer func() {
		if r := recover(); r != nil {
			panicErr = logger.PanicToError(r)
		}
	}()
	values = msgHandler.Method.Func.Call(args)
	return values, nil
}

func buildErrorResponse(msgHandler *network.Handler, code int32) (any, bool) {
	mt := msgHandler.Method.Type
	if mt.NumOut() == 0 {
		return nil, false
	}
	outType := mt.Out(0)
	if outType.Kind() != reflect.Ptr || outType.Elem().Kind() != reflect.Struct {
		return nil, false
	}
	resp := reflect.New(outType.Elem())
	codeField := resp.Elem().FieldByName("Code")
	if !codeField.IsValid() || !codeField.CanSet() || codeField.Kind() != reflect.Int32 {
		return nil, false
	}
	codeField.SetInt(int64(code))
	return resp.Interface(), true
}

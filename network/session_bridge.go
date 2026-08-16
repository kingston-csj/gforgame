package network

import (
	"net"
	"reflect"

	"github.com/forfun/gforgame/codec"
	"github.com/forfun/gforgame/network/protocol"
	sessionpkg "github.com/forfun/gforgame/network/session"
)

type Session = sessionpkg.Session

type PayloadMode = sessionpkg.PayloadMode

const (
	PayloadModeDecode  = sessionpkg.PayloadModeDecode
	PayloadModeRawBody = sessionpkg.PayloadModeRawBody
)

type SerialSessionLoopOptions = sessionpkg.SerialSessionLoopOptions
type DispatchSessionLoopOptions = sessionpkg.DispatchSessionLoopOptions

func HashSessionWorkerIndex(sessionKey string, workerCount int) int {
	return sessionpkg.HashSessionWorkerIndex(sessionKey, workerCount)
}

func ResolveWorkerIndex(ioFrame *protocol.RequestDataFrame, fallbackSessionIdx int, workerCount int) int {
	return sessionpkg.ResolveWorkerIndex(ioFrame, fallbackSessionIdx, workerCount)
}

func ServeSession(session Session, ioDispatch IoDispatch, run func(session Session)) {
	sessionpkg.ServeSession(session, ioDispatch, run)
}

func ServeSessionConn(conn net.Conn, messageCodec codec.MessageCodec, ioDispatch IoDispatch, dispatchWorkers int32, payloadMode PayloadMode) {
	sessionpkg.ServeSessionConn(conn, messageCodec, ioDispatch, dispatchWorkers, payloadMode)
}

const DefaultDirectSessionIdleTimeout = sessionpkg.DefaultDirectSessionIdleTimeout

type messageResolverBridge struct{}

func (messageResolverBridge) GetMessageCmd(msg any) (int32, error) {
	return GetMessageCmd(msg)
}

func (messageResolverBridge) GetMsgName(cmd int32) (string, error) {
	return GetMsgName(cmd)
}

func (messageResolverBridge) GetMessageType(cmd int32) (reflect.Type, error) {
	return GetMessageType(cmd)
}

func init() {
	sessionpkg.SetMessageResolver(messageResolverBridge{})
}

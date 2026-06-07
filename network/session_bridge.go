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

func NewSession(conn net.Conn, messageCodec codec.MessageCodec) *Session {
	return sessionpkg.NewSession(conn, messageCodec)
}

func NewSessionWithProtocol(conn net.Conn, messageCodec codec.MessageCodec, protocolType protocol.ProtocolType) *Session {
	return sessionpkg.NewSessionWithProtocol(conn, messageCodec, protocolType)
}

func RegisterSession(conn net.Conn, s *Session) {
	sessionpkg.RegisterSession(conn, s)
}

func GetSession(conn net.Conn) *Session {
	return sessionpkg.GetSession(conn)
}

func UnregisterSession(conn net.Conn) {
	sessionpkg.UnregisterSession(conn)
}

func CloseAllSessions() {
	sessionpkg.CloseAllSessions()
}

func AddSession(session *Session, playerID string) {
	sessionpkg.AddSession(session, playerID)
}

func RemoveSession(session *Session) {
	sessionpkg.RemoveSession(session)
}

func GetPlayerIDBySession(session *Session) (string, bool) {
	return sessionpkg.GetPlayerIDBySession(session)
}

func GetSessionByPlayerId(playerID string) *Session {
	return sessionpkg.GetSessionByPlayerId(playerID)
}

func GetAllSessions() []*Session {
	return sessionpkg.GetAllSessions()
}

func GetAllOnlinePlayerIds() []string {
	return sessionpkg.GetAllOnlinePlayerIds()
}

func GetAllOnlinePlayerSessions() []*Session {
	return sessionpkg.GetAllOnlinePlayerSessions()
}

func IsOnline(playerID string) bool {
	return sessionpkg.IsOnline(playerID)
}

func HashSessionWorkerIndex(sessionKey string, workerCount int) int {
	return sessionpkg.HashSessionWorkerIndex(sessionKey, workerCount)
}

func ResolveWorkerIndex(ioFrame *protocol.RequestDataFrame, fallbackSessionIdx int, workerCount int) int {
	return sessionpkg.ResolveWorkerIndex(ioFrame, fallbackSessionIdx, workerCount)
}

func ServeSession(session *Session, ioDispatch IoDispatch, run func(session *Session)) {
	sessionpkg.ServeSession(session, ioDispatch, run)
}

func ServeSessionConn(conn net.Conn, messageCodec codec.MessageCodec, ioDispatch IoDispatch, dispatchWorkers int32, payloadMode PayloadMode) {
	sessionpkg.ServeSessionConn(conn, messageCodec, ioDispatch, dispatchWorkers, payloadMode)
}

func RunDispatchSessionLoop(session *Session, ioDispatch IoDispatch, opts *DispatchSessionLoopOptions) {
	sessionpkg.RunDispatchSessionLoop(session, ioDispatch, opts)
}

func RunSerialSessionLoop(session *Session, ioDispatch IoDispatch, opts *SerialSessionLoopOptions) {
	sessionpkg.RunSerialSessionLoop(session, ioDispatch, opts)
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

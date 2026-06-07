package session

import (
	"encoding/json"
	"hash/fnv"

	serverconfig "github.com/forfun/gforgame/config"
	"github.com/forfun/gforgame/network/protocol"
)

// HashSessionWorkerIndex 根据 session 或玩家维度计算 worker 下标。
func HashSessionWorkerIndex(sessionKey string, workerCount int) int {
	if workerCount <= 1 {
		return 0
	}
	h := fnv.New32a()
	_, _ = h.Write([]byte(sessionKey))
	return int(h.Sum32() % uint32(workerCount))
}

// ResolveWorkerIndex 统一决定当前消息应该投递到哪个 worker。
func ResolveWorkerIndex(ioFrame *protocol.RequestDataFrame, fallbackSessionIdx int, workerCount int) int {
	if workerCount <= 1 {
		return 0
	}
	if serverconfig.ServerConfig.UseGateMode {
		playerID := extractPlayerID(ioFrame)
		if playerID != "" {
			return HashSessionWorkerIndex(playerID, workerCount)
		}
	}
	return fallbackSessionIdx
}

func extractPlayerID(ioFrame *protocol.RequestDataFrame) string {
	if ioFrame == nil {
		return ""
	}
	if ioFrame.Header.Payload != "" {
		return ioFrame.Header.Payload
	}
	msg := ioFrame.Msg
	if msg == nil {
		return ""
	}
	if carrier, ok := msg.(playerIDCarrier); ok {
		if playerID := carrier.GetPlayerID(); playerID != "" {
			return playerID
		}
	}
	if playerID := extractPlayerIDFromJSON(msg); playerID != "" {
		return playerID
	}
	return ""
}

type playerIDCarrier interface{ GetPlayerID() string }

func extractPlayerIDFromJSON(msg any) string {
	var raw []byte
	switch v := msg.(type) {
	case []byte:
		raw = v
	case string:
		raw = []byte(v)
	default:
		return ""
	}
	if len(raw) == 0 {
		return ""
	}
	var payload struct {
		PlayerId string `json:"playerId"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return ""
	}
	return payload.PlayerId
}

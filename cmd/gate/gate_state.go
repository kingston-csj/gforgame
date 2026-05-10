package main

import (
	"sync"

	"github.com/forfun/gforgame/codec/json"
	"github.com/forfun/gforgame/network"
)

type backendPool struct {
	serverID     int32
	addr         string
	session      *network.Session
	reconnecting bool
}

var (
	gateMsgCodec        = json.NewSerializer()
	logicIoDispatcher   network.IoDispatch
	backendPools        = make(map[int32]*backendPool)
	backendPoolsMu      sync.RWMutex
	playerServerIDMap   = make(map[string]int32)
	playerServerIDMapMu sync.RWMutex
	outboundQueue       chan *backendOutboundMsg
	outboundNotify      chan struct{}
	outboundStop        chan struct{}
	outboundWg          sync.WaitGroup
)

const (
	logicServerType       uint32 = 1
	backendReconnectDelay        = 1500
)

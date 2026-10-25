package main

import (
	"sync"

	serverconfig "github.com/forfun/gforgame/config"

	"github.com/forfun/gforgame/codec/json"
	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/dispatch"
	"github.com/forfun/gforgame/rpc"
)

type backendPool struct {
	serverID     int32
	addr         string
	session      network.Session
	reconnecting bool
}

func (p *backendPool) IsConnected() bool {
	return p.session != nil && p.session.IsAlive()
}

var (
	gateMsgCodec        = json.NewSerializer()
	logicIoDispatcher   dispatch.IoDispatch
	gameRPCClient       *rpc.Client
	backendPools        = make(map[int32]*backendPool)
	backendPoolsMu      sync.RWMutex
	serverDiscoveryStop = make(chan struct{})
	backendMonitorStop  = make(chan struct{})
	playerServerIDMap   = make(map[string]int32)
	playerServerIDMapMu sync.RWMutex
	outboundQueue       chan *backendOutboundMsg
	outboundNotify      chan struct{}
	outboundStop        chan struct{}
	outboundWg          sync.WaitGroup
)

const (
	logicServerType        uint32 = uint32(serverconfig.NodeTypeGame)
	backendReconnectDelay         = 1500
	backendMonitorInterval        = 3 * 1000
)

package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/forfun/gforgame/common/container/set"
	"github.com/forfun/gforgame/common/logger"
	serverconfig "github.com/forfun/gforgame/config"
	"github.com/forfun/gforgame/gateway/contract"
	"github.com/forfun/gforgame/internal/gatewayadapter"
	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/ws"
)

var (
	// 用于对日志进行限流，防止日志泛滥
	logFilter = set.NewExpireSet(10 * time.Minute)
	// 第一阶段登录协议适配器，先由当前项目协议提供实现。
	gateLoginAdapter contract.ClientLoginAdapter = gatewayadapter.NewProtoClientLoginAdapter()
	// 第一阶段转发协议编解码器，先由当前项目协议提供实现。
	gateTransferCodec contract.TransferCodec = gatewayadapter.NewProtoTransferCodec()
)

func main() {
	logicIoDispatcher = newLogicIoDispatcher()
	onlyLocalDiscovery, ok := serverconfig.GetExtraBool("discovery.onlylocal")
	if !ok {
		onlyLocalDiscovery = true
	}

	var err error
	if onlyLocalDiscovery {
		err = startLocalServerDiscovery()
	} else {
		err = startServerDiscoveryHeartbeat()
	}
	if err != nil {
		panic(err)
	}
	startBackendSessionMonitor()
	startOutboundDispatcher()
	router := network.NewMessageRoute()
	ioDispatcher := &MyMessageDispatch{}
	ioDispatcher.AddHandler(&ClientRouter{router: router})
	codec := gateMsgCodec
	node := ws.NewServer(
		ws.WithAddress(serverconfig.ServerConfig.ServerUrl),
		ws.WithRouter(router),
		ws.WithIoDispatch(ioDispatcher),
		ws.WithCodec(codec),
		ws.WithPayloadMode(network.PayloadModeRawBody),
	)
	err = node.Start()
	if err != nil {
		panic(err)
	}

	sg := make(chan os.Signal)
	signal.Notify(sg, os.Interrupt, os.Kill)
	select {
	case sig := <-sg:
		logger.Info(fmt.Sprintf("game server is closing (signal: %v)", sig))
	case <-node.RunningChan():
		logger.Info(fmt.Sprintf("game server is closing (signal: http)"))
	}
	// 执行所有关服逻辑
	close(serverDiscoveryStop)
	close(backendMonitorStop)
	node.Stop()
	stopOutboundDispatcher()
	closeAllBackendSessions()
	logFilter.Close()
}

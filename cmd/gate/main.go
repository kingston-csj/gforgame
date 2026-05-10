package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/forfun/gforgame/common/logger"
	serverconfig "github.com/forfun/gforgame/config"
	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/ws"
)

func main() {
	logicIoDispatcher = newLogicIoDispatcher()
	startOutboundDispatcher()
	router := network.NewMessageRoute()
	ioDispatcher := &MyMessageDispatch{}
	ioDispatcher.AddHandler(&ClientRouter{router: router})
	codec := gateMsgCodec
	modules := []network.Module{}
	node := ws.NewServer(
		ws.WithAddress(serverconfig.ServerConfig.ServerUrl),
		ws.WithRouter(router),
		ws.WithIoDispatch(ioDispatcher),
		ws.WithCodec(codec),
		ws.WithModules(modules...),
		ws.WithPayloadMode(network.PayloadModeRawBody),
	)
	startLogicConnectScheduler()
	err := node.Start()
	if err != nil {
		panic(err)
	}

	sg := make(chan os.Signal)
	signal.Notify(sg, os.Interrupt, os.Kill)
	select {
	case sig := <-sg:
		logger.Info(fmt.Sprintf("game server is closing (signal: %v)", sig))
	case <-node.Running:
		logger.Info(fmt.Sprintf("game server is closing (signal: http)"))
	}
	// 执行所有关服逻辑
	node.Stop()
	stopOutboundDispatcher()
	closeAllBackendSessions()
}

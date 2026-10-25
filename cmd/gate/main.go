package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/forfun/gforgame/common/container/set"
	serverconfig "github.com/forfun/gforgame/config"
	"github.com/forfun/gforgame/gateway/contract"
	"github.com/forfun/gforgame/internal/gatewayadapter"
	"github.com/forfun/gforgame/internal/infra/net"
	"github.com/forfun/gforgame/internal/logger"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/protocol"
	"github.com/forfun/gforgame/network/ws"
	"github.com/forfun/gforgame/rpc"
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
	nodeId := flag.String("nodeId", "gate1000", "当前节点ID")
	flag.Parse()
	serverconfig.Init(*nodeId)

	onlinePlayerRegistry := net.NewOnlinePlayerRegistry()
	rpcFactory := protocol.NewMessageFactory()
	protos.RegisterGateRPCMessages(rpcFactory)
	gameRPCClient = rpc.NewClient(5*time.Second, rpcFactory)
	defer gameRPCClient.Close()
	for _, req := range []any{&protos.TransferGateToLogic{}, &protos.NotifyPlayerLogoutToGame{}, &protos.NotifyOnlinePlayerToGame{}} {
		if err := gameRPCClient.RegisterResponse(req, &protos.GateRPCResult{}); err != nil {
			panic(err)
		}
	}

	logicIoDispatcher = newLogicIoDispatcher(onlinePlayerRegistry)

	err := startLocalServerDiscoveryHeartbeat()
	if err != nil {
		panic(err)
	}
	startBackendSessionMonitor()
	startOutboundDispatcher()
	factory := protocol.NewMessageFactory()
	protos.RegisterMessages(factory)
	router := network.NewMessageRoute(factory)

	ioDispatcher := &MyMessageDispatch{onlinePlayerRegistry: onlinePlayerRegistry}
	ioDispatcher.AddHandler(&ClientRouter{router: router, onlinePlayerRegistry: onlinePlayerRegistry})
	codec := gateMsgCodec
	node := ws.NewServer(
		ws.WithAddress(serverconfig.GetCurrentNode().ServiceAddr),
		ws.WithRouter(router),
		ws.WithIoDispatch(ioDispatcher),
		ws.WithCodec(codec),
		ws.WithPayloadMode(protocol.PayloadModeRawBody),
	)
	err = node.Start()
	if err != nil {
		panic(err)
	}

	sg := make(chan os.Signal, 1)
	signal.Notify(sg, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sg)
	select {
	case sig := <-sg:
		logger.Info(fmt.Sprintf("game server is closing (signal: %v)", sig))
	case <-node.RunningChan():
		logger.Info("game server is closing (signal: http)")
	}
	// 执行所有关服逻辑
	close(serverDiscoveryStop)
	close(backendMonitorStop)
	node.Stop()
	gameRPCClient.Close()
	stopOutboundDispatcher()
	closeAllBackendSessions()
	logFilter.Close()
}

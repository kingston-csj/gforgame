package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/forfun/gforgame/codec/json"
	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/internal/bootstrap"
	"github.com/forfun/gforgame/internal/config"
	"github.com/forfun/gforgame/internal/infra/net"
	mysqldb "github.com/forfun/gforgame/internal/infra/persistence"
	"github.com/forfun/gforgame/internal/io"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/route"
	serverconfig "github.com/forfun/gforgame/internal/serverconfig"
	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/dispatch"
	"github.com/forfun/gforgame/network/protocol"
	"github.com/forfun/gforgame/network/ws"
	"github.com/forfun/gforgame/rpc"
)

//go:generate go run ../../internal/tools/gen all

// 使用静态路由，避免反射调用，提高性能。
// 代码生成内容见 route_dispatch_gen.go
// 部署前请先执行 go generate .
type generatedRouteInvoker func(msgHandler *network.Handler, playerID string, session network.Session, index int32, msg any) (any, error)

// 注册所有模块的路由
var generatedRouteDispatchers = map[int32]generatedRouteInvoker{}

type MyMessageDispatch struct {
	dispatch.BaseIoDispatch
	onlinePlayerRegistry *net.OnlinePlayerRegistry
}

func (m *MyMessageDispatch) OnSessionCreated(session network.Session) {
	if serverconfig.GetCurrentNode().UseGateMode {
		io.SetGateSession(session)
	}
	logger.Info(fmt.Sprintf("logic session created: %s", session.ToString()))
}

func (m *MyMessageDispatch) OnSessionClosed(session network.Session) {
	logger.Info(fmt.Sprintf("session closed: %s", session.ToString()))
	// 两个session相同，才可以移除在线
	// 否则，顶号会有问题
	unbindPlayer := false
	// 关闭session
	if v, ok := session.GetAttr("id"); ok {
		if sessionPlayerKey, ok := v.(string); ok && sessionPlayerKey != "" {
			if current := m.onlinePlayerRegistry.GetSessionByPlayerID(sessionPlayerKey); current == session {
				unbindPlayer = true
			}
		}
		m.onlinePlayerRegistry.RemovePlayerSession(session, unbindPlayer)
	}
}

func main() {
	nodeId := flag.String("nodeId", "game1001", "当前节点ID")
	flag.Parse()
	serverconfig.Init(*nodeId)

	logger.Info("game server is starting...")
	startTime := time.Now()

	factory := protocol.NewMessageFactory()
	protos.RegisterMessages(factory)
	router := network.NewMessageRoute(factory)

	codec := json.NewSerializer()

	mysqldb.InitMysql()
	bootstrap.InitMysqlDdl()
	config.InitConfig()
	s := bootstrap.InitServices()
	// resetTime := s.System.GetDailyReset().Data.(int64)
	// logger.Info(fmt.Sprintf("daily reset time: %d", resetTime))
	// 各自业务初始化
	bootstrap.InitBusiness(s)
	// 启动系统任务
	bootstrap.StartSchedulers(s)

	var modules = []any{
		route.NewCatalogRoute(s.Catalog, s.PlayerRepo),
		route.NewChatRoute(s.Chat, s.PlayerRepo),
		route.NewFriendRoute(s.Friend, s.PlayerRepo),
		route.NewGmRoute(s.Gm, s.PlayerRepo),
		route.NewHeroRoute(s.Hero, s.PlayerRepo),
		route.NewItemRoute(s.Item, s.PlayerRepo),
		route.NewMailRoute(s.Mail, s.PlayerRepo),
		route.NewMallRoute(s.Mall, s.PlayerRepo),
		route.NewMixtureRoute(s.Mixture, s.PlayerRepo),
		route.NewMonthCardRoute(s.MonthCard, s.PlayerRepo),
		route.NewPlayerRoute(s.Player, s.PlayerRepo),
		route.NewQuestRoute(s.Quest, s.PlayerRepo),
		route.NewRankRoute(s.Rank),
		route.NewRechargeRoute(),
		route.NewSignInRoute(s.SignIn, s.PlayerRepo),
	}
	if err := bootstrap.InitRouteModules(router, modules); err != nil {
		logger.Error("init route modules fail", err)
		panic(err)
	}

	// node := tcp.NewServer(
	// 	tcp.WithAddress(serverconfig.ServerConfig.ServerUrl),
	// 	tcp.WithRouter(router),
	// 	tcp.WithIoDispatch(ioDispatcher),
	// 	tcp.WithCodec(codec),
	// 	tcp.WithDispatchWorkers(8),
	// )

	io.SetOnlinePlayerRegistry(s.OnlinePlayerRegistry)
	io.SetMessageFactory(factory)

	ioDispatcher := &MyMessageDispatch{onlinePlayerRegistry: s.OnlinePlayerRegistry}
	gameTaskHandler := NewGameTaskHandler(router, s.ActorSystem)
	ioDispatcher.AddHandler(gameTaskHandler)

	rpcFactory := protocol.NewMessageFactory()
	protos.RegisterGateRPCMessages(rpcFactory)
	rpcNode := rpc.NewTCPServer(serverconfig.GetCurrentNode().RpcAddr, rpcFactory, codec,
		rpc.WithSessionHooks(ioDispatcher.OnSessionCreated, ioDispatcher.OnSessionClosed))
	if serverconfig.GetCurrentNode().RpcAddr != "" {
		if err := rpcNode.Register(newGateRPCService(route.NewGateRoute(s.Cross), gameTaskHandler, factory)); err != nil {
			panic(err)
		}
		if err := rpcNode.Start(); err != nil {
			panic(err)
		}
		defer rpcNode.Stop()
	}

	// listener := tcp.NewServer(
	// 	tcp.WithAddress(serverconfig.GetCurrentNode().ServiceAddr),
	// 	tcp.WithRouter(router),
	// 	tcp.WithIoDispatch(ioDispatcher),
	// 	tcp.WithCodec(codec),
	// )
	nodde := serverconfig.GetCurrentNode()
	if len(nodde.ServiceAddr) > 0 {
		listener := ws.NewServer(
			ws.WithAddress(nodde.ServiceAddr),
			ws.WithRouter(router),
			ws.WithIoDispatch(ioDispatcher),
			ws.WithCodec(codec),
			ws.WithUseGateway(nodde.UseGateMode),
			ws.WithWsPath("ws"),
		)
		err := listener.Start()
		if err != nil {
			logger.Error("game server start fail", err)
			panic(err)
		}
	}

	// 启动后台http服务
	go func() {
		NewHttpServer(router, s.PlayerRepo, s.Recharge, s.ProfileService)
	}()

	// 自启动，不会有额外消耗
	if len(serverconfig.GetCurrentNode().PprofAddr) > 0 {
		// pprof性能监控
		go func() {
			mux := NewHttpServeMux()
			// 监听并在 0.0.0.0:6060 上启动服务器

			http.ListenAndServe(serverconfig.GetCurrentNode().PprofAddr, mux)
		}()
	}

	endTime := time.Now()
	costSeconds := endTime.Sub(startTime).Seconds()
	logger.Info(fmt.Sprintf("game server is starting at %s, cost %.2fs", serverconfig.GetCurrentNode().ServiceAddr, costSeconds))

	sg := make(chan os.Signal, 1)
	signal.Notify(sg, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sg)
	select {
	case sig := <-sg:
		logger.Info(fmt.Sprintf("game server is closing (signal: %v)", sig))
		// case <-listener.RunningChan():
		// 	logger.Info("game server is closing (signal: http)")
	}
	// 执行所有关服逻辑
	rpcNode.Stop()
	// listener.Stop()

	// 框架模块
	s.DbService.Shutdown()
	// 业务模块

	logger.Info("game server is closed")
}

package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/forfun/gforgame/codec/json"
	"github.com/forfun/gforgame/common/logger"
	serverconfig "github.com/forfun/gforgame/config"
	"github.com/forfun/gforgame/internal/bootstrap"
	"github.com/forfun/gforgame/internal/context"
	"github.com/forfun/gforgame/internal/io"
	"github.com/forfun/gforgame/internal/route"
	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/ws"
)

//go:generate go run ../../internal/tools/gen all

// 使用静态路由，避免反射调用，提高性能。
// 代码生成内容见 route_dispatch_gen.go
// 部署前请先执行 go generate .
type generatedRouteInvoker func(msgHandler *network.Handler, playerID string, session *network.Session, index int32, msg any) (any, error)

// 注册所有模块的路由
var generatedRouteDispatchers = map[int32]generatedRouteInvoker{}

type MyMessageDispatch struct {
	network.BaseIoDispatch
}

func (m *MyMessageDispatch) OnSessionCreated(session *network.Session) {
	if serverconfig.ServerConfig.UseGateMode {
		io.SetGateSession(session)
	}
	logger.Info(fmt.Sprintf("logic session created: %s", session.ToString()))
}

func (m *MyMessageDispatch) OnSessionClosed(session *network.Session) {
	logger.Info(fmt.Sprintf("session closed: %s", session.ToString()))
	// 关闭session
	network.RemoveSession(session)
}

func main() {
	logger.Info(fmt.Sprintf("game server is starting..."))
	startTime := time.Now()

	router := network.NewMessageRoute()
	ioDispatcher := &MyMessageDispatch{}
	// 如果是网关模式，需要添加网关消息转换处理程序
	if serverconfig.ServerConfig.UseGateMode {
		ioDispatcher.AddHandler(NewGateTransformHandler())
	}
	ioDispatcher.AddHandler(&GameTaskHandler{router: router})
	// codec := protobuf.NewSerializer()
	codec := json.NewSerializer()

	// 自动建表
	bootstrap.InitMysqlDdl()
	// 加载配置数据
	bootstrap.InitConfig()
	// 注册所有服务
	s := bootstrap.InitServices()
	// 各自业务初始化
	bootstrap.InitBusiness(s)
	// 启动系统任务
	bootstrap.StartSchedulers()

	var modules = []any{
		route.NewCatalogRoute(s.Catalog, s.Player),
		route.NewChatRoute(s.Chat, s.Player),
		route.NewFriendRoute(s.Friend, s.Player),
		route.NewGmRoute(s.Gm, s.Player),
		route.NewHeroRoute(s.Hero, s.Player),
		route.NewItemRoute(s.Item, s.Player),
		route.NewMailRoute(s.Mail, s.Player),
		route.NewMallRoute(s.Mall, s.Player),
		route.NewMixtureRoute(s.Mixture, s.Player),
		route.NewMonthCardRoute(s.MonthCard, s.Player),
		route.NewPlayerRoute(s.Player),
		route.NewQuestRoute(s.Quest, s.Player),
		route.NewRankRoute(s.Rank),
		route.NewRechargeRoute(),
		route.NewSignInRoute(s.SignIn, s.Player),
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
	// context.GameServer = node
	node := ws.NewServer(
		ws.WithAddress(serverconfig.ServerConfig.ServerUrl),
		ws.WithRouter(router),
		ws.WithIoDispatch(ioDispatcher),
		ws.WithCodec(codec),
		ws.WithUseGateway(serverconfig.ServerConfig.UseGateMode),
		ws.WithWsPath("ws"),
	)
	context.GameServer = node

	err := node.Start()
	if err != nil {
		logger.Error("game server start fail", err)
		panic(err)
	}

	// 启动后台http服务
	go func() {
		context.HttpServer = NewHttpServer()
	}()

	// 自启动，不会有额外消耗
	if len(serverconfig.ServerConfig.PprofAddr) > 0 {
		// pprof性能监控
		go func() {
			mux := NewHttpServeMux()
			// 监听并在 0.0.0.0:6060 上启动服务器
			http.ListenAndServe(serverconfig.ServerConfig.PprofAddr, mux)
		}()
	}

	endTime := time.Now()
	costSeconds := endTime.Sub(startTime).Seconds()
	logger.Info(fmt.Sprintf("game server is starting at %s, cost %.2fs", serverconfig.ServerConfig.ServerUrl, costSeconds))

	sg := make(chan os.Signal, 1)
	signal.Notify(sg, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sg)
	select {
	case sig := <-sg:
		logger.Info(fmt.Sprintf("game server is closing (signal: %v)", sig))
	case <-node.RunningChan():
		logger.Info(fmt.Sprintf("game server is closing (signal: http)"))
	}
	// 执行所有关服逻辑
	node.Stop()
	// 框架模块
	context.DbService.Shutdown()

	logger.Info(fmt.Sprintf("game server is closed"))
}

package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/forfun/gforgame/codec/json"
	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/common/util/jsonutil"
	serverconfig "github.com/forfun/gforgame/config"
	"github.com/forfun/gforgame/internal/bootstrap"
	"github.com/forfun/gforgame/internal/constants"
	"github.com/forfun/gforgame/internal/context"
	"github.com/forfun/gforgame/internal/http"
	"github.com/forfun/gforgame/internal/route"
	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/protocol"
	"github.com/forfun/gforgame/network/ws"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//go:generate go generate ../../internal/bootstrap

type GameTaskHandler struct {
	router *network.MessageRoute
}

// 使用静态路由，避免反射调用，提高性能。
// 代码生成内容见 route_dispatch_gen.go
// 部署前请先执行 go generate .
type generatedRouteInvoker func(msgHandler *network.Handler, session *network.Session, index int32, msg any) (any, error)

var generatedRouteDispatchers = map[int32]generatedRouteInvoker{}

func (g *GameTaskHandler) MessageReceived(session *network.Session, frame *protocol.RequestDataFrame) bool {
	defer func() {
		if r := recover(); r != nil {
			logger.ErrorNoStack(fmt.Errorf("panic recovered: %v", r))
		}
	}()
	msgName, _ := network.GetMsgName(frame.Header.Cmd)
	jsonStr, err := jsonutil.StructToJSON(frame.Msg)
	if err == nil {
		if strings.Index(msgName, "HeartBeat") == -1 {
			id, ok := session.GetAttr("id")
			if !ok {
				id = "anonymous"
			}
			// fmt.Println("接收消息: cmd:", frame.Header.Cmd, " name:", msgName, " 内容:", jsonStr)
			logger.Info(fmt.Sprintf("[%s] 接收消息: cmd:%d, name:%s, 内容:%s", id, frame.Header.Cmd, msgName, jsonStr))
		}
	}

	msgHandler, _ := g.router.GetHandler(frame.Header.Cmd)
	if msgHandler == nil {
		logger.ErrorNoStack(fmt.Errorf("msgHandler is nil: %v", frame.Header.Cmd))
		return false
	}

	resp, handled, dispatchErr, panicErr := callGeneratedRouteHandlerSafely(frame.Header.Cmd, msgHandler, session, frame.Header.Index, frame.Msg)
	if handled {
		if panicErr != nil {
			logger.Error(fmt.Sprintf("generated route handler panic: cmd=%d method=%s", frame.Header.Cmd, msgHandler.Method.Name), panicErr)
			if errorResp, ok := buildErrorResponse(msgHandler, constants.I18N_COMMON_INTERNAL_ERROR); ok {
				session.Send(errorResp, frame.Header.Index)
			}
			return false
		}
		if dispatchErr != nil {
			// 静态分发失败时回退反射调用，保证兼容性。
			logger.ErrorNoStack(fmt.Errorf("generated dispatch failed, fallback to reflect: cmd=%d err=%v", frame.Header.Cmd, dispatchErr))
		} else {
			if resp != nil {
				session.Send(resp, frame.Header.Index)
			}
			return true
		}
	}

	var args []reflect.Value
	if msgHandler.Indindexed {
		args = []reflect.Value{msgHandler.Receiver, reflect.ValueOf(session), reflect.ValueOf(frame.Header.Index), reflect.ValueOf(frame.Msg)}
	} else {
		args = []reflect.Value{msgHandler.Receiver, reflect.ValueOf(session), reflect.ValueOf(frame.Msg)}
	}

	// 反射调用路由处理器，并捕获处理器内部 panic
	values, panicErr := callRouteHandlerSafely(msgHandler, args)
	if panicErr != nil {
		logger.Error(fmt.Sprintf("route handler panic: cmd=%d method=%s", frame.Header.Cmd, msgHandler.Method.Name), panicErr)
		if resp, ok := buildErrorResponse(msgHandler, constants.I18N_COMMON_INTERNAL_ERROR); ok {
			session.Send(resp, frame.Header.Index)
		}
		return false
	}

	if len(values) > 0 {
		err := session.Send(values[0].Interface(), frame.Header.Index)
		if err != nil {
			logger.Error("session.Send: %v", fmt.Errorf("session.Send: %v", err))
			return false
		}
	}
	return true
}

func callGeneratedRouteHandlerSafely(cmd int32, msgHandler *network.Handler, session *network.Session, index int32, msg any) (resp any, handled bool, dispatchErr error, panicErr error) {
	invoker, ok := getGeneratedRouteInvoker(cmd)
	if !ok {
		return nil, false, nil, nil
	}
	handled = true
	defer func() {
		if r := recover(); r != nil {
			panicErr = logger.PanicToError(r)
		}
	}()
	resp, dispatchErr = invoker(msgHandler, session, index, msg)
	return resp, handled, dispatchErr, nil
}

func getGeneratedRouteInvoker(cmd int32) (generatedRouteInvoker, bool) {
	if generatedRouteDispatchers == nil {
		return nil, false
	}
	invoker, ok := generatedRouteDispatchers[cmd]
	return invoker, ok
}

func callRouteHandlerSafely(msgHandler *network.Handler, args []reflect.Value) (values []reflect.Value, panicErr error) {
	defer func() {
		if r := recover(); r != nil {
			panicErr = logger.PanicToError(r)
		}
	}()
	values = msgHandler.Method.Func.Call(args)
	return values, nil
}

func buildErrorResponse(msgHandler *network.Handler, code int32) (any, bool) {
	mt := msgHandler.Method.Type
	if mt.NumOut() == 0 {
		return nil, false
	}
	outType := mt.Out(0)
	if outType.Kind() != reflect.Ptr || outType.Elem().Kind() != reflect.Struct {
		return nil, false
	}
	resp := reflect.New(outType.Elem())
	codeField := resp.Elem().FieldByName("Code")
	if !codeField.IsValid() || !codeField.CanSet() || codeField.Kind() != reflect.Int32 {
		return nil, false
	}
	codeField.SetInt(int64(code))
	return resp.Interface(), true
}

func NewHttpServer() *gin.Engine {
	router := gin.Default()
	// 关闭游戏服务器进
	router.POST("/admin/stop", func(c *gin.Context) {
		http.StopServer(c)
	})
	// 配置 CORS 中间
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有源，生产环境应指定具体域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	err := router.Run(serverconfig.ServerConfig.HttpUrl)
	if err != nil {
		panic(err)
	}

	return router
}

type MyMessageDispatch struct {
	network.BaseIoDispatch
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
	ioDispatcher.AddHandler(&GameTaskHandler{router: router})
	// codec := protobuf.NewSerializer()
	codec := json.NewSerializer()

	// 自动建表
	bootstrap.InitMysqlDdl()
	// 开发环境，导出所有客户端协议
	bootstrap.DevOnlyExportProtocols()
	// 加载配置数据
	bootstrap.InitConfig()
	// 预热服务并完成跨模块注册（避免首次请求时懒初始化副作用）
	bootstrap.InitServices()
	// 各自业务初始化
	bootstrap.InitBusiness()
	// 启动系统任务
	bootstrap.StartSchedulers()

	// 在这里，添加你的模块消息路由
	modules := []network.Module{
		route.NewPlayerRoute(),
		route.NewHeroRoute(),
		route.NewSceneRoute(),
		route.NewQuestRoute(),
		route.NewGmRoute(),
		route.NewSignInRoute(),
		route.NewItemRoute(),
		route.NewMallRoute(),
		route.NewMonthCardRoute(),
		route.NewMailRoute(),
		route.NewRankRoute(),
		route.NewRechargeRoute(),
		route.NewMixtureRoute(),
		route.NewCatalogRoute(),
		route.NewChatRoute(),
		route.NewFriendRoute(),
	}

	node := ws.NewServer(
		ws.WithAddress(serverconfig.ServerConfig.ServerUrl),
		ws.WithRouter(router),
		ws.WithIoDispatch(ioDispatcher),
		ws.WithCodec(codec),
		ws.WithModules(modules...),
		ws.WithWsPath("ws"),
	)
	context.WsServer = node

	err := node.Start()
	if err != nil {
		panic(err)
	}

	// 启动rpc服务
	// if len(serverconfig.ServerConfig.RpcServerUrl) > 0 {
	// 	go func() {
	// 		NewRpcServer(serverconfig.ServerConfig.RpcServerUrl)
	// 	}()
	// }

	// 启动后台http服务
	// go func() {
	// 	context.HttpServer = NewHttpServer()
	// }()

	// pprof性能监控
	// go func() {
	// 	mux := NewHttpServeMux()
	// 	// 监听并在 0.0.0.0:6060 上启动服务器
	// 	http.ListenAndServe(serverconfig.ServerConfig.PprofAddr, mux)
	// }()

	endTime := time.Now()
	cost := endTime.Sub(startTime)
	logger.Info(fmt.Sprintf("game server is starting at " + serverconfig.ServerConfig.ServerUrl + ", cost %.2f seconds", cost.Seconds()))

	sg := make(chan os.Signal, 1)
	signal.Notify(sg, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sg)
	select {
	case sig := <-sg:
		logger.Info(fmt.Sprintf("game server is closing (signal: %v)", sig))
	case <-node.Running:
		logger.Info(fmt.Sprintf("game server is closing (signal: http)"))
	}
	// 执行所有关服逻辑
	node.Stop()
	context.DbService.Shutdown()
	logger.Info(fmt.Sprintf("game server is closed"))
}

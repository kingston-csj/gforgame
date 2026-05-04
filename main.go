package main

import (
	"fmt"
	"io/github/gforgame/codec/json"
	"io/github/gforgame/common/logger"
	"io/github/gforgame/common/util/jsonutil"
	serverconfig "io/github/gforgame/config"
	"io/github/gforgame/examples/bootstrap"
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/http"
	mysqldb "io/github/gforgame/examples/infra/persistence"
	"io/github/gforgame/examples/route"
	"io/github/gforgame/examples/system"
	protocolValidator "io/github/gforgame/examples/validator"
	"io/github/gforgame/network"
	"io/github/gforgame/network/protocol"
	"io/github/gforgame/network/ws"
	protocolexporter "io/github/gforgame/tools/protocol"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type GameTaskHandler struct {
	router *network.MessageRoute
}

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

	// 验证协议的参数，如果校验失败，直接返回错误码
	if msgHandler.NeedValidate {
		validationErrors := protocolValidator.ValidateStruct(frame.Msg)
		if len(validationErrors) > 0 {
			errMsg := protocolValidator.FormatValidationErrors(validationErrors)
			if errMsg != "" {
				// logger.Info(fmt.Sprintf("validation failed for cmd=%d: %s", frame.Header.Cmd, errMsg))
				if resp, ok := buildErrorResponse(msgHandler, constants.I18N_COMMON_PROTOCOL_VALIDATION_FAILED); ok {
					session.Send(resp, frame.Header.Index)
				}
				return false
			}
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
			if err := session.Send(resp, frame.Header.Index); err != nil {
				// logger.Error(fmt.Errorf("session.Send error response failed: %v", err))
			}
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
	// 关闭游戏服务器进程
	router.POST("/admin/stop", func(c *gin.Context) {
		http.StopServer(c)
	})
	// 配置 CORS 中间件
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
	logger.Info(fmt.Sprintf("session closed: %v", session))
	// 关闭session
	network.RemoveSession(session)
}

func main() {
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

	// node := tcp.NewServer(
	// 	tcp.WithAddress(serverconfig.ServerConfig.ServerUrl),
	// 	tcp.WithRouter(router),
	// 	tcp.WithIoDispatch(ioDispatcher),
	// 	tcp.WithCodec(codec),
	// 	tcp.WithModules(modules...),
	// )

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

	// 启动rpc服务器
	// if len(serverconfig.ServerConfig.RpcServerUrl) > 0 {
	// 	go func() {
	// 		NewRpcServer(serverconfig.ServerConfig.RpcServerUrl)
	// 	}()
	// }

	// 启动后台http服务器
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
	logger.Info("game server is starting at " + serverconfig.ServerConfig.ServerUrl + ", cost " + endTime.Sub(startTime).String())

	sg := make(chan os.Signal)
	signal.Notify(sg, os.Interrupt, syscall.SIGTERM)
	select {
	case sig := <-sg:
		logger.Info(fmt.Sprintf("game server is closing (signal: %v)", sig))
	case <-node.Running:
		logger.Info(fmt.Sprintf("game server is closing (signal: http)"))
	}
	// 执行所有关服逻辑
	node.Stop()
}

// 自动建表
func autoCrreateDatabase() {
	// 玩家表
	err := mysqldb.Db.AutoMigrate(&playerdomain.Player{})
	if err != nil {
		panic(err)
	}
	// 好友表
	err = mysqldb.Db.AutoMigrate(&playerdomain.Friend{})
	if err != nil {
		panic(err)
	}
	// 场景表
	err = mysqldb.Db.AutoMigrate(&playerdomain.Scene{})
	if err != nil {
		panic(err)
	}
	// 系统参数表
	err = mysqldb.Db.AutoMigrate(&system.SystemParameterEnt{})
	if err != nil {
		panic(err)
	}
}

func TryExportProtocols() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}
	// 开发环境，导出所有客户端协议
	if env == "dev" {
		// generator := protocolexporter.NewTypeScriptGenerator(
		// 	"protos",
		// 	"tools\\protocol\\output\\typescript\\",
		// 	"tools\\protocol\\templates\\tstemplate.tpl",
		// )
		generator := protocolexporter.NewCSharpGenerator(
			"examples\\protos",
			"tools\\protocol\\output\\csharp\\",
			"tools\\protocol\\templates\\csharptemplate.tpl",
		)

		error := generator.Generate(network.GetMsgName2IdMapper())
		if error != nil {
			panic(error)
		}
		err2 := generator.BaseGenerator.GenerateRegisterFromTags("examples\\protos", "examples\\protos\\register_gen.go", nil)
		if err2 != nil {
			panic(err2)
		}
	}
}

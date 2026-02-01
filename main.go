package main

import (
	"fmt"
	"io/github/gforgame/codec/json"
	serverconfig "io/github/gforgame/config"
	mysqldb "io/github/gforgame/db"
	"io/github/gforgame/examples/activity"
	dataconfig "io/github/gforgame/examples/config"
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/http"
	"io/github/gforgame/examples/route"
	"io/github/gforgame/examples/service/player"
	"io/github/gforgame/examples/system"
	"io/github/gforgame/logger"
	"io/github/gforgame/network"
	"io/github/gforgame/network/protocol"
	"io/github/gforgame/network/tcp"
	protocolexporter "io/github/gforgame/tools/protocol"
	"io/github/gforgame/util/jsonutil"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"reflect"
	"strings"
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
			logger.Error(fmt.Errorf("panic recovered: %v", r))
		}
	}()
	msgName, _ := network.GetMsgName(frame.Header.Cmd)
	jsonStr, err := jsonutil.StructToJSON(frame.Msg)
	if err == nil {
		if strings.Index(msgName, "HeartBeat") == -1 {
			fmt.Println("接收消息: cmd:", frame.Header.Cmd, " name:", msgName, " 内容：", jsonStr)
		}
	}
	

	msgHandler, _ := g.router.GetHandler(frame.Header.Cmd)
	if msgHandler == nil {
		logger.Error3(fmt.Errorf("msgHandler is nil: %v", frame.Header.Cmd))
		return false
	}
	var args []reflect.Value
	if msgHandler.Indindexed {
		args = []reflect.Value{msgHandler.Receiver, reflect.ValueOf(session), reflect.ValueOf(frame.Header.Index), reflect.ValueOf(frame.Msg)}
	} else {
		args = []reflect.Value{msgHandler.Receiver, reflect.ValueOf(session), reflect.ValueOf(frame.Msg)}
	}

	// 反射
	values := msgHandler.Method.Func.Call(args)
	if len(values) > 0 {
		err := session.Send(values[0].Interface(), frame.Header.Index)
		if err != nil {
			logger.Error(fmt.Errorf("session.Send: %v", err))
			return false
		}
	}
	return true
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
	autoCrreateDatabase()
	// 开发环境，导出所有客户端协议
	TryExportProtocols()

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

	node := tcp.NewServer(
		tcp.WithAddress(serverconfig.ServerConfig.ServerUrl),
		tcp.WithRouter(router),
		tcp.WithIoDispatch(ioDispatcher),
		tcp.WithCodec(codec),
		tcp.WithModules(modules...),
	)
	context.TcpServer = node

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

	dataconfig.GetDataManager()

	// itemData := config.QueryById[configdomain.PropData](10000001) 
	// if itemData == nil {
	// 	panic("item data not found")
	// }

	system.StartSystemTask()

	endTime := time.Now()
	logger.Info("game server is starting at " + serverconfig.ServerConfig.ServerUrl + ", cost " + endTime.Sub(startTime).String())

	// rank.GetRankService().QueryRank(rank.PlayerLevelRank, 0, 10)

	// fight.GetFightService().Test()
 
	// 各自业务初始化
	player.GetPlayerService().LoadPlayerProfile()

	activity.GetActivityService().ScheduleAllActivity()

	// p := player.GetPlayerService().GetPlayer("111")
	// logger.Info(p.Name)

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
			"protos",
			"tools\\protocol\\output\\csharp\\",
			"tools\\protocol\\templates\\csharptemplate.tpl",
		)
		
		error := generator.Generate(network.GetMsgName2IdMapper())	
		if error != nil {
			panic(error)
		}
		err2 := generator.BaseGenerator.GenerateRegisterFromTags("protos", "protos\\register_gen.go", nil)
		if err2 != nil {
			panic(err2)
		}
	}
}

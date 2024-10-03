package main

import (
	"fmt"
	"io/github/gforgame/codec/protobuf"
	"io/github/gforgame/config"
	"io/github/gforgame/examples/chat"
	"io/github/gforgame/examples/player"
	"io/github/gforgame/log"
	"io/github/gforgame/network"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
)

type GameTaskHandler struct {
}

var (
	node   *network.Node
	router *gin.Engine
)

func (g *GameTaskHandler) MessageReceived(session *network.Session, frame *network.RequestDataFrame) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Error(r.(error))
		}
	}()
	msgHandler, _ := network.GetHandler(frame.Header.Cmd)
	args := []reflect.Value{msgHandler.Receiver, reflect.ValueOf(session), reflect.ValueOf(frame.Msg)}
	// 反射
	values := msgHandler.Method.Func.Call(args)
	if len(values) > 0 {
		err := session.Send(values[0].Interface())
		if err != nil {
			log.Error(fmt.Errorf("session.Send: %v", err))
			return false
		}
	}
	return true
}

func NewHttpServer() *gin.Engine {
	router := gin.Default()
	// 关闭游戏服务器进程
	router.GET("/admin/stop", func(c *gin.Context) {
		node.Running <- true
	})
	return router
}

func StartHttpServer(router *gin.Engine) {
	err := router.Run(config.ServerConfig.HttpUrl)
	if err != nil {
		panic(err)
	}
}

func main() {

	startTime := time.Now()
	network.RegisterModule(chat.NewRoomService())
	network.RegisterModule(player.NewPlayerService())

	modules := network.ListModules()
	for _, c := range modules {
		c.Init()
		err := network.RegisterMessageHandlers(c)
		if err != nil {
			panic(err)
		}
	}

	ioDispatcher := &network.BaseIoDispatch{}
	ioDispatcher.AddHandler(&GameTaskHandler{})

	// 设置服务器监听的地址和端口
	node = &network.Node{
		Running: make(chan bool),
	}
	codec := protobuf.NewSerializer()
	//codec := &json.NewSerializer()
	err := node.Startup(network.WithAddress(config.ServerConfig.ServerUrl), network.WithIoDispatch(ioDispatcher), network.WithCodec(codec))
	//err := node.Startup(network.WithAddress(config.ServerConfig.ServerUrl), network.WithIoDispatch(ioDispatcher), network.WithCodec(codec), network.WithWebsocket())
	if err != nil {
		panic(err)
	}

	// 启动rpc服务器
	if len(config.ServerConfig.RpcServerUrl) > 0 {
		go func() {
			NewRpcServer(config.ServerConfig.RpcServerUrl)
		}()
	}

	// 启动后台http服务器
	// router = NewHttpServer()
	// go func() {
	// 	StartHttpServer(router)
	// }()

	// pprof性能监控
	// go func() {
	// 	mux := NewHttpServeMux()
	// 	// 监听并在 0.0.0.0:6060 上启动服务器
	// 	http.ListenAndServe(config.ServerConfig.PprofAddr, mux)
	// }()

	endTime := time.Now()
	log.Info("game server is starting, cost " + endTime.Sub(startTime).String())

	sg := make(chan os.Signal)
	signal.Notify(sg, os.Interrupt, os.Kill)
	select {
	case sig := <-sg:
		log.Info(fmt.Sprintf("game server is closing (signal: %v)", sig))
		showDown(modules)
	case <-node.Running:
		log.Info(fmt.Sprintf("game server is closing (signal: http)"))
		showDown(modules)
	}
}

func showDown(modules []network.Module) {
	for _, c := range modules {
		c.Shutdown()
	}
	log.Info("game server has closed...")
}

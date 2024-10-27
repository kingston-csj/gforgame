package main

import (
	"fmt"
	"io/github/gforgame/codec/protobuf"
	"io/github/gforgame/config"
	"io/github/gforgame/examples/chat"
	"io/github/gforgame/examples/player"
	"io/github/gforgame/logger"
	"io/github/gforgame/network"
	"io/github/gforgame/network/protocol"
	"io/github/gforgame/network/tcp"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	node   network.Server
	router *gin.Engine
)

type GameTaskHandler struct {
	router *network.MessageRoute
}

func (g *GameTaskHandler) MessageReceived(session *network.Session, frame *protocol.RequestDataFrame) bool {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(r.(error))
		}
	}()
	msgHandler, _ := g.router.GetHandler(frame.Header.Cmd)
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

// func NewHttpServer() *gin.Engine {
// 	router := gin.Default()
// 	// 关闭游戏服务器进程
// 	router.GET("/admin/stop", func(c *gin.Context) {
// 		node.Running <- true
// 	})
// 	return router
// }

// func StartHttpServer(router *gin.Engine) {
// 	err := router.Run(config.ServerConfig.HttpUrl)
// 	if err != nil {
// 		panic(err)
// 	}
// }

func main() {

	startTime := time.Now()

	router := &network.MessageRoute{Handlers: make(map[int]*network.Handler)}
	ioDispatcher := &network.BaseIoDispatch{}
	ioDispatcher.AddHandler(&GameTaskHandler{router: router})
	codec := protobuf.NewSerializer()
	// codec := json.NewSerializer()

	node := tcp.NewServer(tcp.WithAddress(config.ServerConfig.ServerUrl), tcp.WithRouter(router),
		tcp.WithIoDispatch(ioDispatcher), tcp.WithCodec(codec), tcp.WithModules(chat.NewRoomService(), player.NewPlayerService()))

	// node := ws.NewServer(ws.WithAddress(config.ServerConfig.ServerUrl), ws.WithRouter(router),
	// 	ws.WithIoDispatch(ioDispatcher), ws.WithCodec(codec), ws.WithModules(chat.NewRoomService(), player.NewPlayerService()))

	err := node.Start()
	if err != nil {
		panic(err)
	}

	// 启动rpc服务器
	// if len(config.ServerConfig.RpcServerUrl) > 0 {
	// 	go func() {
	// 		NewRpcServer(config.ServerConfig.RpcServerUrl)
	// 	}()
	// }

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
	logger.Info("game server is starting, cost " + endTime.Sub(startTime).String())

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

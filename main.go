package main

import (
	"fmt"
	"io/github/gforgame/codec/protobuf"
	"io/github/gforgame/config"
	"io/github/gforgame/examples/chat"
	"io/github/gforgame/examples/player"
	"io/github/gforgame/log"
	"io/github/gforgame/network"
	"os"
	"os/signal"
	"reflect"
	"time"
)

type GameTaskHandler struct {
}

func (g *GameTaskHandler) MessageReceived(session *network.Session, frame network.RequestDataFrame) bool {
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
		session.Send(values[0].Interface())
	}

	return true
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
	node := &network.Node{}
	codec := &protobuf.ProtobufCodec{}
	//codec := &json.JsonCodec{}
	err := node.Startup(network.WithAddress(config.ServerConfig.ServerUrl), network.WithIoDispatch(ioDispatcher), network.WithCodec(codec))
	//err := node.Startup(network.WithAddress(config.ServerConfig.ServerUrl), network.WithIoDispatch(ioDispatcher), network.WithCodec(codec), network.WithWebsocket())
	if err != nil {
		panic(err)
	}
	endTime := time.Now()
	log.Info("game server is starting, cost " + endTime.Sub(startTime).String())

	sg := make(chan os.Signal)
	signal.Notify(sg, os.Interrupt, os.Kill)
	sig := <-sg
	log.Info(fmt.Sprintf("game server is closing (signal: %v)", sig))

	for _, c := range modules {
		c.Shutdown()
	}
	log.Info("game server has closed...")
}

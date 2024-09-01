package main

import (
	"fmt"
	"io/github/gforgame/codec"
	"io/github/gforgame/examples/chat"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
	"log"
	"os"
	"os/signal"
	"reflect"
)

type GameTaskHandler struct {
}

func (g *GameTaskHandler) MessageReceived(session *network.Session, frame network.RequestDataFrame) bool {
	msgHandler, _ := network.GetHandler(frame.Header.Cmd)

	args := []reflect.Value{msgHandler.Receiver, reflect.ValueOf(session), reflect.ValueOf(frame.Msg)}
	// 反射
	msgHandler.Method.Func.Call(args)
	return true
}

func main() {
	network.RegisterModule(chat.NewRoomService())

	modules := network.ListModules()
	for _, c := range modules {
		c.Init()
		err := network.RegisterMessageHandlers(c)
		if err != nil {
			panic(err)
		}
	}

	network.RegisterMessage(1001, &protos.ReqJoinRoom{})
	network.RegisterMessage(1002, &protos.ReqChat{})

	fmt.Println(network.GetMessageCmd(&protos.ReqChat{}))
	fmt.Println(network.GetMessageType(1001))

	ioDispatcher := &network.BaseIoDispatch{}
	ioDispatcher.AddHandler(&GameTaskHandler{})

	// 设置服务器监听的地址和端口
	node := &network.Node{}
	err := node.Startup(network.WithAddress("127.0.0.1:9090"), network.WithIoDispatch(ioDispatcher), network.WithCodec(&codec.JsonCodec{}))
	if err != nil {
		panic(err)
	}

	sg := make(chan os.Signal)
	signal.Notify(sg, os.Interrupt, os.Kill)
	sig := <-sg
	log.Println("game server is closing(signal: %v)", sig)

	for _, c := range modules {
		c.Shutdown()
	}

	log.Println("game server has closed...")

}

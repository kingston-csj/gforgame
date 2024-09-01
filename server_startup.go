package main

import (
	"fmt"
	"io/github/gforgame/codec"
	"io/github/gforgame/examples/chat"
	"io/github/gforgame/log"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
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
	log.Printf("----test---")

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

	log.Info(log.APPLICATION, "key", "value")
	log.Printf("game server is starting ...")
	err := node.Startup(network.WithAddress("127.0.0.1:9090"), network.WithIoDispatch(ioDispatcher), network.WithCodec(&codec.JsonCodec{}))
	if err != nil {
		panic(err)
	}

	sg := make(chan os.Signal)
	signal.Notify(sg, os.Interrupt, os.Kill)
	sig := <-sg
	log.Info(log.APPLICATION, "game server is closing(signal: %v)", sig)

	for _, c := range modules {
		c.Shutdown()
	}

	log.Printf("game server has closed...")

}

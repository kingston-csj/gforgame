package main

import (
	"fmt"
	"io/github/gforgame/codec/json"
	"io/github/gforgame/db"
	"io/github/gforgame/examples/cross"
	"io/github/gforgame/examples/player"
	"io/github/gforgame/logger"
	"io/github/gforgame/network"
	"io/github/gforgame/network/client"
	"io/github/gforgame/network/protocol"
	"io/github/gforgame/protos"
	"os"
	"os/signal"
	"reflect"
)

// 实现 RequestCallback 接口的匿名对象
type commonCallback struct{}

func (rc *commonCallback) OnSuccess(result any) {
	fmt.Println("OnSuccess: ", result)
}

func (rc *commonCallback) OnError(err error) {
	fmt.Println("OnError: ", err)
}

type GameTaskHandler struct {
	Router network.MessageRoute
}

func (g *GameTaskHandler) MessageReceived(session *network.Session, frame *protocol.RequestDataFrame) bool {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(r.(error))
		}
	}()

	// 如果有消息流水号，优先走消息回调
	if frame.Header.Index > 0 {
		client.CallBackManager.FillCallBack(frame.Header.Index, frame.Msg)
	} else {
		// 通过方法签名，能否找到消息路由
		msgHandler, err := g.Router.GetHandler(frame.Header.Cmd)
		if err == nil {
			args := []reflect.Value{msgHandler.Receiver, reflect.ValueOf(session), reflect.ValueOf(frame.Msg)}
			// 反射
			values := msgHandler.Method.Func.Call(args)
			if len(values) > 0 {
				err := session.Send(values[0].Interface(), frame.Header.Index)
				if err != nil {
					logger.Error(fmt.Errorf("session.Send: %v", err))
					return false
				}
			}

		} else {
			// 没有找到，直接打印
			fmt.Println("客户端收到消息：(", frame.Header.Cmd, ")", frame.Msg)
		}
	}

	return true
}

func main() {
	p := &player.Player{db.BaseEntity{Id: "123456"}, "gforgame", 999}
	cross.PlayerLoginRemote(p, cross.Island)

	network.RegisterMessage(protos.CmdChatReqJoin, &protos.ReqJoinRoom{})
	network.RegisterMessage(protos.CmdChatReqChat, &protos.ReqChat{})
	network.RegisterMessage(protos.CmdPlayerReqLogin, &protos.ReqPlayerLogin{})
	network.RegisterMessage(protos.CmdPlayerResLogin, &protos.ResPlayerLogin{})
	network.RegisterMessage(protos.CmdPlayerReqCreate, &protos.ReqPlayerCreate{})
	network.RegisterMessage(protos.CmdPlayerResCreate, &protos.ResPlayerCreate{})

	// 服务器地址和端口
	address := "127.0.0.1:9090"
	// msgCodec := protobuf.NewSerializer()
	msgCodec := json.NewSerializer()

	ioDispatcher := &network.BaseIoDispatch{}
	ioDispatcher.AddHandler(&GameTaskHandler{})
	tcpClient := &client.TcpSocketClient{RemoteAddress: address, MsgCodec: msgCodec}
	session, err := tcpClient.OpenSession()
	if err != nil {
		panic(err)
	}

	go func() {
		for frame := range session.DataReceived {
			ioDispatcher.OnMessageReceived(session, frame)
		}
	}()

	// session.Send(&protos.ReqPlayerLogin{Id: "1001"}, 0)
	// client.Callback(session, &protos.ReqPlayerLogin{Id: "1001"}, &commonCallback{})
	r, err := client.Request(session, &protos.ReqPlayerLogin{Id: "1001"})
	if err != nil {
		fmt.Println(err)
	}
	resPlayerLogin := r.(*protos.ResPlayerLogin)
	fmt.Println("客户端收到消息：(", resPlayerLogin, ")")

	// session.Send(&protos.ReqJoinRoom{RoomId: 123, PlayerId: 123}, 0)

	sg := make(chan os.Signal)
	signal.Notify(sg, os.Interrupt, os.Kill)
	_ = <-sg
}

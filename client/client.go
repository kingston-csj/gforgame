package main

import (
	"fmt"
	"io/github/gforgame/codec/protobuf"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
	"log"
	"net"
	"os"
	"os/signal"
	"reflect"
)

func main() {

	network.RegisterMessage(protos.CmdChatReqJoin, &protos.ReqJoinRoom{})
	network.RegisterMessage(protos.CmdChatReqChat, &protos.ReqChat{})
	network.RegisterMessage(protos.CmdPlayerReqLogin, &protos.ReqPlayerLogin{})
	network.RegisterMessage(protos.CmdPlayerResLogin, &protos.ResPlayerLogin{})
	network.RegisterMessage(protos.CmdPlayerReqCreate, &protos.ReqPlayerCreate{})
	network.RegisterMessage(protos.CmdPlayerResCreate, &protos.ResPlayerCreate{})

	// 服务器地址和端口
	address := "127.0.0.1:9090"
	// 连接服务器
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal("连接服务器失败:", err)
	}
	defer conn.Close()
	fmt.Println("已连接到服务器:", address)

	//session := network.NewSession(&conn, &json.JsonCodec{})
	msgCodec := &protobuf.ProtobufCodec{}
	session := network.NewSession(&conn, msgCodec)
	go session.Write()

	go func() {
		// read loop
		buf := make([]byte, 2048)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				log.Printf(fmt.Sprintf("Read message error: %s, session will be closed immediately", err.Error()))
				return
			}
			packets, err := session.ProtocolCodec.Decode(buf[:n])
			if err != nil {
				log.Println(err.Error())
				return
			}
			// process packets decoded
			for _, p := range packets {
				typ, _ := network.GetMessageType(p.Header.Cmd)
				msg := reflect.New(typ.Elem()).Interface()
				msgCodec.Decode(p.Data, msg)
				fmt.Println("客户端收到消息：(", p.Header.Cmd, ")", msg)
			}
		}
	}()
	//req := &protos.ReqJoinRoom{RoomId: 123, PlayerId: 123}
	req := &protos.ReqPlayerLogin{Id: 1001}
	session.Send(req)

	sg := make(chan os.Signal)
	signal.Notify(sg, os.Interrupt, os.Kill)
	_ = <-sg
}

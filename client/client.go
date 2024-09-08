package main

import (
	"fmt"
	"io/github/gforgame/codec/protobuf"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
	"log"
	"net"
)

func main() {
	// 服务器地址和端口
	address := "127.0.0.1:9090"

	network.RegisterMessage(protos.CmdChatReqJoin, &protos.ReqJoinRoom{})
	network.RegisterMessage(protos.CmdChatReqChat, &protos.ReqChat{})
	network.RegisterMessage(protos.CmdPlayerReqLogin, &protos.ReqPlayerLogin{})
	network.RegisterMessage(protos.CmdPlayerResLogin, &protos.ResPlayerLogin{})
	network.RegisterMessage(protos.CmdPlayerReqCreate, &protos.ReqPlayerCreate{})
	network.RegisterMessage(protos.CmdPlayerResCreate, &protos.ResPlayerCreate{})

	// 连接服务器
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal("连接服务器失败:", err)
	}
	defer conn.Close()
	fmt.Println("已连接到服务器:", address)

	//session := network.NewSession(&conn, &json.JsonCodec{})
	session := network.NewSession(&conn, &protobuf.ProtobufCodec{})

	//req := &protos.ReqJoinRoom{RoomId: 123, PlayerId: 123}
	req := &protos.ReqPlayerCreate{Name: "gforgame"}
	session.Send(req)
}

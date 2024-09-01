package main

import (
	"fmt"
	"io/github/gforgame/codec"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
	"log"
	"net"
)

func main() {
	// 服务器地址和端口
	address := "127.0.0.1:9090"

	network.RegisterMessage(1001, protos.ReqJoinRoom{})
	network.RegisterMessage(1002, protos.ReqChat{})
	network.RegisterMessage(1003, protos.ReqPlayerLogin{})

	// 连接服务器
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal("连接服务器失败:", err)
	}
	defer conn.Close()
	fmt.Println("已连接到服务器:", address)

	session := network.NewSession(&conn, &codec.JsonCodec{})

	req := protos.ReqJoinRoom{RoomId: 123, PlayerId: 123}

	session.Send(req)
}

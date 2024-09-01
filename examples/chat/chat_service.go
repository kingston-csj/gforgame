package chat

import (
	"fmt"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
)

type RoomService struct {
	network.Base
}

func NewRoomService() RoomService {
	return RoomService{}
}

func (rs RoomService) Init() {
	network.RegisterMessage(protos.CmdChatReqJoin, &protos.ReqJoinRoom{})
	network.RegisterMessage(protos.CmdChatReqChat, &protos.ReqChat{})
}

func (rs RoomService) JoinRoom(s *network.Session, msg *protos.ReqJoinRoom) error {
	fmt.Println(msg.PlayerId, "加入房间", msg.RoomId)
	return nil
}

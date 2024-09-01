package protos

const (
	CmdChatReqJoin = 1001
	CmdChatReqChat = 1002
)

type ReqPlayerLogin struct {
	Id int64
}

type ReqJoinRoom struct {
	RoomId int64

	PlayerId int64
}

type ReqChat struct {
	Id string
}

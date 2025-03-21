package protos

const (
	CmdChatReqJoin = 1001
	CmdChatReqChat = 1002

	CmdPlayerReqLogin  = 2001
	CmdPlayerResLogin  = 2002
	CmdPlayerReqCreate = 2003
	CmdPlayerResCreate = 2004

	CmdGmReqAction = 3001
	CmdGmResAction = 3002
)

type ReqPlayerLogin struct {
	Id  string
	Pwd string
}

type ResPlayerLogin struct {
	Succ bool
}

type ReqPlayerCreate struct {
	Name string
}

type ResPlayerCreate struct {
	Id int64
}

type ReqJoinRoom struct {
	RoomId int64

	PlayerId int64
}

type ReqChat struct {
	Id string
}

type ResBackpackInfo struct {
	Items []ItemInfo
}

type ItemInfo struct {
	Id    int32
	Count int32
}

type ReqGmAction struct {
	Topic  string
	Params string
}

type ResGmAction struct {
	Code int32
}

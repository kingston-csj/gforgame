package friend

import (
	mysqldb "io/github/gforgame/db"
	"io/github/gforgame/examples/context"
	"io/github/gforgame/examples/domain/player"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"

	"gorm.io/gorm"
)

type FriendController struct {
	network.Base
}

func NewFriendController() *FriendController {
	return &FriendController{}
}

func (c *FriendController) Init() {
	network.RegisterMessage(protos.CmdFriendReqApply, &protos.ReqFriendApply{})
	network.RegisterMessage(protos.CmdFriendResApply, &protos.ResFriendApply{})
	network.RegisterMessage(protos.CmdFriendReqDealApply, &protos.ReqFriendDealApplyRecord{})
	network.RegisterMessage(protos.CmdFriendResDealApply, &protos.ResFriendDealApplyRecord{})
	network.RegisterMessage(protos.CmdFriendReqDelete, &protos.ReqFriendDelete{})
	network.RegisterMessage(protos.CmdFriendResDelete, &protos.ResFriendDelete{})
	network.RegisterMessage(protos.CmdFriendReqQueryFriends, &protos.ReqFriendQueryMyFriends{})
	network.RegisterMessage(protos.CmdFriendResQueryFriends, &protos.ResFriendQueryMyFriends{})
	network.RegisterMessage(protos.CmdFriendReqSearchPlayers, &protos.ReqFriendSearchPlayers{})
	network.RegisterMessage(protos.CmdFriendResSearchPlayers, &protos.ResFriendSearchPlayers{})

	context.EventBus.Subscribe(events.PlayerLogin, func(data interface{}) {
		GetFriendService().RefreshClientInfo(data.(*playerdomain.Player))
	})

	// 自动建表
	err := mysqldb.Db.AutoMigrate(&player.Friend{})
	if err != nil {
		panic(err)
	}

	// 缓存数据读取
	dbLoader := func(key string) (interface{}, error) {
		var p player.Friend
		result := mysqldb.Db.First(&p, "id=?", key)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				// 未找到记录
				return nil, nil
			}
		}
		p.AfterFind(nil)
		return &p, nil
	}
	context.CacheManager.Register("friend", dbLoader)
}

func (c *FriendController) ReqSearchPlayers(s *network.Session, index int, msg *protos.ReqFriendSearchPlayers) *protos.ResFriendSearchPlayers {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	items := GetFriendService().SearchByKey(msg.Key)
	filteredItems := make([]*protos.FriendVo, 0)
	for _, item := range items {
		if item.Id != p.Id && !GetFriendService().IsFriend(p.Id, item.Id) {
			filteredItems = append(filteredItems, item)
		}
	}

	response := &protos.ResFriendSearchPlayers{
		Code:  0,
		Items: filteredItems,
	}
	return response
}

func (c *FriendController) ReqQueryFriends(s *network.Session, index int, msg *protos.ReqFriendQueryMyFriends) *protos.ResFriendQueryMyFriends {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	items := GetFriendService().QueryMyFriendVos(p.Id)
	response := &protos.ResFriendQueryMyFriends{
		Code:  0,
		Items: items,
	}
	return response
}

func (c *FriendController) ReqApply(s *network.Session, index int, msg *protos.ReqFriendApply) *protos.ResFriendApply {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	code := GetFriendService().ApplyFriend(p, msg.FriendId)
	response := &protos.ResFriendApply{
		Code: code,
	}
	return response
}

func (c *FriendController) ReqDealApply(s *network.Session, index int, msg *protos.ReqFriendDealApplyRecord) *protos.ResFriendDealApplyRecord {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	code := GetFriendService().DealApplyRecord(p, msg.ApplyId, msg.Status)
	response := &protos.ResFriendDealApplyRecord{
		Code: code,
	}
	return response
}

func (c *FriendController) ReqDelete(s *network.Session, index int, msg *protos.ReqFriendDelete) *protos.ResFriendDelete {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	code := GetFriendService().DeleteFriend(p, msg.FriendId)
	response := &protos.ResFriendDelete{
		Code: code,
	}
	return response
}

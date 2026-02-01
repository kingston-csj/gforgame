package route

import (
	mysqldb "io/github/gforgame/db"
	"io/github/gforgame/examples/context"
	"io/github/gforgame/examples/domain/player"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/service/friend"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"

	"gorm.io/gorm"
)

type FriendRoute struct {
	network.Base
	service *friend.FriendService
}

func NewFriendRoute() *FriendRoute {
	return &FriendRoute{
		service: friend.GetFriendService(),
	}
}

func (r *FriendRoute) Init() {
	context.EventBus.Subscribe(events.PlayerLogin, func(data interface{}) {
		r.service.RefreshClientInfo(data.(*playerdomain.Player))
	})

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

func (r *FriendRoute) ReqSearchPlayers(s *network.Session, index int32, msg *protos.ReqFriendSearchPlayers) *protos.ResFriendSearchPlayers {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	items := r.service.SearchByKey(msg.Key)
	filteredItems := make([]*protos.FriendVo, 0)
	for _, item := range items {
		if item.Id != p.Id && !r.service.IsFriend(p.Id, item.Id) {
			filteredItems = append(filteredItems, item)
		}
	}

	response := &protos.ResFriendSearchPlayers{
		Code:  0,
		Items: filteredItems,
	}
	return response
}

func (r *FriendRoute) ReqQueryFriends(s *network.Session, index int32, msg *protos.ReqFriendQueryMyFriends) *protos.ResFriendQueryMyFriends {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	items := r.service.QueryMyFriendVos(p.Id)
	response := &protos.ResFriendQueryMyFriends{
		Code:  0,
		Items: items,
	}
	return response
}

func (r *FriendRoute) ReqApply(s *network.Session, index int32, msg *protos.ReqFriendApply) *protos.ResFriendApply {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	code := r.service.ApplyFriend(p, msg.FriendId)
	response := &protos.ResFriendApply{
		Code: code,
	}
	return response
}

func (r *FriendRoute) ReqDealApply(s *network.Session, index int32, msg *protos.ReqFriendDealApplyRecord) *protos.ResFriendDealApplyRecord {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	code := r.service.DealApplyRecord(p, msg.ApplyId, msg.Status)
	response := &protos.ResFriendDealApplyRecord{
		Code: code,
	}
	return response
}

func (r *FriendRoute) ReqDelete(s *network.Session, index int32, msg *protos.ReqFriendDelete) *protos.ResFriendDelete {
	p := network.GetPlayerBySession(s).(*playerdomain.Player)
	code := r.service.DeleteFriend(p, msg.FriendId)
	response := &protos.ResFriendDelete{
		Code: code,
	}
	return response
}

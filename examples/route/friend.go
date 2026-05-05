package route

import (
	"github.com/forfun/gforgame/examples/context"
	"github.com/forfun/gforgame/examples/domain/player"
	playerdomain "github.com/forfun/gforgame/examples/domain/player"
	"github.com/forfun/gforgame/examples/events"
	mysqldb "github.com/forfun/gforgame/examples/infra/persistence"
	"github.com/forfun/gforgame/examples/protos"
	"github.com/forfun/gforgame/examples/service/friend"
	playerservice "github.com/forfun/gforgame/examples/service/player"
	"github.com/forfun/gforgame/network"

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
	p := playerservice.GetPlayerService().GetPlayerBySession(s)
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
	p := playerservice.GetPlayerService().GetPlayerBySession(s)
	items := r.service.QueryMyFriendVos(p.Id)
	response := &protos.ResFriendQueryMyFriends{
		Code:  0,
		Items: items,
	}
	return response
}

func (r *FriendRoute) ReqApply(s *network.Session, index int32, msg *protos.ReqFriendApply) *protos.ResFriendApply {
	p := playerservice.GetPlayerService().GetPlayerBySession(s)
	code := r.service.ApplyFriend(p, msg.FriendId)
	response := &protos.ResFriendApply{
		Code: code,
	}
	return response
}

func (r *FriendRoute) ReqDealApply(s *network.Session, index int32, msg *protos.ReqFriendDealApplyRecord) *protos.ResFriendDealApplyRecord {
	p := playerservice.GetPlayerService().GetPlayerBySession(s)
	code := r.service.DealApplyRecord(p, msg.ApplyId, msg.Status)
	response := &protos.ResFriendDealApplyRecord{
		Code: code,
	}
	return response
}

func (r *FriendRoute) ReqDelete(s *network.Session, index int32, msg *protos.ReqFriendDelete) *protos.ResFriendDelete {
	p := playerservice.GetPlayerService().GetPlayerBySession(s)
	code := r.service.DeleteFriend(p, msg.FriendId)
	response := &protos.ResFriendDelete{
		Code: code,
	}
	return response
}

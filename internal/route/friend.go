package route

import (
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/friend"
	playerService "github.com/forfun/gforgame/internal/service/player"
)

type FriendRoute struct {
	service *friend.FriendService
	player  *playerService.PlayerService
}

func NewFriendRoute(service *friend.FriendService, player *playerService.PlayerService) *FriendRoute {
	return &FriendRoute{
		service: service,
		player:  player,
	}
}

func (r *FriendRoute) ReqSearchPlayers(playerId string, index int32, msg *protos.ReqFriendSearchPlayers) *protos.ResFriendSearchPlayers {
	p := r.player.GetPlayer(playerId)
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

func (r *FriendRoute) ReqQueryFriends(playerId string, index int32, msg *protos.ReqFriendQueryMyFriends) *protos.ResFriendQueryMyFriends {
	p := r.player.GetPlayer(playerId)
	items := r.service.QueryMyFriendVos(p.Id)
	response := &protos.ResFriendQueryMyFriends{
		Code:  0,
		Items: items,
	}
	return response
}

func (r *FriendRoute) ReqApply(playerId string, index int32, msg *protos.ReqFriendApply) *protos.ResFriendApply {
	p := r.player.GetPlayer(playerId)
	code := r.service.ApplyFriend(p, msg.FriendId)
	response := &protos.ResFriendApply{
		Code: code,
	}
	return response
}

func (r *FriendRoute) ReqDealApply(playerId string, index int32, msg *protos.ReqFriendDealApplyRecord) *protos.ResFriendDealApplyRecord {
	p := r.player.GetPlayer(playerId)
	code := r.service.DealApplyRecord(p, msg.ApplyId, msg.Status)
	response := &protos.ResFriendDealApplyRecord{
		Code: code,
	}
	return response
}

func (r *FriendRoute) ReqDelete(playerId string, index int32, msg *protos.ReqFriendDelete) *protos.ResFriendDelete {
	p := r.player.GetPlayer(playerId)
	code := r.service.DeleteFriend(p, msg.FriendId)
	response := &protos.ResFriendDelete{
		Code: code,
	}
	return response
}

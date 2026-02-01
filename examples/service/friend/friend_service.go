package friend

import (
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/context"
	"io/github/gforgame/examples/domain/player"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/io"
	playerservice "io/github/gforgame/examples/service/player"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
	"io/github/gforgame/util"
	"sync"
	"time"
)

type FriendService struct {
}

var (
	instance *FriendService
	once     sync.Once
)

func GetFriendService() *FriendService {
	once.Do(func() {
		instance = &FriendService{}
	})
	return instance
}

func (s *FriendService) GetFriendEnt(playerId string) *player.Friend {
	cache, _ := context.CacheManager.GetCache("friend")
	cacheEntity, err := cache.Get(playerId)
	if err != nil {
		return nil
	}
	if cacheEntity == nil {
		return nil
	}
	friend, _ := cacheEntity.(*playerdomain.Friend)
	return friend
}

func (s *FriendService) GetFriendEntOrCreate(playerId string) *player.Friend {
	friend := s.GetFriendEnt(playerId)
	if friend == nil {
		friend = &player.Friend{}
		friend.Id = playerId
		friend.AfterFind(nil)
	}
	return friend
}

// 查询未处理申请记录
func (s *FriendService) QueryApplyRecords(playerId string) []*playerdomain.FriendApplyItem {
	friend := s.GetFriendEntOrCreate(playerId)
	if friend == nil {
		return nil
	}
	applies := make([]*playerdomain.FriendApplyItem, 0)
	for _, apply := range friend.Applies {
		if apply.Status == 0 && apply.TargetId == playerId {
			applies = append(applies, apply)
		}
	}
	return applies
}

// 查询我的好友
func (s *FriendService) QueryMyFriends(playerId string) []string {
	friend := s.GetFriendEntOrCreate(playerId)
	if friend == nil {
		return nil
	}
	friends := make([]string, 0)
	if friend != nil {
		for friendId := range friend.Friends {
			friends = append(friends, friendId)
		}
	}
	return friends
}

func (s *FriendService) QueryMyFriendVos(playerId string) []*protos.FriendVo {
	friend := s.GetFriendEntOrCreate(playerId)
	if friend == nil {
		return nil
	}
	friends := make([]*protos.FriendVo, 0)
	for friendId := range friend.Friends {
		friend := playerservice.GetPlayerService().GetPlayerProfileById(friendId)
		friends = append(friends, &protos.FriendVo{
			Id:       friendId,
			Name:     friend.Name,
			Fighting: int64(friend.Fight),
		})
	}
	return friends
}

func (s *FriendService) IsFriend(playerId string, friendId string) bool {
	friend := s.GetFriendEnt(playerId)
	if friend == nil {
		return false
	}
	return friend.IsFriend(friendId)
}

// 模糊搜索玩家 (key可能为名字或id)
func (s *FriendService) SearchByKey(key string) []*protos.FriendVo {
	playerIds := playerservice.GetPlayerService().FuzzySearchPlayers(key)
	friends := make([]*protos.FriendVo, 0)
	for _, playerId := range playerIds {
		profile := playerservice.GetPlayerService().GetPlayerProfileById(playerId)
		friends = append(friends, &protos.FriendVo{
			Id:       playerId,
			Name:     profile.Name,
			Fighting: int64(profile.Fight),
		})
	}
	// 如果是id,添加到结果中
	playerByName := playerservice.GetPlayerService().GetPlayerProfileById(key)
	if playerByName != nil {
		friends = append(friends, &protos.FriendVo{
			Id:       playerByName.Id,
			Name:     playerByName.Name,
			Fighting: int64(playerByName.Fight),
		})
	}
	return friends
}

func (s *FriendService) RefreshClientInfo(player *playerdomain.Player) {
	applyItems := s.QueryApplyRecords(player.Id)
	applyVos := make([]*protos.FriendApplyVo, 0, len(applyItems))
	for _, apply := range applyItems {
		fromPlayer := playerservice.GetPlayerService().GetPlayerProfileById(apply.FromId)
		applyVos = append(applyVos, &protos.FriendApplyVo{
			FromId:     fromPlayer.Id,
			FromName:   fromPlayer.Name,
			TargetId:   apply.TargetId,
			TargetName: fromPlayer.Name,
			Status:     int(apply.Status),
			Time:       apply.Time,
		})
	}

	friendVos := s.QueryMyFriendVos(player.Id)

	pushFriendInfo := &protos.PushFriendInfo{
		ApplyItems:  applyVos,
		FriendItems: friendVos,
		FriendSum:   len(friendVos),
	}

	io.NotifyPlayer(player, pushFriendInfo)
}

// 申请好友
func (s *FriendService) ApplyFriend(player *playerdomain.Player, friendId string) int {
	targetPlayer := playerservice.GetPlayerService().GetPlayerProfileById(friendId)
	if targetPlayer == nil {
		return constants.I18N_COMMON_NOT_FOUND
	}

	fromFriendEnt := s.GetFriendEntOrCreate(player.Id)
	if fromFriendEnt.IsFriend(friendId) {
		return constants.I18N_FRIEND_TIPS1
	}
	if fromFriendEnt.HasApplied(friendId) {
		return constants.I18N_FRIEND_TIPS2
	}

	applyItem := &playerdomain.FriendApplyItem{
		FromId:   player.Id,
		TargetId: friendId,
		Time:     time.Now().Unix(),
		Id:       util.GetNextID(),
	}
	fromFriendEnt.Applies[applyItem.Id] = applyItem
	s.SaveFriend(fromFriendEnt)
	// 复制一份给对方
	// 在线，考虑线程问题
	if network.IsOnline(friendId) {
		session := network.GetSessionByPlayerId(friendId)
		if session != nil {
			session.AsynTasks <- func() {
				targetFriendEnt := s.GetFriendEntOrCreate(friendId)
				targetFriendEnt.Applies[applyItem.Id] = applyItem
				s.SaveFriend(targetFriendEnt)
				s.RefreshClientInfo(playerservice.GetPlayerService().GetPlayer(friendId))
			}
		}
	} else {
		// 离线，直接修改数据
		targetFriendEnt := s.GetFriendEntOrCreate(friendId)
		targetFriendEnt.Applies[applyItem.Id] = applyItem
		s.SaveFriend(targetFriendEnt)
	}

	return 0
}

func (s *FriendService) DealApplyRecord(player *playerdomain.Player, applyId string, status int) int {
	friend := s.GetFriendEntOrCreate(player.Id)
	players := make([]*playerdomain.Player, 0)
	applyIds := make([]string, 0)
	if util.IsBlankString(applyId) {
		for applyId := range friend.Applies {
			applyIds = append(applyIds, applyId)
		}
	} else {
		apply := friend.Applies[applyId]
		if apply == nil {
			return constants.I18N_COMMON_NOT_FOUND
		}
		if apply.Status != 0 {
			return constants.I18N_COMMON_ILLEGAL_PARAMS
		}
		applyIds = append(applyIds, applyId)
	}

	for _, applyId := range applyIds {
		apply := friend.Applies[applyId]
		// 处理自己的申请
		s.dealApplyRecord0(player, applyId, apply.TargetId, status)

		// 处理对方的申请
		if network.IsOnline(apply.FromId) {
			session := network.GetSessionByPlayerId(apply.FromId)
			if session != nil {
				players = append(players, playerservice.GetPlayerService().GetPlayer(apply.FromId))
				session.AsynTasks <- func() {
					s.dealApplyRecord0(playerservice.GetPlayerService().GetPlayer(apply.TargetId), applyId, player.Id, status)
				}
			}
		} else {
			s.dealApplyRecord0(playerservice.GetPlayerService().GetPlayer(apply.TargetId), applyId, player.Id, status)
		}
	}

	if status == constants.FriendApplyStatusAgree {
		for _, player := range players {
			s.RefreshClientInfo(player)
		}
	}
	return 0
}

func (s *FriendService) dealApplyRecord0(owner *playerdomain.Player, applyId string, target string, status int) {
	friend := s.GetFriendEntOrCreate(owner.Id)
	apply := friend.Applies[applyId]
	if apply == nil {
		return
	}
	if apply.Status != 0 {
		return
	}
	friend.Applies[applyId].Status = int32(status)
	s.SaveFriend(friend)

	if status == constants.FriendApplyStatusAgree {
		friend.AddFriend(target)
		friend.ClearApply(owner.Id, target)
	}
}

func (s *FriendService) DeleteFriend(player *playerdomain.Player, friendId string) int {
	friend := s.GetFriendEntOrCreate(player.Id)
	if friend.IsFriend(friendId) {
		//  A删除B，B也删除A，客户端界面未刷新，直接返回成功即可
		return 0
	}
	friend.RemoveFriend(friendId)
	s.SaveFriend(friend)
	return 0
}

// 保存数据
func (s *FriendService) SaveFriend(friend *player.Friend) {
	cache, _ := context.CacheManager.GetCache("friend")
	cache.Set(friend.Id, friend)
	context.DbService.SaveToDb(friend)
}

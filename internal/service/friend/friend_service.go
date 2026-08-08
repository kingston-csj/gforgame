package friend

import (
	"time"

	"github.com/forfun/gforgame/common/eventbus"
	"github.com/forfun/gforgame/common/util/conv"
	"github.com/forfun/gforgame/internal/constants"
	"github.com/forfun/gforgame/internal/domain/player"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/events"
	"github.com/forfun/gforgame/internal/idgen"
	playerrepo "github.com/forfun/gforgame/internal/infra/repository/player"
	"github.com/forfun/gforgame/internal/io"
	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/internal/service/dispatch"
	mailservice "github.com/forfun/gforgame/internal/service/mail"
	"github.com/forfun/gforgame/network"
)

// 好友模块
type FriendService struct {
	playerRepo *playerrepo.PlayerRepository
	profile    *playerrepo.PlayerProfileService
	friendRepo playerdomain.FriendRepository
	mail       *mailservice.MailService
}

func NewFriendService(playerRepo *playerrepo.PlayerRepository, profile *playerrepo.PlayerProfileService, friendRepo playerdomain.FriendRepository, mail *mailservice.MailService) *FriendService {
	return &FriendService{
		playerRepo: playerRepo,
		profile:    profile,
		friendRepo: friendRepo,
		mail:       mail,
	}
}

func (s *FriendService) Init() {
	eventbus.Default().Subscribe(events.PlayerLogin, func(data any) {
		s.RefreshClientInfo(data.(*playerdomain.Player))
	})
}

func (s *FriendService) GetFriendEntOrCreate(playerId string) *playerdomain.Friend {
	friend := s.friendRepo.GetFriendEnt(playerId)
	if friend == nil {
		friend = &playerdomain.Friend{}
		friend.Id = playerId
		friend.AfterLoad()
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
		if friendId == playerId {
			continue
		}
		friend := s.profile.GetPlayerProfileById(friendId)
		friends = append(friends, &protos.FriendVo{
			Id:       friendId,
			Name:     friend.Name,
			Head:     friend.Head,
			Fighting: int64(friend.Fight),
		})
	}
	return friends
}

func (s *FriendService) IsFriend(playerId string, friendId string) bool {
	friend := s.friendRepo.GetFriendEnt(playerId)
	if friend == nil {
		return false
	}
	return friend.IsFriend(friendId)
}

// 模糊搜索玩家 (key可能为名字或id)
func (s *FriendService) SearchByKey(key string) []*protos.FriendVo {
	playerIds := s.profile.FuzzySearchPlayers(key)
	friends := make([]*protos.FriendVo, 0)
	for _, playerId := range playerIds {
		profile := s.profile.GetPlayerProfileById(playerId)
		friends = append(friends, &protos.FriendVo{
			Id:       playerId,
			Name:     profile.Name,
			Head:     profile.Head,
			Fighting: int64(profile.Fight),
		})
	}
	// 如果是id,添加到结果中
	playerByName := s.profile.GetPlayerProfileById(key)
	if playerByName != nil {
		friends = append(friends, &protos.FriendVo{
			Id:       playerByName.Id,
			Name:     playerByName.Name,
			Head:     playerByName.Head,
			Fighting: int64(playerByName.Fight),
		})
	}
	return friends
}

func (s *FriendService) RefreshClientInfo(player *playerdomain.Player) {
	applyItems := s.QueryApplyRecords(player.Id)
	applyVos := make([]*protos.FriendApplyVo, 0, len(applyItems))
	for _, apply := range applyItems {
		fromPlayer := s.profile.GetPlayerProfileById(apply.FromId)
		applyVos = append(applyVos, &protos.FriendApplyVo{
			FromId:     fromPlayer.Id,
			FromName:   fromPlayer.Name,
			FromHead:   fromPlayer.Head,
			TargetId:   apply.TargetId,
			TargetName: fromPlayer.Name,
			Status:     apply.Status,
			Time:       apply.Time,
		})
	}

	friendVos := s.QueryMyFriendVos(player.Id)

	pushFriendInfo := &protos.PushFriendInfo{
		ApplyItems:  applyVos,
		FriendItems: friendVos,
		FriendSum:   int32(len(friendVos)),
	}

	io.NotifyPlayer(player, pushFriendInfo)
}

// 申请好友
func (s *FriendService) ApplyFriend(player *playerdomain.Player, friendId string) int32 {
	targetPlayer := s.profile.GetPlayerProfileById(friendId)
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

	fromApplyItem := &playerdomain.FriendApplyItem{
		FromId:   player.Id,
		TargetId: friendId,
		Time:     time.Now().UnixMilli(),
		Id:       idgen.GetNextID(),
	}
	fromFriendEnt.Applies[fromApplyItem.Id] = fromApplyItem
	s.SaveFriend(fromFriendEnt)
	// 复制一份给对方(浅拷贝)
	targetApplyItem := *fromApplyItem
	// 在线，考虑线程问题
	if network.IsOnline(friendId) {
		dispatch.DispatchPlayerTask(friendId, func() {
			targetFriendEnt := s.GetFriendEntOrCreate(friendId)
			targetFriendEnt.AddApply(&targetApplyItem)
			s.SaveFriend(targetFriendEnt)
			s.RefreshClientInfo(s.playerRepo.GetPlayer(friendId))
		})
	} else {
		// 离线，直接修改数据库
		targetFriendEnt := s.GetFriendEntOrCreate(friendId)
		targetFriendEnt.AddApply(&targetApplyItem)
		s.SaveFriend(targetFriendEnt)
	}

	return 0
}

func (s *FriendService) DealApplyRecord(player *playerdomain.Player, applyId string, status int32) int32 {
	friend := s.GetFriendEntOrCreate(player.Id)
	applyIds := make([]string, 0)
	if conv.IsBlankString(applyId) {
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
		s.dealApplyRecord0(player, applyId, apply.FromId, status)
		target := s.playerRepo.GetPlayer(apply.FromId)
		// 处理对方的申请
		if network.IsOnline(apply.FromId) {
			dispatch.DispatchPlayerTask(apply.FromId, func() {
				s.dealApplyRecord0(target, applyId, player.Id, status)
				if status == constants.FriendApplyStatusAgree {
					s.RefreshClientInfo(target)
				}
			})
		} else {
			s.dealApplyRecord0(target, applyId, player.Id, status)
		}
	}

	if status == constants.FriendApplyStatusAgree {
		s.RefreshClientInfo(player)
	}
	return 0
}

func (s *FriendService) dealApplyRecord0(owner *playerdomain.Player, applyId string, target string, status int32) {
	friend := s.GetFriendEntOrCreate(owner.Id)
	apply := friend.Applies[applyId]
	if apply == nil {
		return
	}
	if apply.Status != 0 {
		return
	}
	friend.Applies[applyId].Status = int32(status)

	if status == constants.FriendApplyStatusAgree {
		friend.AddFriend(target)
		friend.ClearApply(owner.Id, target)
	}
	s.SaveFriend(friend)
}

func (s *FriendService) DeleteFriend(player *playerdomain.Player, friendId string) int32 {
	friend := s.GetFriendEntOrCreate(player.Id)
	if !friend.IsFriend(friendId) {
		//  A删除B，B也删除A，客户端界面未刷新，直接返回成功即可
		return 0
	}
	friend.RemoveFriend(friendId)
	s.SaveFriend(friend)
	s.RefreshClientInfo(player)
	// 处理对方
	task := func() {
		targetFriendEnt := s.GetFriendEntOrCreate(friendId)
		targetFriendEnt.RemoveFriend(player.Id)
		s.SaveFriend(targetFriendEnt)
		s.RefreshClientInfo(s.playerRepo.GetPlayer(friendId))
	}
	if network.IsOnline(friendId) {
		dispatch.DispatchPlayerTask(friendId, task)
	} else {
		task()
	}

	return 0
}

// 保存数据
func (s *FriendService) SaveFriend(friend *player.Friend) {
	s.friendRepo.SaveFriend(friend)
}

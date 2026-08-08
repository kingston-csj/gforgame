package player

import (
	"strings"

	"github.com/forfun/gforgame/common/container/hashmap"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"go.uber.org/dig"
)

// PlayerProfileService 玩家资料索引服务：维护 id->profile 与 id<->name 的内存索引，
// 供聊天、好友、竞技场等模块按 id/名字快速查询玩家展示资料。
// 它不负责玩家聚合的缓存与持久化，仅持有资料摘要索引，db 预加载委托 MySQLPlayerRepository。
type PlayerProfileService struct {
	inner          *MySQLPlayerRepository
	playerProfiles *hashmap.ConcurrentMap[string, *playerdomain.PlayerProfile]
	idNameMapper   *hashmap.SyncDualHashMap[string, string]
}

type PlayerProfileServiceParams struct {
	dig.In

	Inner *MySQLPlayerRepository
}

func NewPlayerProfileService(params PlayerProfileServiceParams) *PlayerProfileService {
	return &PlayerProfileService{
		inner:          params.Inner,
		playerProfiles: hashmap.NewConcurrentMap[string, *playerdomain.PlayerProfile](),
		idNameMapper:   hashmap.NewSyncDualHashMap[string, string](),
	}
}

// LoadPlayerProfiles 预热玩家资料索引。
func (s *PlayerProfileService) LoadPlayerProfiles() {
	profiles := s.inner.LoadPlayerProfilesFromDb()
	s.playerProfiles = hashmap.NewConcurrentMap[string, *playerdomain.PlayerProfile]()
	s.idNameMapper.Clear()
	for _, profile := range profiles {
		s.SavePlayerProfile(profile)
	}
}

func (s *PlayerProfileService) GetPlayerProfileById(playerId string) *playerdomain.PlayerProfile {
	v, ok := s.playerProfiles.Get(playerId)
	if ok {
		return v
	}
	return nil
}

func (s *PlayerProfileService) GetAllPlayerProfiles() []*playerdomain.PlayerProfile {
	result := make([]*playerdomain.PlayerProfile, 0, s.playerProfiles.Count())
	for _, profile := range s.playerProfiles.Values() {
		result = append(result, profile)
	}
	return result
}

func (s *PlayerProfileService) FuzzySearchPlayers(name string) []string {
	playerIds := make([]string, 0)
	profiles := s.playerProfiles.Values()
	for _, profile := range profiles {
		if strings.Contains(profile.Name, name) {
			playerIds = append(playerIds, profile.Id)
		}
	}
	return playerIds
}

// SavePlayerProfile 写入资料索引并维护名字映射。
func (s *PlayerProfileService) SavePlayerProfile(profile *playerdomain.PlayerProfile) {
	if profile == nil {
		return
	}
	s.playerProfiles.Set(profile.Id, profile)
	s.rebindPlayerName(profile.Id, profile.Name)
}

// UpdatePlayerProfile 更新资料索引中的展示字段。
func (s *PlayerProfileService) UpdatePlayerProfile(playerId string, name string, head int32) {
	profile, ok := s.playerProfiles.Get(playerId)
	if !ok || profile == nil {
		profile = &playerdomain.PlayerProfile{Id: playerId}
	}
	profile.Name = name
	profile.Head = head
	s.playerProfiles.Set(playerId, profile)
	s.rebindPlayerName(playerId, name)
}

func (s *PlayerProfileService) IsPlayerNameTaken(name string) bool {
	_, ok := s.idNameMapper.GetByValue(name)
	return ok
}

// rebindPlayerName 保证 id 和名字索引始终一一对应。
func (s *PlayerProfileService) rebindPlayerName(playerId string, name string) {
	if name == "" {
		return
	}
	if currentName, ok := s.idNameMapper.GetByKey(playerId); ok {
		if currentName == name {
			return
		}
		s.idNameMapper.DeleteByKey(playerId)
	}
	if currentPlayerId, ok := s.idNameMapper.GetByValue(name); ok && currentPlayerId != playerId {
		return
	}
	_ = s.idNameMapper.Put(playerId, name)
}

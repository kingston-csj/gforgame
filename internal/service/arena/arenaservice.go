package arena

import (
	"math"
	"sync"

	"github.com/forfun/gforgame/internal/config"
	"github.com/forfun/gforgame/internal/config/container"
	"github.com/forfun/gforgame/internal/constants"
	"github.com/forfun/gforgame/internal/context"
	player "github.com/forfun/gforgame/internal/domain/player"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
	"github.com/forfun/gforgame/internal/events"
	"github.com/forfun/gforgame/internal/idgen"
	"github.com/forfun/gforgame/internal/protos"
	heroService "github.com/forfun/gforgame/internal/service/hero"
	mailService "github.com/forfun/gforgame/internal/service/mail"
	"github.com/forfun/gforgame/internal/service/monthcard"
	playerService "github.com/forfun/gforgame/internal/service/player"
	"github.com/forfun/gforgame/internal/service/rank"
)

type ArenaService struct {
}

var (
	instance *ArenaService
	once     sync.Once
)

func GetArenaService() *ArenaService {
	once.Do(func() {
		instance = &ArenaService{}
	})
	return instance
}

// 申请挑战
func (s *ArenaService) Apply(player *player.Player, targetId string) int32 {
	target := playerService.GetPlayerService().GetPlayer(targetId)
	if target == nil {
		return constants.I18N_COMMON_NOT_FOUND
	}
	return 0
}

func queryDefenseTeam(player *player.Player) []*protos.HeroInfo {
	heroInfos := make([]*protos.HeroInfo, 0)
	for _, hero := range player.HeroBox.Heros {
		heroInfos = append(heroInfos, heroService.ToHeroVo(hero))
	}
	return heroInfos
}

// 每天免费战斗次数
func getTodayFreeTimes(player *player.Player) int32 {
	commonContainer := config.GetSpecificContainer[*container.CommonContainer]()
	// 每日竞技场战斗次数
	arenaDailyTimes := commonContainer.GetInt32Value("arenaDailyTimes")
	return arenaDailyTimes + monthcard.GetMonthCardService().GetExtraArenaTimes(player)
}

func (s *ArenaService) FightEnd(player *player.Player, target string, win bool) *protos.ResArenaFightEnd{
	res := &protos.ResArenaFightEnd{}
	targetPlayer := playerService.GetPlayerService().GetPlayer(target)
	if targetPlayer == nil {
		res.Code = constants.I18N_COMMON_NOT_FOUND
		return res
	}
	challengeTimes := player.DailyReset.ArenaTimes
	if challengeTimes < getTodayFreeTimes(player) {
		// 优先扣免费次数
		player.DailyReset.ArenaTimes++
	} else {
		if player.ArenaBox.Ticket < 0 {
			res.Code = constants.I18N_ARENA_TIPS1
			return res
		} else {
			player.ArenaBox.Ticket--
			player.DailyReset.ArenaTimes++
		}
	}
	res.MyInitScore = player.ArenaScore
	// 挑战者,增加积分,每日次数
	score1 := calcSettleScore(player, targetPlayer, win)
	newScore1 := player.ArenaScore + score1
	player.ArenaBox.ChallengeTimes++
	context.EventBus.Publish(events.AreaScoreChanged, &events.AreaScoreChangedEvent{
		PlayerEvent: events.PlayerEvent{
			Player: player,
		},
		Score:       score1,
	})
	context.EventBus.Publish(events.PassArena, &events.PassArenaEvent{
		PlayerEvent: events.PlayerEvent{
			Player: player,
		},
	})

	rankInfo1 := rank.GetRankService().GetMyRankInfo(rank.PlayerArenaRank, player.Id)
	rankParams1 := string(rankInfo1.Order)
	if rankInfo1.Order <= 0 {
		rankParams1 = "未上榜"
	}
	mailId := Ternary(win, constants.MailIdArenaFightWin, constants.MailIdArenaFightLose)
	mailService.GetMailService().SendSimpleMail2(player, mailId, 
		targetPlayer.Name, rankParams1, string(newScore1), rankParams1)
	addFightRecord(player, targetPlayer, score1, true, win)
	res.TargetInitScore = targetPlayer.ArenaScore

	// TODO 线程问题
	// 被挑战者,增加积分
	score2 := calcSettleScore(targetPlayer, player, !win)
	newScore2 := targetPlayer.ArenaScore + score2
	rankInfo2 := rank.GetRankService().GetMyRankInfo(rank.PlayerArenaRank, targetPlayer.Id)
	rankParams2 := string(rankInfo2.Order)
	if rankInfo2.Order <= 0 {
		rankParams2 = "未上榜"
	}
	mailId2 := Ternary(win, constants.MailIdArenaFightLose, constants.MailIdArenaFightWin)
	mailService.GetMailService().SendSimpleMail2(targetPlayer, mailId2, 
		player.Name, rankParams2, string(newScore2), rankParams2)
	addFightRecord(targetPlayer, player, score2, false, !win)
	context.EventBus.Publish(events.AreaScoreChanged, &events.AreaScoreChangedEvent{
		PlayerEvent: events.PlayerEvent{
			Player: targetPlayer,
		},
		Score:       score2,
	})
	
	res.MyChangedScore = score1
	res.TargetChangedScore = score2
	return res
}

// 封装通用的二选一函数（泛型版，支持任意类型）
func Ternary[V any](condition bool, valTrue V, valFalse V) V {
	if condition {
		return valTrue
	}
	return valFalse
}

func addFightRecord(player *player.Player, target *player.Player, score int32, isAttack bool,  win bool) {
	winner := player.Id
	if win {
		winner = target.Id
	}
	record := &playerdomain.MatchRecord{
		Id:         idgen.GetNextID(),
		OpponentId: target.Id,
		Score:      score,
		IsAttack:   isAttack,
		Winner:     winner,
	}
	player.ArenaBox.MatchRecords = append(player.ArenaBox.MatchRecords, record)
}

func calcSettleScore(player *player.Player, target *player.Player, win bool) int32 {
 	E := 1.0 / (1 + math.Pow(10, (float64(target.ArenaScore) - float64(player.ArenaScore)) / 400.0));
	score := 32 * (1.0 - E);
	if win {
		score = 32 * (1.0 - E);
	} else {
		score = 32 * -E;
	}
	return int32(math.Round(score));
}
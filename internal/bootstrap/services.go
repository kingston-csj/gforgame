package bootstrap

import (
	"reflect"

	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/internal/service/activity"
	"github.com/forfun/gforgame/internal/service/arena"

	"github.com/forfun/gforgame/internal/service/catalog"
	"github.com/forfun/gforgame/internal/service/chat"

	"github.com/forfun/gforgame/internal/service/friend"
	"github.com/forfun/gforgame/internal/service/gm"
	"github.com/forfun/gforgame/internal/service/hero"
	"github.com/forfun/gforgame/internal/service/item"
	"github.com/forfun/gforgame/internal/service/mail"
	"github.com/forfun/gforgame/internal/service/mall"
	"github.com/forfun/gforgame/internal/service/mixture"
	"github.com/forfun/gforgame/internal/service/monthcard"
	"github.com/forfun/gforgame/internal/service/player"

	"github.com/forfun/gforgame/internal/service/quest"
	"github.com/forfun/gforgame/internal/service/rank"
	"github.com/forfun/gforgame/internal/service/recharge"

	"github.com/forfun/gforgame/internal/service/signin"

	"github.com/forfun/gforgame/internal/service/vip"
)

type Services struct {
	Activity  *activity.ActivityService
	Arena     *arena.ArenaService
	Catalog   *catalog.CatalogService
	Chat      *chat.ChatService
	Friend    *friend.FriendService
	Gm        *gm.GmService
	Hero      *hero.HeroService
	Item      *item.ItemService
	Mail      *mail.MailService
	Mall      *mall.MallService
	Mixture   *mixture.MixtureService
	MonthCard *monthcard.MonthCardService
	Player    *player.PlayerService
	Quest     *quest.QuestService
	Rank      *rank.RankService
	Recharge  *recharge.RechargeService
	SignIn    *signin.SignInService
	Vip       *vip.VipService
}


// ServiceModule 定义 service 启动期初始化能力。
type ServiceModule interface {
	Init()
}

// InitServices 预热服务并完成跨模块注册（reward/consume ops 等）。
func InitServices() *Services {
	logger.Info("InitServices")
	s := &Services{}
	s.Mail = mail.NewMailService()
	s.Catalog = catalog.NewCatalogService()
	
	s.Quest = quest.NewQuestService()
	s.Player = player.NewPlayerService(s.Quest)
	s.Item = item.NewItemService(s.Player, s.Catalog)
	s.Hero = hero.NewHeroService(s.Player, s.Item)
	s.Friend = friend.NewFriendService(s.Player, s.Mail)
	s.Chat = chat.NewChatService(s.Player, s.Friend)
	s.MonthCard = monthcard.NewMonthCardService(s.Mail)
	s.Rank = rank.NewRankService(s.Player)
	s.Recharge = recharge.NewRechargeService()
	s.Vip = vip.NewVipService()
	s.Mall = mall.NewMallService()
	s.Mail = mail.NewMailService()
	s.Mixture = mixture.NewMixtureService()
	s.SignIn = signin.NewSignInService()
	s.Activity = activity.NewActivityService()
	s.Arena = arena.NewArenaService(s.Player, s.Rank, s.Mail)
	s.Gm = gm.NewGmService(&gm.GmDependencies{
		Player:    s.Player,
		Item:      s.Item,
		Quest:     s.Quest,
		Recharge:  s.Recharge,
		Mail:      s.Mail,
	})

	s.InitServiceModules()
	return s
}

// InitServiceModules 统一执行 service 启动初始化。
func (s *Services) InitServiceModules() {
	v := reflect.ValueOf(s).Elem()
	for i := 0; i < v.Type().NumField(); i++ {
		field := v.Field(i)
		if field.Kind() != reflect.Ptr {
			continue
		}
		if field.IsNil() {
			continue
		}
		// 判断是否实现了 ServiceModule 接口
		if module, ok := field.Interface().(ServiceModule); ok {
			module.Init()
		}
	}
}
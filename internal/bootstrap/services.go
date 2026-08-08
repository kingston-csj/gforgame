package bootstrap

import (
	"reflect"

	"github.com/forfun/gforgame/cache"
	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/internal/infra/persistence"
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

	"go.uber.org/dig"
)

type Services struct {
	dig.In    // 显式声明：这是一个「依赖输入组」
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

	// 创建 dig 容器
	c := dig.New()
	// 注册所有服务（只需要写构造函数！）
	registerServices(c)
	// 直接取出完整 Services（dig.In 参数对象只能用值类型依赖，不能用指针）
	var services Services
	if err := c.Invoke(func(built Services) {
		services = built
	}); err != nil {
		logger.Error("Failed to initialize services: %v", err)
		panic(err)
	}

	services.InitServiceModules()
	return &services
}

func registerServices(c *dig.Container) {
	// 基础服务
	_ = c.Provide(cache.NewCacheManager)
	_ = c.Provide(persistence.NewAsyncDbService)

	_ = c.Provide(item.NewItemService)
	_ = c.Provide(mail.NewMailService)
	_ = c.Provide(catalog.NewCatalogService)

	_ = c.Provide(quest.NewQuestService)
	_ = c.Provide(player.NewPlayerService)
	_ = c.Provide(item.NewItemService)
	_ = c.Provide(hero.NewHeroService)
	_ = c.Provide(friend.NewFriendService)
	_ = c.Provide(chat.NewChatService)
	_ = c.Provide(monthcard.NewMonthCardService)
	_ = c.Provide(rank.NewRankService)
	_ = c.Provide(recharge.NewRechargeService)
	_ = c.Provide(vip.NewVipService)
	_ = c.Provide(mall.NewMallService)
	_ = c.Provide(mixture.NewMixtureService)
	_ = c.Provide(signin.NewSignInService)
	_ = c.Provide(activity.NewActivityService)
	_ = c.Provide(arena.NewArenaService)

	// GM 特殊处理
	_ = c.Provide(newGmService)
}

func newGmService(
	player *player.PlayerService,
	item *item.ItemService,
	quest *quest.QuestService,
	recharge *recharge.RechargeService,
	mail *mail.MailService,
) *gm.GmService {
	return gm.NewGmService(&gm.GmDependencies{
		Player:   player,
		Item:     item,
		Quest:    quest,
		Recharge: recharge,
		Mail:     mail,
	})
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

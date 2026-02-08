package gm

import (
	"fmt"
	"io/github/gforgame/common"
	mysqldb "io/github/gforgame/db"
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/consume"
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/events"
	"io/github/gforgame/examples/reward"
	"io/github/gforgame/examples/service/item"
	"io/github/gforgame/examples/service/mail"
	playerservice "io/github/gforgame/examples/service/player"
	questservice "io/github/gforgame/examples/service/quest"
	"io/github/gforgame/examples/service/recharge"
	"io/github/gforgame/examples/service/scene"
	"io/github/gforgame/examples/system"
	"io/github/gforgame/logger"
	"io/github/gforgame/util"
	"io/github/gforgame/util/jsonutil"
	"sort"
	"strings"
	"sync"
)

type GmHandler func(player *playerdomain.Player, params string) *common.BusinessRequestException

type GmCommand struct {
	Topic       string
	Description string
	Example     string
	Handler     GmHandler
}

// GM模块
type GmService struct {
	commands map[string]*GmCommand
}

var (
	instance *GmService
	once     sync.Once
)

func GetGmService() *GmService {
	once.Do(func() {
		instance = &GmService{
			commands: make(map[string]*GmCommand),
		}
		instance.init()
	})
	return instance
}

func (s *GmService) init() {
	s.Register("help", "查看所有GM命令", "help", s.handleHelp)
	s.Register("reset", "重置玩家数据", "reset", handleReset)
	s.Register("level", "修改等级", "level 100", handleLevel)
	s.Register("add_items", "添加物品", "add_items 1001=1;1002=2", handleAddItems)
	s.Register("remove_items", "移除物品", "remove_items 1001=1", handleRemoveItems)
	s.Register("add_diamond", "添加钻石", "add_diamond 1000", handleAddDiamond)
	s.Register("add_gold", "添加金币", "add_gold 1000", handleAddGold)
	s.Register("quest", "完成任务", "quest 1001", handleQuest)
	s.Register("recharge", "模拟充值", "recharge 1", handleRecharge)
	s.Register("add_scene_items", "添加场景物品", "add_scene_items 1001=1", handleAddSceneItems)
	s.Register("remove_scene_items", "移除场景物品", "remove_scene_items 1001=1", handleRemoveSceneItems)
	s.Register("daily_reset", "触发每日重置", "daily_reset", handleDailyReset)
	s.Register("add_mail", "添加邮件", "add_mail 1001", handleAddMail)
	s.Register("clone", "克隆玩家", "clone 1001", handleClone)
}

func (s *GmService) Register(topic, desc, example string, handler GmHandler) {
	s.commands[topic] = &GmCommand{
		Topic:       topic,
		Description: desc,
		Example:     example,
		Handler:     handler,
	}
}

func (s *GmService) Dispatch(player *playerdomain.Player, topic string, params string) *common.BusinessRequestException {
	defer func() {
		if err := recover(); err != nil {
			logger.Error2("gm dispatch fail", err.(error))
		}
	}()

	cmd, ok := s.commands[topic]
	if !ok {
		logger.Error3(fmt.Sprintf("gm command not found: %s", topic))
		return common.NewBusinessRequestException(constants.I18N_GM_UNKNOWN_COMMAND)
	}

	err := cmd.Handler(player, params)

	// 触发玩家变更事件
	context.EventBus.Publish(events.PlayerEntityChange, player)

	return err
}

// ================= GM Handlers =================

func (s *GmService) handleHelp(player *playerdomain.Player, params string) *common.BusinessRequestException {
	var sb strings.Builder
	sb.WriteString("\n=== GM Commands ===\n")
	
	// 按Topic排序输出
	var topics []string
	for topic := range s.commands {
		topics = append(topics, topic)
	}
	sort.Strings(topics)

	for _, topic := range topics {
		cmd := s.commands[topic]
		sb.WriteString(fmt.Sprintf("%-20s : %s \n\tExample: %s\n", cmd.Topic, cmd.Description, cmd.Example))
	}
	logger.Info(sb.String())
	return nil
}

func handleReset(player *playerdomain.Player, params string) *common.BusinessRequestException {
	player.Reset()
	var scenes []playerdomain.Scene
	err := mysqldb.Db.Where(fmt.Sprintf("id like '%s%%'", player.Id)).Find(&scenes).Error
	if err != nil {
		logger.Error2("gm reset scene fail", err)
		return common.NewBusinessRequestException(constants.I18N_COMMON_INTERNAL_ERROR)
	}
	for _, item := range scenes {
		sceneId := item.Id[len(player.Id)+1:]
		cacheScene := scene.GetSceneService().GetSceneRecord(player.Id, sceneId)
		cacheScene.Data = ""
		scene.GetSceneService().SaveScene(cacheScene)
	}
	return nil
}

func handleLevel(player *playerdomain.Player, params string) *common.BusinessRequestException {
	player.Level = util.Int32Value(params)
	playerservice.GetPlayerService().GetPlayerProfileById(player.Id)
	return nil
}

func handleAddItems(player *playerdomain.Player, params string) *common.BusinessRequestException {
	itemIdMap, err := util.ToIntIntMap(params, ";", "=")
	if err != nil {
		return common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}
	for itemId, itemNum := range itemIdMap {
		item.GetItemService().AddByModelId(player, itemId, itemNum)
	}
	return nil
}

func handleRemoveItems(player *playerdomain.Player, params string) *common.BusinessRequestException {
	itemIdMap, err := util.ToIntIntMap(params, ";", "=")
	if err != nil {
		return common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}
	consums := &consume.AndConsume{}
	for itemId, itemNum := range itemIdMap {
		consums.Add(&consume.ItemConsume{
			ItemId: itemId,
			Amount: itemNum,
		})
	}
	if err := consums.Verify(player); err != nil {
		return err.(*common.BusinessRequestException)
	}
	consums.Consume(player, constants.ActionType_Gm)
	return nil
}

func handleAddDiamond(player *playerdomain.Player, params string) *common.BusinessRequestException {
	count, _ := util.StringToInt32(params)
	reward := &reward.CurrencyReward{
		Currency: "diamond",
		Amount:   count,
	}
	reward.Reward(player, constants.ActionType_Gm)
	return nil
}

func handleAddGold(player *playerdomain.Player, params string) *common.BusinessRequestException {
	count, _ := util.StringToInt32(params)
	reward := &reward.CurrencyReward{
		Currency: "gold",
		Amount:   count,
	}
	reward.Reward(player, constants.ActionType_Gm)
	return nil
}

func handleQuest(player *playerdomain.Player, params string) *common.BusinessRequestException {
	questId, _ := util.StringToInt32(params)
	questservice.GetQuestService().GmFinish(player, questId)
	return nil
}

func handleRecharge(player *playerdomain.Player, params string) *common.BusinessRequestException {
	rechargeId, _ := util.StringToInt32(params)
	recharge.GetRechargeService().Recharge(player, rechargeId)
	return nil
}

func handleAddSceneItems(player *playerdomain.Player, params string) *common.BusinessRequestException {
	itemIdMap, err := util.ToIntIntMap(params, ";", "=")
	if err != nil {
		return common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}
	for itemId, itemNum := range itemIdMap {
		err := item.GetSceneItemService().AddByModelId(player, itemId, itemNum)
		if err != nil {
			return common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
		}
	}
	return nil
}

func handleRemoveSceneItems(player *playerdomain.Player, params string) *common.BusinessRequestException {
	itemIdMap, err := util.ToIntIntMap(params, ";", "=")
	if err != nil {
		return common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}
	for itemId, itemNum := range itemIdMap {
		item.GetSceneItemService().UseByModelId(player, itemId, itemNum)
	}
	return nil
}

func handleDailyReset(player *playerdomain.Player, params string) *common.BusinessRequestException {
	system.PerformDailyUpdate()
	return nil
}


func handleAddMail(player *playerdomain.Player, params string) *common.BusinessRequestException {
	mail.GetMailService().SendSimpleMail(player, util.Int32Value(params))
	return nil
}

func handleClone(player *playerdomain.Player, params string) *common.BusinessRequestException {
	targetId := params
	if player.Id == targetId {
		return common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}
	// 复制玩家数据
	var to playerdomain.Player
	json, err := jsonutil.StructToJSON(player)
	if err != nil {
		return common.NewBusinessRequestException(constants.I18N_COMMON_INTERNAL_ERROR)
	}
	err = jsonutil.JsonToStruct(json, &to)
	if err != nil {
		return common.NewBusinessRequestException(constants.I18N_COMMON_INTERNAL_ERROR)
	}
	to.Id = targetId
	to.Name = playerservice.GetPlayerService().RandomName()
	playerservice.GetPlayerService().SavePlayer(&to)

	// 复制场景数据
	var scenes []playerdomain.Scene
	err = mysqldb.Db.Where(fmt.Sprintf("id like '%s%%'", player.Id)).Find(&scenes).Error
	if err != nil {
		logger.Error2("gm reset scene fail", err)
		return common.NewBusinessRequestException(constants.I18N_COMMON_INTERNAL_ERROR)
	}
	for _, item := range scenes {
		sceneId := item.Id[len(player.Id)+1:]
		fromScene := scene.GetSceneService().GetSceneRecord(player.Id, sceneId)
		toScene := scene.GetSceneService().GetOrCreateScene(targetId, sceneId)
		toScene.Data = fromScene.Data
		scene.GetSceneService().SaveScene(toScene)
	}


	return nil
}
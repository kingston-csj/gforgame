package gm

import (
	"fmt"
	"sort"
	"strings"

	commonerrors "github.com/forfun/gforgame/common/errors"
	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/internal/constants"
	playerdomain "github.com/forfun/gforgame/internal/domain/player"
)

type GmHandler func(player *playerdomain.Player, params string) *commonerrors.BusinessError

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

type GmRegistrar interface {
	RegisterTo(gm *GmService)
}

func NewGmService(deps *GmDependencies) *GmService {
	resolved := buildGmDependencies(deps)
	service := &GmService{
		commands: make(map[string]*GmCommand),
	}
	service.init(resolved)
	return service
}

func (s *GmService) init(deps *GmDependencies) {
	s.registerHandlers(
		NewSystemGmHandler(),
		NewPlayerGmHandler(deps.Player),
		NewItemGmHandler(deps.Item),
		NewQuestGmHandler(deps.Quest),
		NewRechargeGmHandler(deps.Recharge),
		NewMailGmHandler(deps.Mail),
	)
}

func (s *GmService) registerHandlers(handlers ...GmRegistrar) {
	for _, handler := range handlers {
		handler.RegisterTo(s)
	}
}

func (s *GmService) Register(topic, desc, example string, handler GmHandler) {
	s.commands[topic] = &GmCommand{
		Topic:       topic,
		Description: desc,
		Example:     example,
		Handler:     handler,
	}
}

func (s *GmService) Dispatch(player *playerdomain.Player, topic string, params string) *commonerrors.BusinessError {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("gm dispatch fail", err.(error))
		}
	}()

	cmd, ok := s.commands[topic]
	if !ok {
		logger.ErrorNoStack(fmt.Sprintf("gm command not found: %s", topic))
		return commonerrors.NewBusinessError(constants.I18N_GM_UNKNOWN_COMMAND)
	}
	// 去掉各各尾的换行符
	params = strings.TrimSuffix(params, "\r\n")
	params = strings.TrimSuffix(params, "\n")
	err := cmd.Handler(player, params)

	return err
}

func (s *GmService) handleHelp(player *playerdomain.Player, params string) *commonerrors.BusinessError {
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

package quest

import (
	"sync"

	"io/github/gforgame/common"

	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/constants"
	configdomain "io/github/gforgame/examples/domain/config"
	playerdomain "io/github/gforgame/examples/domain/player"
)

type QuestService struct {
	directors map[int32]QuestDirector 
	handlers map[int32]QuestHandler
}

var (
	instance *QuestService
	once     sync.Once
	
)

func GetQuestService() *QuestService {
	once.Do(func() {
		instance = &QuestService{}

		// 注册所有任务分类
		instance.directors = make(map[int32]QuestDirector)
		instance.directors[constants.QuestCategoryDaily] = &MainQuestDirector{}

		// 注册所有任务类型
		instance.handlers[constants.QuestTypeRecruit] = &RecruitQuestHandler{}
		for _,handler := range instance.handlers {
			handler.Init()
		}
	})
	return instance
}

func (s *QuestService) AcceptQuest(player *playerdomain.Player, questId int32) (*playerdomain.Quest, error) {
	if player.QuestBox.HasReceivedQuest(questId) {
		return nil, common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}
	questData := config.QueryById[configdomain.QuestData](questId)
	if questData == nil {
		return nil, common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)
	}
	quest := playerdomain.Quest{
		Id:   questId,
		Type: questData.Type,
		Progress: 0,
	}
	return &quest, nil
}

func GetQuestDirector(catalog int32) QuestDirector {
	item,ok := instance.directors[catalog] 
    if !ok {
        panic("quest director not found")
    }
    return item
}

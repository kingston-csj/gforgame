package player

import (
	"errors"
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/constants"
	configdomain "io/github/gforgame/examples/domain/config"
)

type QuestBox struct {

	// 当前接取的任务列表
	Doing map[int32]*Quest

	// 已完成的任务列表
	Finished map[int32]bool
}

func (q *QuestBox) AcceptNewQuest(quest *Quest) error {
	if _, ok := q.Doing[quest.Id]; ok {
		return errors.New("quest already accepting")
	}
	q.Doing[quest.Id] = quest
	return nil
}

// SelectUnFinishedQuestsByType 根据任务类型获取未完成任务列表
func (q *QuestBox) SelectUnFinishedQuestsByType(questType int32) []*Quest {
	var quests []*Quest
	for _, quest := range q.Doing {
		if quest.Type == questType && quest.Status == constants.QuestStatusInit {
			quests = append(quests, quest)
		}
	}
	return quests
}


// SelectUnFinishedQuestsByCategory 根据任务分类获取未完成任务列表
func (q *QuestBox) SelectUnFinishedQuestsByCategory(questCategory int32) []*Quest {
	var quests []*Quest
	for _, quest := range q.Doing {
		questData := config.QueryById[configdomain.QuestData](quest.Id)
		if questData.Category == questCategory && quest.Status == constants.QuestStatusInit {
			quests = append(quests, quest)
		}
	}
	return quests
}

// ClearQuestsByCategory 清除指定分类的任务
func (q *QuestBox) ClearQuestsByCategory(questCategory int32) {
	for id, quest := range q.Doing {
		questData := config.QueryById[configdomain.QuestData](quest.Id)
		if questData.Category == questCategory {
			delete(q.Doing, id)
		}
	}
}

// HasReceivedQuest 是否接取了指定任务
func (q *QuestBox) HasReceivedQuest(questId int32) bool {
	_, ok := q.Doing[questId]
	return ok
}

// AddFinishedQuest 添加已完成任务
func (q *QuestBox) AddFinishedQuest(questId int32) {
	q.Finished[questId] = true
}


// IsFinished 是否完成了指定任务
func (q *QuestBox) IsFinished(questId int32) bool {
	_, ok := q.Finished[questId]
	return ok
}


// RemoveQuest 移除指定任务
func (q *QuestBox) RemoveQuest(questId int32) {
	delete(q.Doing, questId)
}

// GetCurrentMainQuest 获取当前接取的主线任务
func (q *QuestBox) GetCurrentMainQuest() *Quest {
	for _, quest := range q.Doing {
		questData := config.QueryById[configdomain.QuestData](quest.Id)
		if questData.Category == constants.QuestCategoryMain && quest.Status == constants.QuestStatusInit {
			return quest
		}
	}
	return nil
}


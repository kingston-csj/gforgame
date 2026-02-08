package quest

import (
	"io/github/gforgame/examples/config"
	"io/github/gforgame/examples/config/container"
	"io/github/gforgame/examples/constants"
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/io"
	"io/github/gforgame/protos"
	"io/github/gforgame/util"
	"io/github/gforgame/util/timeutil"
	"time"
)

// 委托任务类别
type EntrustQuestDirector struct {
	baseQuestDirector
}

func NewEntrustQuestDirector() *EntrustQuestDirector {
	return &EntrustQuestDirector{}
}

/// 实现QuestDirector接口

func (d *EntrustQuestDirector) OnPlayerLogin(player *playerdomain.Player) {
	c := config.GetSpecificContainer[*container.QuestContainer]()
	questDatas := c.GetRecordsByIndex("Category", d.GetCategoryType())
	now := time.Now().UnixMilli()
	for _, questData := range questDatas {
		if !player.QuestBox.HasReceivedQuest(questData.Id) {
			GetQuestService().AcceptQuest(player, questData.Id)
		} else {
			quest := player.QuestBox.GetQuest(questData.Id)
			if quest.Status == constants.QuestStatusDoing {
				cost := now - quest.AcceptTime
				// 超过时间，自动完成
				if cost >= int64(util.Int32Value(questData.Extra) * int32(timeutil.MILLIS_PER_SECOND))  {
					quest.AddProgress(1)
					handler := GetQuestHandler(questData.Type)
					if handler != nil {
						handler.CheckProgress(player, quest)
					}
				} else {
					diff := int64(util.Int32Value(questData.Extra)) * timeutil.MILLIS_PER_SECOND - cost
					context.TaskScheduler.Schedule(func() {
						quest.AddProgress(1)
						handler := GetQuestHandler(questData.Type)
						if handler != nil {
							handler.CheckProgress(player, quest)
						}
					}, diff)
				}
			}
		}
	}

	quests := player.QuestBox.SelectUnFinishedQuestsByCategory(d.GetCategoryType())
	questVos := make([]*protos.QuestVo, 0, len(quests))
	for _, quest := range quests {
		questVos = append(questVos, quest.ToVo())
	}
	push := &protos.PushQuestEntrustInfo{
		Quests: questVos,
	}
	io.NotifyPlayer(player, push)
}

// 玩家完成任务后触发
func (d *EntrustQuestDirector) AfterTakeReward(player *playerdomain.Player, quest *playerdomain.Quest) {
	// 任务完成，重置该任务
	quest.Reset()
	// 解除绑定对应的英雄
	for _, hero := range player.HeroBox.Heros {
		if hero.EntrustQuestId == quest.Id {
			hero.EntrustQuestId = 0
		}
	}
	d.OnPlayerLogin(player)
}

// 获取任务类型
func (d *EntrustQuestDirector) GetCategoryType() int32 {
	return constants.QuestCategoryEntrust
}
 








package quest

import (
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
)

type QuestController struct {
	network.Base
}

func NewQuestController() *QuestController {
	return &QuestController{}
}


func (c *QuestController) Init() {
	network.RegisterMessage(protos.CMD_REQ_QUEST_REWARD, &protos.ReqQuestTakeReward{})
	network.RegisterMessage(protos.CMD_REQ_QUEST_PROGRESS_REWARD, &protos.ReqQuestTakeProgressReward{})

	network.RegisterMessage(protos.CMD_PUSH_QUEST_AUTO_REWARD, &protos.PushQuestAutoTakeReward{})
	network.RegisterMessage(protos.CMD_PUSH_DAILY_QUEST, &protos.PushQuestDailyInfo{})
	network.RegisterMessage(protos.CMD_PUSH_UPDATE_QUEST, &protos.PushQuestRefreshVo{})
	network.RegisterMessage(protos.CMD_RES_REPLACE_QUEST, &protos.PushQuestReplace{})
	network.RegisterMessage(protos.CMD_PUSH_WEEKLY_QUEST, &protos.PushQuestWeeklyInfo{})
	network.RegisterMessage(protos.CMD_REQ_QUEST_ALL_REWARD, &protos.ReqQuestTakeAllRewards{})
	network.RegisterMessage(protos.CMD_RES_QUEST_ALL_REWARD, &protos.ResQuestTakeAllRewards{})
	network.RegisterMessage(protos.CMD_RES_QUEST_PROGRESS_REWARD, &protos.ResQuestTakeProgressReward{})
	network.RegisterMessage(protos.CMD_RES_QUEST_REWARD, &protos.ResQuestTakeReward{})
}
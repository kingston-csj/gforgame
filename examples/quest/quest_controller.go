package quest

import "io/github/gforgame/network"

type QuestController struct {
	network.Base
}

func NewQuestController() *QuestController {
	return &QuestController{}
}

func (c *QuestController) Init() {

}
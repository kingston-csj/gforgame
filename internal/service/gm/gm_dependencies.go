package gm

import (
	"fmt"

	"github.com/forfun/gforgame/internal/service/item"
	"github.com/forfun/gforgame/internal/service/mail"
	playerservice "github.com/forfun/gforgame/internal/service/player"
	questservice "github.com/forfun/gforgame/internal/service/quest"
	"github.com/forfun/gforgame/internal/service/recharge"
)

type GmDependencies struct {
	Player    *playerservice.PlayerService
	Item      *item.ItemService
	Quest     *questservice.QuestService
	Recharge  *recharge.RechargeService
	Mail      *mail.MailService
}

var defaultDependencies *GmDependencies

func SetDefaultDependencies(deps *GmDependencies) {
	defaultDependencies = deps
}

func buildGmDependencies(deps *GmDependencies) *GmDependencies {
	if deps == nil {
		deps = &GmDependencies{}
	}
	if defaultDependencies != nil {
		if deps.Player == nil {
			deps.Player = defaultDependencies.Player
		}
		if deps.Item == nil {
			deps.Item = defaultDependencies.Item
		}
		if deps.Quest == nil {
			deps.Quest = defaultDependencies.Quest
		}
		if deps.Recharge == nil {
			deps.Recharge = defaultDependencies.Recharge
		}
		if deps.Mail == nil {
			deps.Mail = defaultDependencies.Mail
		}
	}
	mustNotNil("gm.Player", deps.Player)
	mustNotNil("gm.Quest", deps.Quest)
	mustNotNil("gm.Recharge", deps.Recharge)
	mustNotNil("gm.Mail", deps.Mail)
	return deps
}

func mustNotNil(name string, value any) {
	if value == nil {
		panic(fmt.Sprintf("%s dependency is not configured", name))
	}
}

package system

import (
	"fmt"
	playerdomain "io/github/gforgame/examples/domain/player"
	sessionmanager "io/github/gforgame/examples/session"
	"io/github/gforgame/logger"
	"time"

	"github.com/go-co-op/gocron"
)

func zeroTask() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(fmt.Errorf("零点任务执行错误: %v", err))
		}
	}()
	now := time.Now()
	fmt.Println("零点任务执行:", now) 

	sessions := sessionmanager.GetAllSessions()
	for _, session := range sessions {
		session.AsynTasks <- func() {
			p := sessionmanager.GetPlayerBySession(session)
			if p == nil {
				return
			}
			player := p.(*playerdomain.Player)
			player.DailyReset.Reset(now.Unix())
			fmt.Println("执行玩家零点任务:", now)
		}
	}
}

func StartSystemTask() {
	go func() {
		s := gocron.NewScheduler(time.UTC)
		// 每天0点执行任务
		s.Every(1).Day().At("00:00").Do(zeroTask)
		s.StartBlocking()
	}()
}

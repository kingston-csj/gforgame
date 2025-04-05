package player

import (
	"fmt"

	mysqldb "io/github/gforgame/db"
	"io/github/gforgame/examples/context"
	playerdomain "io/github/gforgame/examples/domain/player"
	"io/github/gforgame/examples/session"
	"io/github/gforgame/logger"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
	"io/github/gforgame/util"

	"gorm.io/gorm"
)

type PlayerController struct {
	network.Base
}

func NewPlayerController() *PlayerController {
	return &PlayerController{}
}

func (ps *PlayerController) Init() {
	network.RegisterMessage(protos.CmdPlayerReqLogin, &protos.ReqPlayerLogin{})
	network.RegisterMessage(protos.CmdPlayerResLogin, &protos.ResPlayerLogin{})

	network.RegisterMessage(protos.CmdPlayerReqCreate, &protos.ReqPlayerCreate{})
	network.RegisterMessage(protos.CmdPlayerResCreate, &protos.ResPlayerCreate{})

	// 自动建表
	err := mysqldb.Db.AutoMigrate(&playerdomain.Player{})
	if err != nil {
		panic(err)
	}

	// 缓存数据读取
	dbLoader := func(key string) (interface{}, error) {
		var p playerdomain.Player
		result := mysqldb.Db.First(&p, "id=?", key)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				// 未找到记录
				return nil, nil
			}
		}
		p.AfterFind(nil)
		return &p, nil
	}
	context.CacheManager.Register("player", dbLoader)
}

func (ps *PlayerController) ReqLogin(s *network.Session, index int, msg *protos.ReqPlayerLogin) {
	player := GetPlayerService().GetOrCreatePlayer(msg.Id)
	fmt.Println("登录成功，id为：", player.Id)

	// 添加session
	session.AddSession(s, player)

	s.Send(&protos.ResPlayerLogin{Succ: true}, index)

	context.EventBus.Publish("player_login", player)
}

func (ps *PlayerController) ReqCreate(s *network.Session, msg *protos.ReqPlayerCreate) {
	id := util.GetNextId()
	player := &playerdomain.Player{}
	player.Id = id
	player.Name = msg.Name
	mysqldb.Db.Create(&player)

	logger.Log(logger.Player, "Id", player.Id, "name", player.Name)
	fmt.Printf(player.Name)
}

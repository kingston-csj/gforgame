package player

import (
	"fmt"
	mysqldb "io/github/gforgame/db"
	"io/github/gforgame/log"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
	"io/github/gforgame/util"
)

type PlayerService struct {
	network.Base
}

func NewPlayerService() PlayerService {
	return PlayerService{}
}

func (rs PlayerService) Init() {
	network.RegisterMessage(protos.CmdPlayerReqLogin, &protos.ReqPlayerLogin{})
	network.RegisterMessage(protos.CmdPlayerResLogin, &protos.ResPlayerLogin{})

	network.RegisterMessage(protos.CmdPlayerReqCreate, &protos.ReqPlayerCreate{})
	network.RegisterMessage(protos.CmdPlayerResCreate, &protos.ResPlayerCreate{})

	// 自动建表
	mysqldb.Db.AutoMigrate(&Player{})
}

func (rs PlayerService) ReqLogin(s *network.Session, msg *protos.ReqPlayerLogin) interface{} {
	var player Player
	mysqldb.Db.First(&player, "id=?", 1001)
	fmt.Println(msg.Id, "登录成功，姓名为：", player.Name)
	//return &protos.ResPlayerLogin{Succ: true}
	return nil
}

func (rs PlayerService) ReqCreate(s *network.Session, msg *protos.ReqPlayerCreate) {
	id := util.GetNextId()
	player := &Player{Id: id, Name: msg.Name}
	mysqldb.Db.Create(&player)

	log.Log(log.Player, "Id", player.Id, "name", player.Name)

	fmt.Printf(player.Name)
}

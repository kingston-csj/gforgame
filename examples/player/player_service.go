package player

import (
	"fmt"
	mysqldb "io/github/gforgame/db"
	"io/github/gforgame/examples/context"
	"io/github/gforgame/log"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
	"io/github/gforgame/util"
	"strconv"
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

	dbLoader := func(key string) (interface{}, error) {
		var p Player
		mysqldb.Db.First(&p, "id=?", key)
		return &p, nil
	}
	context.CacheManager.Register("player", dbLoader)
}

func (rs PlayerService) ReqLogin(s *network.Session, msg *protos.ReqPlayerLogin) interface{} {
	cache, err := context.CacheManager.GetCache("player")
	if err != nil {
		log.Error(err)
	}
	entity, _ := cache.Get(strconv.FormatInt(msg.Id, 10))
	var player, _ = entity.(*Player)
	fmt.Println(msg.Id, "登录成功，姓名为：", player.Name)
	return &protos.ResPlayerLogin{Succ: true}
}

func (rs PlayerService) ReqCreate(s *network.Session, msg *protos.ReqPlayerCreate) {
	id := util.GetNextId()
	player := &Player{Id: id, Name: msg.Name}
	mysqldb.Db.Create(&player)

	log.Log(log.Player, "Id", player.Id, "name", player.Name)

	fmt.Printf(player.Name)
}

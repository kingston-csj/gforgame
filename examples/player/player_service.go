package player

import (
	"fmt"
	mysqldb "io/github/gforgame/db"
	"io/github/gforgame/examples/context"
	"io/github/gforgame/log"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
	"io/github/gforgame/util"
)

type Service struct {
	network.Base
}

func NewPlayerService() *Service {
	return &Service{}
}

func (ps *Service) Init() {
	network.RegisterMessage(protos.CmdPlayerReqLogin, &protos.ReqPlayerLogin{})
	network.RegisterMessage(protos.CmdPlayerResLogin, &protos.ResPlayerLogin{})

	network.RegisterMessage(protos.CmdPlayerReqCreate, &protos.ReqPlayerCreate{})
	network.RegisterMessage(protos.CmdPlayerResCreate, &protos.ResPlayerCreate{})

	// 自动建表
	err := mysqldb.Db.AutoMigrate(&Player{})
	if err != nil {
		panic(err)
	}

	// 缓存数据读取
	dbLoader := func(key string) (interface{}, error) {
		var p Player
		mysqldb.Db.First(&p, "id=?", key)
		return &p, nil
	}
	context.CacheManager.Register("player", dbLoader)
}

func (ps *Service) ReqLogin(s *network.Session, msg *protos.ReqPlayerLogin) interface{} {
	cache, err := context.CacheManager.GetCache("player")
	if err != nil {
		log.Error(err)
	}
	cacheEntity, _ := cache.Get(msg.Id)
	player, _ := cacheEntity.(*Player)
	player.Name = "hello,gforgame"
	player.Level = 999

	ps.SavePlayer(player)

	fmt.Println(msg.Id, "登录成功，姓名为：", player.Name)
	return &protos.ResPlayerLogin{Succ: true}
}

func (ps *Service) ReqCreate(s *network.Session, msg *protos.ReqPlayerCreate) {
	id := util.GetNextId()
	player := &Player{}
	player.Id = id
	player.Name = msg.Name
	mysqldb.Db.Create(&player)

	log.Log(log.Player, "Id", player.Id, "name", player.Name)
	fmt.Printf(player.Name)
}

func (ps *Service) SavePlayer(player interface{}) {
	entity, ok := player.(mysqldb.Entity)
	if !ok {
		panic("not player")
	}
	context.DbService.SaveToDb(entity)
}

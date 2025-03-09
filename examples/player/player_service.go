package player

import (
	"errors"
	"fmt"
	mysqldb "io/github/gforgame/db"
	"io/github/gforgame/examples/context"
	"io/github/gforgame/logger"
	"io/github/gforgame/network"
	"io/github/gforgame/protos"
	"io/github/gforgame/util"
)

var ErrNotFound = errors.New("record not found")
var ErrCast = errors.New("cast exception")

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

func (ps *Service) ReqLogin(s *network.Session, index int, msg *protos.ReqPlayerLogin) interface{} {
	cache, err := context.CacheManager.GetCache("player")
	if err != nil {
		return &protos.ResPlayerLogin{Succ: false}
	}
	playerId := msg.Id
	cacheEntity, err := cache.Get(playerId)
	var player *Player
	if errors.Is(err, ErrNotFound) {
		// 新增玩家
		player = &Player{}
		player.Name = ""
		player.Level = 1
		player.Id = playerId
		cache.Set(playerId, player)

		ps.SavePlayer(player)
		return &protos.ResPlayerLogin{Succ: true}
	} else if err != nil {
		return &protos.ResPlayerLogin{Succ: false}
	} else {
		player, _ = cacheEntity.(*Player)
	}

	fmt.Println("登录成功，id为：", player.Id)
	return &protos.ResPlayerLogin{Succ: true}
}

func (ps *Service) ReqCreate(s *network.Session, msg *protos.ReqPlayerCreate) {
	id := util.GetNextId()
	player := &Player{}
	player.Id = id
	player.Name = msg.Name
	mysqldb.Db.Create(&player)

	logger.Log(logger.Player, "Id", player.Id, "name", player.Name)
	fmt.Printf(player.Name)
}

func (ps *Service) SavePlayer(player interface{}) {
	entity, ok := player.(mysqldb.Entity)
	if !ok {
		panic("not player")
	}
	context.DbService.SaveToDb(entity)
}

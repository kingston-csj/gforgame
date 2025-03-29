package cross

import (
	"errors"
	"strconv"

	playerdomain "io/github/gforgame/examples/domain/player"
)

type TransferType int

const (
	Island TransferType = iota // 抢占公岛，类型为0
	Moba                       // Moba pvp
)

var services map[TransferType]Transfer = make(map[TransferType]Transfer)

type Transfer interface {
	// 能否登录到跨服节点
	CanTransfer(p *playerdomain.Player) int
	// 获取目标服务器id
	GetTargetServerId(p *playerdomain.Player) uint32
	// 本服玩家进入跨服场景
	LocalEnterScene(p *playerdomain.Player) error
	// 玩家跨服成功后进入跨服场景
	RemoteEnterScene(p *playerdomain.Player) error
}

func RegisterTransfer(kind TransferType, t Transfer) {
	_, existed := services[kind]
	if existed {
		panic("cmd duplicated: " + strconv.Itoa(int(kind)))
	}
	services[kind] = t
}

func GetTransfer(kind TransferType) (Transfer, error) {
	value, ok := services[kind]
	if ok {
		return value, nil
	} else {
		return nil, errors.New("GetMessageCmd not found")
	}
}

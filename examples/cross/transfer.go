package cross

import (
	"errors"
	"io/github/gforgame/examples/player"
	"strconv"
)

type TransferType int

const (
	Island TransferType = iota // 抢占公岛，类型为0
	Moba                       // Moba pvp
)

var (
	services map[TransferType]Transfer = make(map[TransferType]Transfer)
)

type Transfer interface {
	// 能否登录到跨服节点
	CanTransfer(p *player.Player) int
	// 获取目标服务器id
	GetTargetServerId(p *player.Player) uint32
	// 本服玩家进入跨服场景
	LocalEnterScene(p *player.Player) error
	// 玩家跨服成功后进入跨服场景
	RemoteEnterScene(p *player.Player) error
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

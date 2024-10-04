package cross

import (
	"context"
	"fmt"
	"io/github/gforgame/codec/json"
	"io/github/gforgame/config"
	"io/github/gforgame/examples/player"
	"io/github/gforgame/logger"
	"io/github/gforgame/network/rpc"
)

func PlayerLoginRemote(p *player.Player, kind TransferType) {
	t, err := GetTransfer(kind)
	if err != nil {
		logger.Error(err)
		return
	}
	if result := t.CanTransfer(p); result > 0 {
		return
	}
	err = transferToRemote(p, t)
	if err != nil {
		logger.Error(err)
		return
	}
}

func transferToRemote(p *player.Player, t Transfer) error {
	localSid := config.ServerConfig.ServerId
	targetSid := t.GetTargetServerId(p)

	// 目标服是本服，直接进场景
	if localSid == targetSid {
		err := t.LocalEnterScene(p)
		if err != nil {
			return err
		}
	} else {
		// 目标服不是本服，走跨服登录
		err := doEnterCrossServer(p, t)
		if err != nil {
			return err
		}
	}

	return nil
}

func doEnterCrossServer(p *player.Player, t Transfer) error {
	codec := json.NewSerializer()
	data, err := codec.Encode(p)
	if err != nil {
		return err
	}

	rpcClient, err := rpc.GetOrCreateClient(1001)
	if err != nil {
		panic(err)
	}
	resp, err := rpcClient.EnterRemote(context.Background(), &rpc.PlayerCrossRequest{Data: data})
	if err != nil {
		return err
	}
	fmt.Println("rpc客户端收到消息：", resp.Message)
	return nil
}

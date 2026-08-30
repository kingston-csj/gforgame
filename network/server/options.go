package server

import (
	"github.com/forfun/gforgame/codec"
	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/dispatch"
	"github.com/forfun/gforgame/network/protocol"
)

type BaseServerOptions struct {
	Name            string // 服务器名称
	ServiceAddr     string // current server service address
	MessageCodec    codec.MessageCodec
	IoDispatch      dispatch.IoDispatch
	Router          *network.MessageRoute // 路由表
	PayloadMode     protocol.PayloadMode
	UseGateway      bool // 是否使用网关模式
}

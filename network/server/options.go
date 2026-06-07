package server

import (
	"github.com/forfun/gforgame/codec"
	"github.com/forfun/gforgame/network"
)

type BaseServerOptions struct {
	Name            string // 服务器名称
	ServiceAddr     string // current server service address
	MessageCodec    codec.MessageCodec
	IoDispatch      network.IoDispatch
	DispatchWorkers int32 // 网关模式下， 处理io任务的worker程数
	Router          *network.MessageRoute // 路由表
	PayloadMode     network.PayloadMode
	UseGateway      bool // 是否使用网关模式
}

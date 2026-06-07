package ws

import (
	"github.com/forfun/gforgame/codec"
	"github.com/forfun/gforgame/network"
	serverpkg "github.com/forfun/gforgame/network/server"
)

type Options struct {
	serverpkg.BaseServerOptions
	wsPath string
}

type Option func(*Options)

// WithAddress 指定ip地址
func WithAddress(addr string) Option {
	return func(opt *Options) {
		opt.ServiceAddr = addr
	}
}

// WithIoDispatch 消息处理链
func WithIoDispatch(dispatch network.IoDispatch) Option {
	return func(opt *Options) {
		opt.IoDispatch = dispatch
	}
}

// WithCodec 消息编解码
func WithCodec(codec codec.MessageCodec) Option {
	return func(opt *Options) {
		opt.MessageCodec = codec
	}
}

// WithWsPath 设置websocket的路径
func WithWsPath(path string) Option {
	return func(opt *Options) {
		opt.wsPath = path
	}
}

// WithRouter 消息路由器
func WithRouter(r *network.MessageRoute) Option {
	return func(opt *Options) {
		opt.Router = r
	}
}

// WithPayloadMode 设置消息体处理模式（解析 or 原始转发）
func WithPayloadMode(mode network.PayloadMode) Option {
	return func(opt *Options) {
		opt.PayloadMode = mode
	}
}

// WithDispatchWorkers 设置每条连接的消息消费 worker 数。
// 值 <= 0 时按 1 处理。
func WithDispatchWorkers(n int32) Option {
	return func(opt *Options) {
		opt.DispatchWorkers = n
	}
}

// WithUseGateway 设置是否使用网关模式
func WithUseGateway(useGateway bool) Option {
	return func(opt *Options) {
		opt.UseGateway = useGateway
	}
}
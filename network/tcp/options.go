package tcp

import (
	"io/github/gforgame/codec"
	"io/github/gforgame/network"
)

type Options struct {
	Name         string // 服务器名称
	ServiceAddr  string // current server service address (RPC)
	MessageCodec codec.MessageCodec
	IoDispatch   network.IoDispatch
	modules      []network.Module
	Router       *network.MessageRoute
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

// WithRouter 消息路由器
func WithRouter(r *network.MessageRoute) Option {
	return func(opt *Options) {
		opt.Router = r
	}
}

// WithModules 注册消息路由
func WithModules(ms ...network.Module) Option {
	return func(opt *Options) {
		opt.modules = append(opt.modules, ms...)
	}
}

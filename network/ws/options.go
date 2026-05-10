package ws

import (
	"github.com/forfun/gforgame/codec"
	"github.com/forfun/gforgame/network"
)

type Options struct {
	Name         string // 服务器名称
	ServiceAddr  string // current server service address (RPC)
	MessageCodec codec.MessageCodec
	IoDispatch   network.IoDispatch
	wsPath       string
	modules      []network.Module
	Router       *network.MessageRoute
	payloadMode  network.PayloadMode
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

// WithModules 注册消息路由
func WithModules(ms ...network.Module) Option {
	return func(opt *Options) {
		opt.modules = append(opt.modules, ms...)
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
		opt.payloadMode = mode
	}
}

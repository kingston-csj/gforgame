package network

import "io/github/gforgame/codec"

type Options struct {
	Name         string // 服务器名称
	ServiceAddr  string // current server service address (RPC)
	MessageCodec codec.MessageCodec
	IoDispatch   *BaseIoDispatch
	isWebsocket  bool
	wsPath       string
}

type Option func(*Options)

// WithAddress 指定ip地址
func WithAddress(addr string) Option {
	return func(opt *Options) {
		opt.ServiceAddr = addr
	}
}

// WithIoDispatch 消息处理链
func WithIoDispatch(dispatch *BaseIoDispatch) Option {
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

// WithWebsocket 设置为websocket
func WithWebsocket() Option {
	return func(opt *Options) {
		opt.isWebsocket = true
	}
}

// WithWsPath 设置websocket的路径
func WithWsPath(path string) Option {
	return func(opt *Options) {
		opt.wsPath = path
	}
}

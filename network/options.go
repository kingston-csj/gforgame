package network

import "io/github/gforgame/codec"

type Options struct {
	Name         string // 服务器名称
	ServiceAddr  string // current server service address (RPC)
	MessageCodec codec.MessageCodec
	IoDispatch   *BaseIoDispatch
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

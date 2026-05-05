package main

import (
	"context"
	"net"

	"github.com/forfun/gforgame/codec/json"
	playerdomain "github.com/forfun/gforgame/examples/domain/player"
	"github.com/forfun/gforgame/network/rpc"

	"google.golang.org/grpc"
)

type server struct {
	rpc.UnimplementedRpcServer
}

func (s *server) EnterRemote(ctx context.Context, in *rpc.PlayerCrossRequest) (*rpc.PlayerCrossReply, error) {
	codec := json.NewSerializer()
	p := playerdomain.Player{}
	codec.Decode(in.Data, &p)
	return &rpc.PlayerCrossReply{Message: "name: " + p.Name}, nil
}

func NewRpcServer(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	rpc.RegisterRpcServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

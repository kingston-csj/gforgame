package rpc

import (
	"fmt"
	"io/github/gforgame/log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	mu      sync.RWMutex
	clients map[uint32]RpcClient
)

func init() {
	clients = make(map[uint32]RpcClient)
}

func GetOrCreateClient(sid uint32) (RpcClient, error) {
	mu.RLock()
	s, found := clients[sid]
	mu.RUnlock()
	if !found {
		mu.Lock()
		defer mu.Unlock()
		s, found = clients[sid]
		// 双重检测
		if !found {
			conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return nil, fmt.Errorf("new rpc client failed %v", err)
			}

			c := NewRpcClient(conn)
			log.Info(fmt.Sprintf("connect to rpc server %d)", sid))
			clients[sid] = c
			return c, nil
		}
	}

	return s, nil
}

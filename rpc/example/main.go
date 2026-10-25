// Run with: go run ./rpc/example
package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	jsoncodec "github.com/forfun/gforgame/codec/json"
	"github.com/forfun/gforgame/network/dispatch"
	"github.com/forfun/gforgame/network/protocol"
	"github.com/forfun/gforgame/network/session"
	"github.com/forfun/gforgame/rpc"
)

type ReqExample struct{ Name string }
type ResExample struct{ Greeting string }
type Handler struct{}

func (*Handler) RPC_Handle(req *ReqExample) *ResExample {
	return &ResExample{Greeting: "Hello, " + req.Name}
}

type callback struct{ done chan struct{} }

func (c *callback) OnSuccess(response any) {
	fmt.Println("async:", response.(*ResExample).Greeting)
	close(c.done)
}
func (c *callback) OnFail(err error) {
	fmt.Println("async error:", err)
	close(c.done)
}

func main() {
	// Each peer registers these commands during startup, before reading messages.
	factory := protocol.NewMessageFactory()
	factory.RegisterMessage(900001, &ReqExample{})
	factory.RegisterMessage(900002, &ResExample{})
	factory.RegisterMessage(900003, &rpc.ErrorResponse{})
	client := rpc.NewClient(3*time.Second, factory)
	defer client.Close()
	if err := client.RegisterResponse(&ReqExample{}, &ResExample{}); err != nil {
		log.Fatal(err)
	}
	server := rpc.NewServer(factory)
	if err := server.Register(&Handler{}); err != nil {
		log.Fatal(err)
	}

	// net.Pipe keeps this example self contained. For TCP use net.Dial/net.Accept;
	// session, codec, dispatch and RPC setup stay the same.
	a, b := net.Pipe()
	left := session.NewSession(a, jsoncodec.NewSerializer(), factory)
	right := session.NewSession(b, jsoncodec.NewSerializer(), factory)
	clientDispatch := &dispatch.BaseIoDispatch{}
	clientDispatch.AddHandler(client)
	serverDispatch := &dispatch.BaseIoDispatch{}
	serverDispatch.AddHandler(server)
	opts := &dispatch.SessionLoopOpts{IdleCheckInterval: time.Minute, IdleTimeout: 10 * time.Minute}
	var loops sync.WaitGroup
	loops.Add(2)
	go func() { defer loops.Done(); dispatch.ServeSession(left, clientDispatch, opts) }()
	go func() { defer loops.Done(); dispatch.ServeSession(right, serverDispatch, opts) }()
	defer func() { left.Close(); right.Close(); loops.Wait() }()
	response, err := client.Request(left, &ReqExample{Name: "RPC"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("sync:", response.(*ResExample).Greeting)
	cb := &callback{done: make(chan struct{})}
	if err := client.Callback(left, &ReqExample{Name: "Callback"}, cb); err != nil {
		log.Fatal(err)
	}
	<-cb.done
}

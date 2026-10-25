package rpc

import (
	"errors"
	"sync"

	"github.com/forfun/gforgame/codec"
	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/dispatch"
	"github.com/forfun/gforgame/network/protocol"
	"github.com/forfun/gforgame/network/tcp"
)

// TCPServer combines service registration, TCP transport and session dispatch.
// Register messages and services before Start. A stopped server cannot be reused.
type TCPServer struct {
	*Server
	node       *tcp.TcpServer
	dispatcher *tcpDispatch
	mu         sync.Mutex
	started    bool
	stopped    bool
}

type TCPServerOption func(*tcpDispatch)

// WithSessionHooks attaches application lifecycle callbacks without coupling the
// RPC transport to player registries or application services.
func WithSessionHooks(created, closed func(network.Session)) TCPServerOption {
	return func(d *tcpDispatch) {
		d.created, d.closed = created, closed
	}
}

func NewTCPServer(addr string, factory *protocol.MessageFactory, codec codec.MessageCodec, options ...TCPServerOption) *TCPServer {
	service := NewServer(factory)
	d := &tcpDispatch{sessions: make(map[network.Session]struct{})}
	d.AddHandler(service)
	for _, option := range options {
		option(d)
	}
	return &TCPServer{
		Server:     service,
		dispatcher: d,
		node: tcp.NewServer(
			tcp.WithAddress(addr),
			tcp.WithRouter(network.NewMessageRoute(factory)),
			tcp.WithCodec(codec),
			tcp.WithIoDispatch(d),
		),
	}
}

func (s *TCPServer) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.stopped {
		return errors.New("rpc: TCP server has stopped")
	}
	if s.started {
		return errors.New("rpc: TCP server has already started")
	}
	if err := s.node.Start(); err != nil {
		return err
	}
	s.started = true
	return nil
}

// Stop closes the listener and all accepted RPC sessions. It is idempotent.
func (s *TCPServer) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stopped = true
	s.node.Stop()
	s.dispatcher.stop()
}

type tcpDispatch struct {
	dispatch.BaseIoDispatch
	created, closed func(network.Session)
	mu              sync.Mutex
	sessions        map[network.Session]struct{}
	stopped         bool
}

func (d *tcpDispatch) OnSessionCreated(session network.Session) {
	d.mu.Lock()
	if d.stopped {
		d.mu.Unlock()
		session.Close()
		return
	}
	d.sessions[session] = struct{}{}
	d.mu.Unlock()
	if d.created != nil {
		d.created(session)
	}
}

func (d *tcpDispatch) OnSessionClosed(session network.Session) {
	d.mu.Lock()
	_, tracked := d.sessions[session]
	delete(d.sessions, session)
	d.mu.Unlock()
	if tracked && d.closed != nil {
		d.closed(session)
	}
}

func (d *tcpDispatch) stop() {
	d.mu.Lock()
	d.stopped = true
	sessions := make([]network.Session, 0, len(d.sessions))
	for session := range d.sessions {
		sessions = append(sessions, session)
	}
	d.mu.Unlock()
	for _, session := range sessions {
		session.Close()
	}
}

package tcp

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"log/slog"
	"net"
	"sync"

	serverconfig "github.com/forfun/gforgame/config"
	"github.com/forfun/gforgame/network"
	"github.com/forfun/gforgame/network/protocol"
)

type TcpServer struct {
	Options
	Name     string // 服务器名称
	Running  chan bool
	listener net.Listener
	stopOnce sync.Once
}

func NewServer(opts ...Option) *TcpServer {
	opt := Options{DispatchWorkers: 1}
	for _, option := range opts {
		option(&opt)
	}
	if opt.DispatchWorkers <= 0 {
		opt.DispatchWorkers = 1
	}

	s := &TcpServer{
		Options: opt,
		Running: make(chan bool),
	}

	return s
}

func (s *TcpServer) Start() error {
	if s.ServiceAddr == "" {
		return errors.New("service address cannot be empty")
	}

	listener, err := net.Listen("tcp", s.ServiceAddr)
	if err != nil {
		return err
	}
	s.listener = listener

	go func() {
		s.startListen()
	}()

	return nil
}

func (s *TcpServer) Addr() string {
	return s.ServiceAddr
}

// Enable current server accept connection
func (s *TcpServer) startListen() {
	if s.listener == nil {
		slog.Error("tcp listener is nil")
		return
	}

	defer func() {
		_ = s.listener.Close()
	}()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			// listener 被关闭时会返回错误，此时应退出 accept 循环
			if errors.Is(err, net.ErrClosed) {
				return
			}
			slog.Error(fmt.Sprintf("new tcp conn failed %v", err))
			return
		}
		go onClientConnected(s, conn)
	}
}

// 处理客户端连接，包括socket,websocket
func onClientConnected(node *TcpServer, conn net.Conn) {
	defer func() {
		slog.Info("客户端连接关闭", "remoteAddr", conn.RemoteAddr().String(), "localAddr", conn.LocalAddr().String())
		// 处理客户端网络断开
		s := network.GetSession(conn)
		node.IoDispatch.OnSessionClosed(s)
		network.UnregisterSession(conn)
		_ = conn.Close()
	}()

	ioSession := network.NewSession(conn, node.MessageCodec)
	network.RegisterSession(conn, ioSession)

	// session created hook
	node.IoDispatch.OnSessionCreated(ioSession)

	// 异步读写数据
	go ioSession.Read()
	go ioSession.Write()

	workerCount := node.DispatchWorkers
	if workerCount <= 1 {
		// 直连模式：单worker快速路径：避免额外队列、哈希与goroutine开销。
		for {
			select {
			case task := <-ioSession.AsynTasks:
				task()
			case ioFrame := <-ioSession.DataReceived:
				if ioFrame != nil {
					node.IoDispatch.OnMessageReceived(ioSession, ioFrame)
				}
			case <-ioSession.Die:
				return
			}
		}
	}

	workerQueues := make([]chan *protocol.RequestDataFrame, workerCount)
	var dispatchWg sync.WaitGroup
	dispatchWg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		workerQueues[i] = make(chan *protocol.RequestDataFrame, 256)
		go func(workerIdx int) {
			defer dispatchWg.Done()
			for {
				select {
				case ioFrame := <-workerQueues[workerIdx]:
					if ioFrame != nil {
						node.IoDispatch.OnMessageReceived(ioSession, ioFrame)
					}
				case <-ioSession.Die:
					return
				}
			}
		}(i)
	}
	sessionWorkerIdx := hashSessionWorkerIndex(conn.RemoteAddr().String(), workerCount)

	// 主循环处理异步任务、分发消息和连接生命周期。
	// 直连模式按 session 固定路由；网关模式优先按玩家维度路由。
	for {
		select {
		case task := <-ioSession.AsynTasks:
			task()
		case ioFrame := <-ioSession.DataReceived:
			if ioFrame == nil {
				continue
			}
			// 将任务总线的消息分发到对应的worker
			targetWorkerIdx := resolveWorkerIndex(ioFrame, sessionWorkerIdx, workerCount)
			select {
			case workerQueues[targetWorkerIdx] <- ioFrame:
			case <-ioSession.Die:
				dispatchWg.Wait()
				return
			}
		case <-ioSession.Die:
			dispatchWg.Wait()
			return
		}
	}
}

func hashSessionWorkerIndex(sessionKey string, workerCount int) int {
	if workerCount <= 1 {
		return 0
	}
	h := fnv.New32a()
	_, _ = h.Write([]byte(sessionKey))
	return int(h.Sum32() % uint32(workerCount))
}

// 提取玩家ID，用于路由到指定worker
// 目前的做法效率低，后续可以优化
func resolveWorkerIndex(ioFrame *protocol.RequestDataFrame, fallbackSessionIdx int, workerCount int) int {
	if workerCount <= 1 {
		return 0
	}
	if serverconfig.ServerConfig.UseGateMode {
		playerID := extractPlayerID(ioFrame)
		if playerID != "" {
			return hashSessionWorkerIndex(playerID, workerCount)
		}
	}
	return fallbackSessionIdx
}

func extractPlayerID(ioFrame *protocol.RequestDataFrame) string {
	if ioFrame == nil {
		return ""
	}
	if ioFrame.Header.Payload != "" {
		return ioFrame.Header.Payload
	}
	msg := ioFrame.Msg
	if msg == nil {
		return ""
	}
	// 优先通过协议接口获取，避免反射依赖业务结构体字段。
	if carrier, ok := msg.(playerIDCarrier); ok {
		if playerID := carrier.GetPlayerID(); playerID != "" {
			return playerID
		}
	}
	// 兜底处理：兼容原始 JSON 场景。
	if playerID := extractPlayerIDFromJSON(msg); playerID != "" {
		return playerID
	}
	return ""
}

type playerIDCarrier interface{ GetPlayerID() string }

func extractPlayerIDFromJSON(msg any) string {
	var raw []byte
	switch v := msg.(type) {
	case []byte:
		raw = v
	case string:
		raw = []byte(v)
	default:
		return ""
	}
	if len(raw) == 0 {
		return ""
	}
	var payload struct {
		PlayerId string `json:"playerId"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return ""
	}
	return payload.PlayerId
}

func (n *TcpServer) Stop() {
	n.stopOnce.Do(func() {
		if n.listener != nil {
			_ = n.listener.Close()
		}
		network.CloseAllSessions()
	})
}

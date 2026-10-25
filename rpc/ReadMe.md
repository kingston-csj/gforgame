# RPC

基于 `network.Session` 的请求响应层，不依赖 `network/client`。

需要独立的 TCP RPC 监听时，可直接使用 `TCPServer`，统一管理注册、分发和连接生命周期：

```go
// factory 先注册 rpc.ErrorResponse 等公共协议；codec 由应用指定。
// Register 从方法签名发现协议，自动注册实现 rpc.CommandMessage 的请求和应答。
node := rpc.NewTCPServer(addr, factory, codec,
    rpc.WithSessionHooks(onSessionCreated, onSessionClosed)) // hooks 可省略
if err := node.Register(service); err != nil {
    return err
}
if err := node.Start(); err != nil {
    return err
}
defer node.Stop()
```

同一个节点可多次 `Register` 注册不同 service，应在启动前完成。
`Stop` 同时关闭监听和已接入的 RPC 连接，可重复调用；停止后不能再次启动。
地址是否为空、是否启用 RPC 由应用配置层决定；通用 RPC 层不读取游戏配置。
游戏服的网关协议与 Actor 适配在 `cmd/game/gate_rpc.go`，启动入口只负责组装和注册。

运行完整示例：

```sh
go run ./rpc/example
go test ./rpc/...
```

示例通过 `net.Pipe` 连接两个真实 network session，走 JSON 编解码、二进制协议和 dispatch，打印同步、异步调用结果。接入 TCP 时使用 `net.Dial` / `net.Accept` 创建连接，再按示例创建 session 和启动分发循环。

## 定义服务

```go
type ReqExample struct { Name string }
type ResExample struct { Greeting string }
type Handler struct{}

func (*Handler) RPC_Handle(req *ReqExample) *ResExample {
    return &ResExample{Greeting: "Hello, " + req.Name}
}

// 也支持：
// func (*Handler) RPC_Handle(req *ReqExample) (*ResExample, error)
// 可选接收当前连接：
// func (*Handler) RPC_Handle(session network.Session, req *ReqExample) (*ResExample, error)
```

服务方法必须以 `RPC_` 为前缀。没有此前缀的方法不参与注册；带此前缀但签名不符合以上两种形式的方法会导致注册失败，不会安装该次注册的任何方法。

启动时，两端使用相同的消息命令，命令由应用分配，RPC 不占用硬编码命令。
请求和响应可实现 `rpc.CommandMessage`，由服务端 `Register` 和客户端
`RegisterResponse` 自动注册（已手动注册的类型也继续支持）：

```go
func (*ReqExample) RPCCommand() int32 { return 900001 }
func (*ResExample) RPCCommand() int32 { return 900002 }

factory := protocol.NewMessageFactory()
factory.RegisterMessage(900003, &rpc.ErrorResponse{})
```

消息注册必须在收发消息之前完成，并将同一工厂传给该端的 Session、RPC Server 和 RPC Client，例如 `session.NewSession(conn, codec, factory)`。两端各自创建工厂并使用一致的命令映射。第一版错误响应是普通 Go 结构体，示例使用 `codec/json`；纯 protobuf 编解码器不能直接编码该结构体，需要另行适配错误消息。

## 注册和调用

```go
server := rpc.NewServer(factory)
if err := server.Register(&Handler{}); err != nil {
    return err
}

client := rpc.NewClient(5 * time.Second, factory) // <= 0 使用默认 5 秒
// 应在该 client 的全部会话结束后关闭，唤醒尚未完成的调用。
defer client.Close()
if err := client.RegisterResponse(&ReqExample{}, &ResExample{}); err != nil {
    return err
}

response, err := client.Request(session, &ReqExample{Name: "RPC"})
if err != nil {
    return err
}
fmt.Println(response.(*ResExample).Greeting)

// callback 实现 OnSuccess(any) / OnFail(error)。
if err := client.Callback(session, &ReqExample{Name: "RPC"}, callback); err != nil {
    return err
}
```

客户端分发管线首先添加 `client`，服务端添加 `server`。它们都实现 `dispatch.MessageHandler`，未知命令会继续向后分发：

```go
clientDispatch := &dispatch.BaseIoDispatch{}
clientDispatch.AddHandler(client)

serverDispatch := &dispatch.BaseIoDispatch{}
serverDispatch.AddHandler(server)
```

按示例使用 `dispatch.ServeSession` 启动 session 的读写和消息循环。同一个 session 只启动一套循环，不要重复读取 `DataReceivedChan`。使用现有循环时，将 RPC handler 接入现有 pipeline 即可。

## 行为约定

- 一个 session 使用一个 Client 管理全部 RPC 请求；一个 Client 可以管理多条 session。请求编号从 1 开始，0 留给普通消息。
- 通过请求编号、session 的关闭通道身份和响应 Cmd 关联响应，响应乱序不影响匹配。请求和响应须使用不同的命令；同一 client 中不能把响应命令再作为请求命令。
- 先注册 pending 再发送，响应、超时、断线、Client.Close 只有一个能完成调用；完成后删除 pending。迟到或重复响应不会进入业务处理器。
- `Callback` 参数校验或发送失败时直接返回错误，不调用 callback；返回 nil 后，最终结果在独立 goroutine 中回调一次。回调需要自行保证并发安全。
- 超时可以通过 `errors.Is(err, rpc.ErrTimeout)` 判断；断线使用 `rpc.ErrSessionClosed`；client 关闭使用 `rpc.ErrClientClosed`；远端错误使用 `errors.As(err, &remote)`，其中 remote 为 `*rpc.ErrorResponse`。
- handler 返回 error 时发送统一错误响应；返回 nil 响应时也报错。panic 在服务端记录堆栈，客户端只收到通用错误。发送响应失败时关闭 session，令调用方通过断线或超时结束等待。
- handler 在当前分发 goroutine 中串行执行，不要在其中同步等待同一 session 的另一次 RPC，否则响应无法得到分发。耗时任务应由应用安排执行方式。
- 定时器在发送前启动，但 `Session.Send` 本身可能因发送队列满而阻塞；当前接口无法取消正在阻塞的 Send。超时不代表服务端任务取消，也不提供自动重试。
- 注册服务和响应映射应在启动阶段完成。未注册的远端请求会交给后续业务管线，没有处理器时调用方将超时。

package contract

// GateLoginRequest 定义网关识别登录请求所需的最小字段集合。
type GateLoginRequest interface {
	GetPlayerID() string
	GetServerID() int32
}

// GateTransferMessage 定义网关内部转发包的最小契约。
type GateTransferMessage interface {
	GetPlayerID() string
	GetTransferCmd() int32
	GetTransferIndex() int32
	GetTransferBody() []byte
}

// ClientLoginAdapter 抽象网关接入登录协议所需的最小能力。
// 这部分通常由具体游戏项目提供实现。
type ClientLoginAdapter interface {
	LoginCmd() int32
	DecodeLoginRequest(body []byte) (GateLoginRequest, error)
	NewReplacingLoginPush() any
}

// TransferCodec 抽象网关与 logic 之间的转发协议。
// 跨语言协作时，Java 逻辑服只需要对齐这里对应的协议格式即可。
type TransferCodec interface {
	TransferCmd() int32
	NewTransferMessage(playerID string, cmd int32, index int32, body []byte) any
	ParseTransferMessage(msg any) (GateTransferMessage, error)
}

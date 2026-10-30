package serverconfig

// NodeType 节点类型，用 int 便于路由匹配和协议分发
type NodeType int

const (
	NodeTypeUnknown NodeType = 0 // 未知节点类型
	NodeTypeGame    NodeType = 1 // 游戏逻辑服
	NodeTypeBattle  NodeType = 2 // 战斗服
	NodeTypeCenter  NodeType = 8 // 网关服
	NodeTypeGateway NodeType = 9 // 游戏逻辑服
)

// String 方便日志打印
func (t NodeType) String() string {
	switch t {
	case NodeTypeGateway:
		return "gateway"
	case NodeTypeGame:
		return "game"
	case NodeTypeBattle:
		return "battle"
	case NodeTypeCenter:
		return "center"
	default:
		return "unknown"
	}
}

// ClusterConfig 集群全局配置
type ClusterConfig struct {
	Name      string            `json:"name" mapstructure:"name"`
	Discovery string            `json:"discovery" mapstructure:"discovery"` // 例如local, etcd, consul等
	Meta      map[string]string `json:"meta" mapstructure:"meta"`
}

// NodeConfig 单个节点配置
type NodeConfig struct {
	NodeId      string         `json:"nodeId" mapstructure:"nodeId"`
	ServerId    int32          `json:"serverId" mapstructure:"serverId"`
	Type        NodeType       `json:"type" mapstructure:"type"`
	UseGateMode bool           `json:"useGateMode" mapstructure:"useGateMode"` // 是否使用网关模式
	RpcAddr     string         `json:"rpcAddr" mapstructure:"rpcAddr"`         // RPC 地址
	ServiceAddr string         `json:"serviceAddr" mapstructure:"serviceAddr"` // 服务地址，例如socket/websocket,接收客户端协议
	HttpAddr    string         `json:"httpAddr" mapstructure:"httpAddr"`       // HTTP 地址
	PprofAddr   string         `json:"pprofAddr" mapstructure:"pprofAddr"`     // pprof 地址
	Extra       map[string]any `json:"extra" mapstructure:"extra"`
}

// AppConfig 配置文件根结构
type AppConfig struct {
	Cluster  ClusterConfig `json:"cluster" mapstructure:"cluster"`   // 集群全局配置
	NodeList []NodeConfig  `json:"nodeList" mapstructure:"nodeList"` // 节点列表
	NodeId   string        `json:"-" mapstructure:"-"`               // 当前节点 ID，仅由启动参数设置
}

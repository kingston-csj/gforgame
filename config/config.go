package config

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/common/util/conv"
	"github.com/forfun/gforgame/common/util/pathutil"
	"github.com/spf13/viper"
)

var (
	appConfig   AppConfig
	appConfigMu sync.RWMutex
	configOnce  sync.Once
)

// GetAppConfig 返回进程唯一的完整配置。Init 后供各模块只读访问。
func GetAppConfig() *AppConfig {
	return &appConfig
}

// GetCurrentNode 返回单例 NodeList 中的当前节点；初始化前返回 nil。
func GetCurrentNode() *NodeConfig {
	appConfigMu.RLock()
	defer appConfigMu.RUnlock()
	for i := range appConfig.NodeList {
		if appConfig.NodeList[i].NodeId == appConfig.NodeId {
			return &appConfig.NodeList[i]
		}
	}
	return nil
}

type ServerNodeInfo struct {
	ServerId uint32
	Type     uint32
	Addr     string
}

// DynamicServerNode 用于运行时注册动态发现到的服务器节点。
type DynamicServerNode struct {
	Id          uint32
	Type        uint32
	UseGateMode bool
	Addr        string
	HttpAddr    string
	PprofAddr   string
}

// DiscoveryProvider 扩展集群发现方式，例如 etcd；实现方负责获取节点配置。
// 在 Init 前注册，local 模式始终只使用 config.yml 的 nodeList。
type DiscoveryProvider func(ClusterConfig) ([]NodeConfig, error)

var discoveryProviders = map[string]DiscoveryProvider{}

func RegisterDiscoveryProvider(name string, provider DiscoveryProvider) {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" || name == "local" || provider == nil {
		panic("服务发现扩展必须指定非 local 名称和有效 provider")
	}
	discoveryProviders[name] = provider
}

func loadConfig(path string) (AppConfig, error) {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return AppConfig{}, fmt.Errorf("读取 config.yml 失败: %w", err)
	}
	var app AppConfig
	if err := v.Unmarshal(&app); err != nil {
		return AppConfig{}, fmt.Errorf("解析 config.yml 失败: %w", err)
	}
	app.Cluster.Discovery = strings.ToLower(strings.TrimSpace(app.Cluster.Discovery))
	switch app.Cluster.Discovery {
	case "local":
	case "":
		return AppConfig{}, fmt.Errorf("config.yml 必须指定 cluster.discovery")
	default:
		provider, ok := discoveryProviders[app.Cluster.Discovery]
		if !ok {
			return AppConfig{}, fmt.Errorf("尚未注册 cluster.discovery=%s 的发现实现（etcd 等扩展需在 Init 前注册）", app.Cluster.Discovery)
		}
		nodes, err := provider(app.Cluster)
		if err != nil {
			return AppConfig{}, fmt.Errorf("发现集群节点失败: %w", err)
		}
		app.NodeList = nodes
	}
	if len(app.NodeList) == 0 {
		return AppConfig{}, fmt.Errorf("nodeList 为空")
	}
	nodeIDs := make(map[string]bool, len(app.NodeList))
	serverIDs := make(map[int32]bool, len(app.NodeList))
	for _, node := range app.NodeList {
		if strings.TrimSpace(node.NodeId) == "" || node.ServerId <= 0 || node.Type <= NodeTypeUnknown {
			return AppConfig{}, fmt.Errorf("无效节点配置 nodeId=%q serverId=%d type=%d", node.NodeId, node.ServerId, node.Type)
		}
		if nodeIDs[node.NodeId] || serverIDs[node.ServerId] {
			return AppConfig{}, fmt.Errorf("节点 nodeId 或 serverId 重复: %s/%d", node.NodeId, node.ServerId)
		}
		nodeIDs[node.NodeId], serverIDs[node.ServerId] = true, true
	}
	return app, nil
}

// Init 根据启动参数 nodeId 选择当前节点，只读取 config/config.yml。
func Init(nodeID string) {
	configOnce.Do(func() {
		path, ok := pathutil.ResolveExistingRelativeFile(filepath.Join("config", "config.yml"))
		if !ok {
			panic("config/config.yml not found")
		}
		app, err := loadConfig(path)
		if err != nil {
			panic(err)
		}
		if err := applyConfig(app, nodeID); err != nil {
			panic(err)
		}
		logger.Info(fmt.Sprintf("服务配置加载完成 discovery=%s nodeId=%s serverId=%d serverAddr=%s httpAddr=%s",
			appConfig.Cluster.Discovery, appConfig.NodeId, GetCurrentNode().ServerId, GetCurrentNode().ServiceAddr, GetCurrentNode().HttpAddr))
	})
}

func applyConfig(app AppConfig, nodeID string) error {
	var current *NodeConfig
	for i := range app.NodeList {
		node := &app.NodeList[i]
		if node.NodeId == nodeID {
			current = node
		}
	}
	if current == nil {
		return fmt.Errorf("未找到 nodeId=%q 对应的节点，请检查 config.yml 的 nodeList 或 -nodeId 参数", nodeID)
	}
	app.NodeId = nodeID
	appConfigMu.Lock()
	appConfig = app
	appConfigMu.Unlock()
	return nil
}

// --- Extra 工具方法 ---

func GetExtraString(key string) (string, bool) {
	val, ok := getCurrentExtraValue(key)
	if !ok || val == nil {
		return "", false
	}
	return conv.StringValue(val), true
}

func GetExtraInt(key string) (int, bool) {
	val, ok := getCurrentExtraValue(key)
	if !ok || val == nil {
		return 0, false
	}
	return conv.IntValue(val), true
}

func GetExtraBool(key string) (bool, bool) {
	val, ok := getCurrentExtraValue(key)
	if !ok || val == nil {
		return false, false
	}
	return conv.BooleanValue(val), true
}

func getCurrentExtraValue(key string) (any, bool) {
	node := GetCurrentNode()
	if node == nil {
		return nil, false
	}
	return getExtraValueByPath(node.Extra, key)
}

func getExtraStringFromMap(extra map[string]any, key string) string {
	val, ok := getExtraValueByPath(extra, key)
	if !ok || val == nil {
		return ""
	}
	switch v := val.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		return ""
	}
}

func getExtraValueByPath(extra map[string]any, key string) (any, bool) {
	if extra == nil || key == "" {
		return nil, false
	}
	parts := strings.Split(key, ".")
	var current any = extra
	for _, part := range parts {
		m, ok := normalizeAnyMap(current)
		if !ok {
			return nil, false
		}
		next, ok := m[part]
		if !ok {
			return nil, false
		}
		current = next
	}
	return current, true
}

func normalizeAnyMap(v any) (map[string]any, bool) {
	switch m := v.(type) {
	case map[string]any:
		return m, true
	case map[any]any:
		result := make(map[string]any, len(m))
		for k, val := range m {
			ks, ok := k.(string)
			if !ok {
				return nil, false
			}
			result[ks] = val
		}
		return result, true
	default:
		return nil, false
	}
}

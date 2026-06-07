package config

import (
	"embed"
	"fmt"

	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/common/util/conv"
	"github.com/forfun/gforgame/common/util/pathutil"
	"github.com/spf13/viper"
)

type Config struct {
	// 服务器id
	ServerId uint32
	// 服务地址(socket/websocket)
	ServerUrl string
	// 是否使用网关模式（true: 客户端->gate->logic，false: 客户端直连logic）
	UseGateMode bool
	//后端管理地址
	HttpUrl string
	// pprof性能监测地址
	PprofAddr string
	// 未匹配到强类型字段的配置，按需手动读取
	Extra map[string]any
}

var configFS embed.FS

var (
	ServerConfig Config
	allServers   []serverNode
	allServersMu sync.RWMutex
	dynamicIDs   = make(map[uint32]struct{})
)

type serverNode struct {
	Id          uint32 `mapstructure:"id"`
	Type        uint32 `mapstructure:"type"`
	UseGateMode bool   `mapstructure:"useGateMode"`
	Addr        string `mapstructure:"addr"`
	HttpAddr    string `mapstructure:"httpAddr"`
	PprofAddr   string `mapstructure:"pprofAddr"`
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

// roleConfig 用于承接角色级配置（config-gate/config-game）。
// 已定义字段自动注入，未定义字段进入 Extra。
type roleConfig struct {
	Server struct {
		Id uint32 `mapstructure:"id"`
	} `mapstructure:"server"`
	Extra map[string]any `mapstructure:",remain"`
}

// 配置读取规则：
// 加载顺序: default -> 角色配置(config-gate/config-logic)
// 先加载的配置会被后加载的同名配置所替换！！！！
// 1.优先读default.yml文件，应用程序内部配置，项目打包成二进制后不嵌入该配置
// 2.根据ENV读取 config-{ENV}.yml（如 gate/logic）
func init() {
	// 创建 Viper 实例
	v := viper.New()
	v.SetConfigType("yml")

	// 默认配置优先读取外部文件，便于打包后与二进制一起部署。
	if defaultConfigFile := resolveDefaultConfigPath(); defaultConfigFile != "" {
		v.SetConfigFile(defaultConfigFile)
		if err := v.ReadInConfig(); err != nil {
			panic(fmt.Errorf("failed to read config file %s: %w", defaultConfigFile, err))
		}
		logger.Info(fmt.Sprintf("已加载默认配置: %s", defaultConfigFile))
	} else {
		// 外部 default.yml 不存在时，回退到二进制内嵌默认配置。
		f, err := configFS.Open("default.yml")
		if err != nil {
			panic(fmt.Errorf("failed to open embedded config file: %w", err))
		}
		defer f.Close()
		if err := v.ReadConfig(f); err != nil {
			panic(fmt.Errorf("failed to read embedded config: %w", err))
		}
		logger.Info("已加载内嵌默认配置")
	}

	// 允许 Viper 读取环境变量
	v.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)

	// 获取环境变量，确定要加载的配置文件
	env := os.Getenv("ENV")
	if env == "" {
		env = "game"
	}
	v.SetConfigName("config-" + env)
	if configFile := resolveEnvConfigPath(env); configFile != "" {
		v.SetConfigFile(configFile)
	}
	// 再次读取配置文件，这次是根据环境变量，使用合并配置的方法确保旧配置被替换
	if err := v.MergeInConfig(); err != nil {
		logger.ErrorNoStack(fmt.Sprintf("加载环境配置失败，继续使用默认配置: %v", err))
	} else {
		logger.Info(fmt.Sprintf("已加载环境配置: %s", v.ConfigFileUsed()))
	}
	if err := v.UnmarshalKey("servers", &allServers); err != nil {
		panic(fmt.Errorf("解析 servers 节点失败: %w", err))
	}
	var rc roleConfig
	if err := v.Unmarshal(&rc); err != nil {
		logger.ErrorNoStack(fmt.Sprintf("解析角色配置失败，继续使用兜底逻辑: %v", err))
	}

	serverID := rc.Server.Id
	if serverID == 0 {
		serverID = resolveServerID()
	}
	if serverID == 0 {
		panic(fmt.Errorf("配置项 server.id 为空，请在 config-%s.yml 或环境变量 SERVER_ID 中提供", env))
	}
	currentNode := resolveCurrentServerNode(serverID)
	if currentNode == nil {
		panic(fmt.Errorf("未找到 server.id=%d 对应的节点配置，请检查 default.yml 的 servers 列表", serverID))
	}

	ServerConfig = Config{
		ServerId:    currentNode.Id,
		ServerUrl:   currentNode.Addr,
		UseGateMode: currentNode.UseGateMode,
		HttpUrl:     currentNode.HttpAddr,
		PprofAddr:   currentNode.PprofAddr,
		Extra:       rc.Extra,
	}

	logger.Info(fmt.Sprintf(
		"服务配置加载完成 env=%s type=%d serverId=%d serverAddr=%s httpAddr=%s",
		env,
		currentNode.Type,
		ServerConfig.ServerId,
		ServerConfig.ServerUrl,
		ServerConfig.HttpUrl,
	))
}

func resolveCurrentServerNode(serverID uint32) *serverNode {
	allServersMu.RLock()
	defer allServersMu.RUnlock()
	for i := range allServers {
		if allServers[i].Id == serverID {
			node := allServers[i]
			return &node
		}
	}
	return nil
}

func resolveServerID() uint32 {
	if raw := strings.TrimSpace(os.Getenv("SERVER_ID")); raw != "" {
		var id uint32
		_, err := fmt.Sscanf(raw, "%d", &id)
		if err == nil {
			return id
		}
	}
	return 0
}

func resolveEnvConfigPath(env string) string {
	fileName := fmt.Sprintf("config-%s.yml", env)
	if path, ok := pathutil.ResolveExistingRelativeFile(filepath.Join("config", fileName)); ok {
		return path
	}
	return ""
}

func resolveDefaultConfigPath() string {
	const fileName = "default.yml"
	if path, ok := pathutil.ResolveExistingRelativeFile(filepath.Join("config", fileName)); ok {
		return path
	}
	return ""
}

func GetServerAddrByType(serverType uint32) string {
	allServersMu.RLock()
	defer allServersMu.RUnlock()
	for i := range allServers {
		if allServers[i].Type == serverType {
			return allServers[i].Addr
		}
	}
	return ""
}

func GetServerTypeByID(serverID uint32) (uint32, bool) {
	allServersMu.RLock()
	defer allServersMu.RUnlock()
	for i := range allServers {
		if allServers[i].Id == serverID {
			return allServers[i].Type, true
		}
	}
	return 0, false
}

func GetServersByType(serverType uint32) []ServerNodeInfo {
	allServersMu.RLock()
	defer allServersMu.RUnlock()
	result := make([]ServerNodeInfo, 0)
	for i := range allServers {
		if allServers[i].Type != serverType {
			continue
		}
		result = append(result, ServerNodeInfo{
			ServerId: allServers[i].Id,
			Type:     allServers[i].Type,
			Addr:     allServers[i].Addr,
		})
	}
	return result
}

// RegisterDynamicServers 运行时注册或更新服务器节点。
// 同 id 节点会被覆盖，不存在的节点会被追加。
func RegisterDynamicServers(nodes []DynamicServerNode) {
	if len(nodes) == 0 {
		return
	}

	allServersMu.Lock()
	defer allServersMu.Unlock()

	indexByID := make(map[uint32]int, len(allServers))
	for i := range allServers {
		indexByID[allServers[i].Id] = i
	}

	for _, node := range nodes {
		server := serverNode{
			Id:          node.Id,
			Type:        node.Type,
			UseGateMode: node.UseGateMode,
			Addr:        node.Addr,
			HttpAddr:    node.HttpAddr,
			PprofAddr:   node.PprofAddr,
		}
		if idx, ok := indexByID[node.Id]; ok {
			allServers[idx] = server
			continue
		}
		allServers = append(allServers, server)
		indexByID[node.Id] = len(allServers) - 1
	}
}

// SyncDynamicServers 使用最新列表全量同步动态节点。
// 旧的动态节点会被移除，静态配置节点保持不变。
func SyncDynamicServers(nodes []DynamicServerNode) {
	allServersMu.Lock()
	defer allServersMu.Unlock()

	filtered := make([]serverNode, 0, len(allServers))
	for _, server := range allServers {
		if _, ok := dynamicIDs[server.Id]; ok {
			continue
		}
		filtered = append(filtered, server)
	}
	allServers = filtered
	dynamicIDs = make(map[uint32]struct{}, len(nodes))

	indexByID := make(map[uint32]int, len(allServers))
	for i := range allServers {
		indexByID[allServers[i].Id] = i
	}

	for _, node := range nodes {
		server := serverNode{
			Id:          node.Id,
			Type:        node.Type,
			UseGateMode: node.UseGateMode,
			Addr:        node.Addr,
			HttpAddr:    node.HttpAddr,
			PprofAddr:   node.PprofAddr,
		}
		if idx, ok := indexByID[node.Id]; ok {
			allServers[idx] = server
		} else {
			allServers = append(allServers, server)
			indexByID[node.Id] = len(allServers) - 1
		}
		dynamicIDs[node.Id] = struct{}{}
	}
}

func GetExtraString(key string) (string, bool) {
	val, ok := getExtraValueByPath(ServerConfig.Extra, key)
	if !ok || val == nil {
		return "", false
	}
	return conv.StringValue(val), true
}

func GetExtraInt(key string) (int, bool) {
	val, ok := getExtraValueByPath(ServerConfig.Extra, key)
	if !ok || val == nil {
		return 0, false
	}
	return conv.IntValue(val), true
}

func GetExtraBool(key string) (bool, bool) {
	val, ok := getExtraValueByPath(ServerConfig.Extra, key)
	if !ok || val == nil {
		return false, false
	}
	return conv.BooleanValue(val), true
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

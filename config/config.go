package config

import (
	"embed"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

//go:embed default.yml
var configFS embed.FS

var (
	ServerConfig Config
	allServers   []serverNode
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

// roleConfig 用于承接角色级配置（config-gate/config-logic）。
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
// 1.优先读default.yml文件，应用程序内部配置，项目打包成二进制可执行文件也会嵌入该配置
// 2.根据ENV读取 config-{ENV}.yml（如 gate/logic）
func init() {
	// 创建 Viper 实例
	v := viper.New()
	v.SetConfigType("yml")
	// 打包后的二进制文件也要
	f, err := configFS.Open("default.yml")
	if err != nil {
		panic(fmt.Errorf("failed to open config file: %w", err))
	}
	defer f.Close()
	// 从 io.Reader 读取配置
	if err := v.ReadConfig(f); err != nil {
		panic(fmt.Errorf("failed to read config: %w", err))
	}

	// 允许 Viper 读取环境变量
	v.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)

	// 获取环境变量，确定要加载的配置文件
	env := os.Getenv("ENV")
	if env == "" {
		env = "logic"
	}
	v.SetConfigName("config-" + env)
	if configFile := resolveEnvConfigPath(env); configFile != "" {
		v.SetConfigFile(configFile)
	}
	// 再次读取配置文件，这次是根据环境变量，使用合并配置的方法确保旧配置被替换
	if err := v.MergeInConfig(); err != nil {
		slog.Error("加载环境配置失败，继续使用默认配置", "env", env, "err", err)
	} else {
		slog.Info("已加载环境配置", "env", env, "configFile", v.ConfigFileUsed())
	}
	if err := v.UnmarshalKey("servers", &allServers); err != nil {
		panic(fmt.Errorf("解析 servers 节点失败: %w", err))
	}
	var rc roleConfig
	if err := v.Unmarshal(&rc); err != nil {
		slog.Error("解析角色配置失败，继续使用兜底逻辑", "err", err)
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

	slog.Info("服务配置加载完成",
		"env", env,
		"type", currentNode.Type,
		"serverId", ServerConfig.ServerId,
		"serverAddr", ServerConfig.ServerUrl,
		"httpAddr", ServerConfig.HttpUrl,
	)
}

func resolveCurrentServerNode(serverID uint32) *serverNode {
	for i := range allServers {
		if allServers[i].Id == serverID {
			return &allServers[i]
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
	if exePath, err := os.Executable(); err == nil {
		if path, ok := findConfigFileFromBase(filepath.Dir(exePath), fileName); ok {
			return path
		}
	}
	if cwd, err := os.Getwd(); err == nil {
		if path, ok := findConfigFileFromBase(cwd, fileName); ok {
			return path
		}
	}
	return ""
}

func findConfigFileFromBase(baseDir, fileName string) (string, bool) {
	dir := filepath.Clean(baseDir)
	for {
		candidate := filepath.Join(dir, "config", fileName)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", false
}

func GetServerAddrByType(serverType uint32) string {
	for i := range allServers {
		if allServers[i].Type == serverType {
			return allServers[i].Addr
		}
	}
	return ""
}

func GetServerTypeByID(serverID uint32) (uint32, bool) {
	for i := range allServers {
		if allServers[i].Id == serverID {
			return allServers[i].Type, true
		}
	}
	return 0, false
}

func GetServersByType(serverType uint32) []ServerNodeInfo {
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

func GetExtraString(key string) (string, bool) {
	val, ok := getExtraValueByPath(ServerConfig.Extra, key)
	if !ok || val == nil {
		return "", false
	}
	switch v := val.(type) {
	case string:
		return v, true
	case fmt.Stringer:
		return v.String(), true
	case int:
		return strconv.Itoa(v), true
	case int8:
		return strconv.FormatInt(int64(v), 10), true
	case int16:
		return strconv.FormatInt(int64(v), 10), true
	case int32:
		return strconv.FormatInt(int64(v), 10), true
	case int64:
		return strconv.FormatInt(v, 10), true
	case uint:
		return strconv.FormatUint(uint64(v), 10), true
	case uint8:
		return strconv.FormatUint(uint64(v), 10), true
	case uint16:
		return strconv.FormatUint(uint64(v), 10), true
	case uint32:
		return strconv.FormatUint(uint64(v), 10), true
	case uint64:
		return strconv.FormatUint(v, 10), true
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), true
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), true
	case bool:
		return strconv.FormatBool(v), true
	default:
		return "", false
	}
}

func GetExtraInt(key string) (int, bool) {
	val, ok := getExtraValueByPath(ServerConfig.Extra, key)
	if !ok || val == nil {
		return 0, false
	}
	switch v := val.(type) {
	case int:
		return v, true
	case int8:
		return int(v), true
	case int16:
		return int(v), true
	case int32:
		return int(v), true
	case int64:
		return int(v), true
	case uint:
		return int(v), true
	case uint8:
		return int(v), true
	case uint16:
		return int(v), true
	case uint32:
		return int(v), true
	case uint64:
		return int(v), true
	case float32:
		return int(v), true
	case float64:
		return int(v), true
	case string:
		n, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return 0, false
		}
		return n, true
	default:
		return 0, false
	}
}

func GetExtraBool(key string) (bool, bool) {
	val, ok := getExtraValueByPath(ServerConfig.Extra, key)
	if !ok || val == nil {
		return false, false
	}
	switch v := val.(type) {
	case bool:
		return v, true
	case string:
		b, err := strconv.ParseBool(strings.TrimSpace(v))
		if err != nil {
			return false, false
		}
		return b, true
	case int:
		return v != 0, true
	case int8:
		return v != 0, true
	case int16:
		return v != 0, true
	case int32:
		return v != 0, true
	case int64:
		return v != 0, true
	case uint:
		return v != 0, true
	case uint8:
		return v != 0, true
	case uint16:
		return v != 0, true
	case uint32:
		return v != 0, true
	case uint64:
		return v != 0, true
	case float32:
		return v != 0, true
	case float64:
		return v != 0, true
	default:
		return false, false
	}
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

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/forfun/gforgame/common/logger"
	serverconfig "github.com/forfun/gforgame/config"
)

const serverDiscoveryRefreshInterval = 5 * time.Minute

// discoveredServerListResponse 对应服务发现接口的返回结构。
type discoveredServerListResponse struct {
	Code int                      `json:"code"`
	Msg  any                      `json:"msg"`
	Data discoveredServerListData `json:"data"`
}

type discoveredServerListData struct {
	Servers    []DiscoveredServer `json:"servers"`
	TotalCount int                `json:"totalCount"`
}

// DiscoveredServer 表示服务发现接口返回的单个服务器节点。
type DiscoveredServer struct {
	ID       int32   `json:"id"`
	Name     string  `json:"name"`
	IP       *string `json:"ip"`
	Port     int     `json:"port"`
	HTTPPort int     `json:"httpPort"`
	UseGate  int     `json:"useGate"`
}

// FetchServerList 使用 HTTP GET 请求服务发现接口并解析服务器列表。
func FetchServerList() ([]DiscoveredServer, error) {
	webURL, ok := serverconfig.GetExtraString("discovery.apiurl")
	if !ok {
		return nil, fmt.Errorf("discovery.apiurl is not set")
	}
	return FetchServerListFromURL(webURL)
}

// FetchServerListFromURL 从指定 URL 获取并解析服务器列表。
func FetchServerListFromURL(rawURL string) ([]DiscoveredServer, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(rawURL)
	if err != nil {
		return nil, fmt.Errorf("请求服务器列表失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return nil, fmt.Errorf("请求服务器列表失败: status=%d body=%s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取服务器列表响应失败: %w", err)
	}

	var result discoveredServerListResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析服务器列表响应失败: %w", err)
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("服务器列表接口返回失败: code=%d msg=%v", result.Code, result.Msg)
	}

	return filterValidDiscoveredServers(result.Data.Servers), nil
}

// RegisterDiscoveredServers 将发现到的节点转换为运行时配置节点并注册。
func RegisterDiscoveredServers(servers []DiscoveredServer) {
	serverconfig.SyncDynamicServers(convertDiscoveredServers(servers))
}

func convertDiscoveredServers(servers []DiscoveredServer) []serverconfig.DynamicServerNode {
	result := make([]serverconfig.DynamicServerNode, 0, len(servers))
	for _, server := range servers {
		if server.ID == 0 || server.IP == nil {
			continue
		}
		result = append(result, serverconfig.DynamicServerNode{
			Id:          uint32(server.ID),
			Type:        logicServerType,
			UseGateMode: server.UseGate != 0,
			Addr:        buildDiscoveredServerAddr(server),
			HttpAddr:    buildDiscoveredServerHTTPAddr(server),
		})
	}
	return result
}

func buildDiscoveredServerAddr(server DiscoveredServer) string {
	if server.IP == nil {
		return ""
	}
	if server.Port > 0 {
		return fmt.Sprintf("%s:%d", *server.IP, server.Port)
	}
	if server.HTTPPort > 0 {
		return fmt.Sprintf("%s:%d", *server.IP, server.HTTPPort)
	}
	return ""
}

func buildDiscoveredServerHTTPAddr(server DiscoveredServer) string {
	if server.IP == nil || server.HTTPPort <= 0 {
		return ""
	}
	return fmt.Sprintf("%s:%d", *server.IP, server.HTTPPort)
}

// syncDiscoveredServers 执行一次服务发现同步，并按最新结果收敛后端连接。
func syncDiscoveredServers() error {
	serverList, err := FetchServerList()
	if err != nil {
		return err
	}
	RegisterDiscoveredServers(serverList)
	syncDiscoveredBackendPools(serverList)
	logger.Info(fmt.Sprintf("discovered %d servers", len(serverList)))
	return nil
}

// startServerDiscoveryHeartbeat 启动服务发现心跳。
// 启动时会先同步一次，后续每 5 分钟同步一次。
func startServerDiscoveryHeartbeat() error {
	if err := syncDiscoveredServers(); err != nil {
		return err
	}
	go func() {
		ticker := time.NewTicker(serverDiscoveryRefreshInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := syncDiscoveredServers(); err != nil {
					logger.ErrorNoStack(fmt.Errorf("refresh discovered servers failed: %v", err))
				}
			case <-serverDiscoveryStop:
				return
			}
		}
	}()
	return nil
}

// startLocalServerDiscovery 仅基于本地配置同步一次后端节点。
func startLocalServerDiscovery() error {
	syncLocalBackendPools()
	logger.Info("onlyLocal=true，已按本地配置完成一次后端节点同步")
	return nil
}

// filterValidDiscoveredServers 过滤掉总计行或缺少关键字段的无效记录。
func filterValidDiscoveredServers(servers []DiscoveredServer) []DiscoveredServer {
	result := make([]DiscoveredServer, 0, len(servers))
	for _, server := range servers {
		if server.ID == 0 || server.IP == nil || server.UseGate == 0 {
			continue
		}
		result = append(result, server)
	}
	return result
}

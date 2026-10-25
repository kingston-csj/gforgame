package config

import "strconv"

func GetServerTypeByID(serverID uint32) (uint32, bool) {
	appConfigMu.RLock()
	defer appConfigMu.RUnlock()
	for _, node := range appConfig.NodeList {
		if uint32(node.ServerId) == serverID {
			return uint32(node.Type), true
		}
	}
	return 0, false
}

func GetServersByType(serverType uint32) []ServerNodeInfo {
	appConfigMu.RLock()
	defer appConfigMu.RUnlock()
	result := make([]ServerNodeInfo, 0)
	for _, node := range appConfig.NodeList {
		if uint32(node.Type) == serverType {
			result = append(result, ServerNodeInfo{
				ServerId: uint32(node.ServerId), Type: uint32(node.Type), Addr: node.ServiceAddr,
			})
		}
	}
	return result
}

// RegisterDynamicServers 将发现节点直接注册到单例 NodeList。
// 同 serverId 节点保留其 nodeId 和 extra，仅更新路由字段。
func RegisterDynamicServers(nodes []DynamicServerNode) {
	appConfigMu.Lock()
	defer appConfigMu.Unlock()
	for _, node := range nodes {
		index := -1
		for i := range appConfig.NodeList {
			if uint32(appConfig.NodeList[i].ServerId) == node.Id {
				index = i
				break
			}
		}
		if index == -1 {
			appConfig.NodeList = append(appConfig.NodeList, NodeConfig{NodeId: strconv.FormatUint(uint64(node.Id), 10)})
			index = len(appConfig.NodeList) - 1
		}
		current := &appConfig.NodeList[index]
		current.ServerId = int32(node.Id)
		current.Type = NodeType(node.Type)
		current.UseGateMode = node.UseGateMode
		current.ServiceAddr = node.Addr
		current.HttpAddr = node.HttpAddr
		current.PprofAddr = node.PprofAddr
	}
}

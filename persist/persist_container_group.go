package persist

import (
	"hash/fnv"
	"sync"
)

// PersistContainerGroup 持久化容器组（分片并发持久化）
// 子容器可以是任意 PersistContainer 实现（QueueContainer、DelayContainer 等），
// 由调用方在构造时决定每个分片的具体实现。
type PersistContainerGroup struct {
	group []PersistContainer
	name  string
}

// NewPersistContainerGroup 创建通用容器组，调用方自行提供每个分片的容器实现
// name: 名称前缀
// group: 子容器数组（任意 PersistContainer 实现）
func NewPersistContainerGroup(name string, group []PersistContainer) *PersistContainerGroup {
	return &PersistContainerGroup{
		group: group,
		name:  name,
	}
}

// NewQueueContainerGroup 便捷构造：使用 N 个 QueueContainer 作为分片
// name: 名称
// savingStrategy: 保存策略
// workers: 并发协程数（分片数）
func NewQueueContainerGroup(name string, savingStrategy SavingStrategy, workers int) *PersistContainerGroup {
	group := make([]PersistContainer, workers)
	for i := range group {
		group[i] = NewQueueContainer(name, savingStrategy)
	}
	return NewPersistContainerGroup(name, group)
}

// Receive 接收实体，根据 entity.GetId() 哈希取模路由到对应分片
func (g *PersistContainerGroup) Receive(entity Entity) {
	index := g.hashIndex(entity.GetId())
	g.group[index].Receive(entity)
}

// Size 返回所有分片总大小
func (g *PersistContainerGroup) Size() int {
	size := 0
	for _, container := range g.group {
		size += container.Size()
	}
	return size
}

// ShutdownGraceful 优雅关闭所有分片
func (g *PersistContainerGroup) ShutdownGraceful() {
	var wg sync.WaitGroup
	for _, container := range g.group {
		wg.Add(1)
		go func(c PersistContainer) {
			defer wg.Done()
			c.ShutdownGraceful()
		}(container)
	}
	wg.Wait()
}

// hashIndex 根据 key 哈希取模
func (g *PersistContainerGroup) hashIndex(key string) int {
	fnvHash := fnv.New32a()
	_, _ = fnvHash.Write([]byte(key))
	return int(fnvHash.Sum32()) % len(g.group)
}

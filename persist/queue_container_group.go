package persist

import (
	"hash/fnv"
	"sync"
)

// QueueContainerGroup 队列容器组（分片并发持久化）
type QueueContainerGroup struct {
	group []*QueueContainer // 容器数组
	name  string
}

// NewQueueContainerGroup 创建分片持久化容器组
// name: 名称
// savingStrategy: 保存策略
// workers: 并发协程数（分片数）
func NewQueueContainerGroup(name string, savingStrategy SavingStrategy, workers int) *QueueContainerGroup {
	group := make([]*QueueContainer, workers)

	for i := 0; i < workers; i++ {
		work := NewQueueContainer(name, savingStrategy)
		group[i] = work
	}

	return &QueueContainerGroup{
		group: group,
		name:  name + "-group",
	}
}

// Receive 接收实体，根据 entity.GetId() 哈希取模路由到对应分片
func (g *QueueContainerGroup) Receive(entity Entity) {
	index := g.hashIndex(entity.GetId())
	g.group[index].Receive(entity)
}

// Size 返回所有队列总大小
func (g *QueueContainerGroup) Size() int {
	size := 0
	for _, container := range g.group {
		size += container.Size()
	}
	return size
}

// ShutdownGraceful 优雅关闭所有容器
func (g *QueueContainerGroup) ShutdownGraceful() {
	var wg sync.WaitGroup
	for _, container := range g.group {
		wg.Add(1)
		go func(c *QueueContainer) {
			defer wg.Done()
			c.ShutdownGraceful()
		}(container)
	}
	wg.Wait()
}

// hashIndex 根据 key 哈希取模
func (g *QueueContainerGroup) hashIndex(key string) int {
	fnvHash := fnv.New32a()
	_, _ = fnvHash.Write([]byte(key))
	return int(fnvHash.Sum32()) % len(g.group)
}
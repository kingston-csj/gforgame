package trie

// NodeContainer 节点容器接口，定义子节点的操作
type NodeContainer interface {
	Get(c rune) *TrieNode       // 根据字符获取子节点
	Add(node *TrieNode)         // 添加子节点
	Remove(c rune) *TrieNode    // 删除子节点并返回被删除的节点
	Size() int                  // 获取子节点数量
	GetAll() []*TrieNode        // 获取所有子节点
	Transform() NodeContainer   // 将Map容器转换为Array容器
}

// MapNodeContainer 基于map的节点容器
type MapNodeContainer struct {
	nodes map[rune]*TrieNode
}

func NewMapNodeContainer() *MapNodeContainer {
	return &MapNodeContainer{
		nodes: make(map[rune]*TrieNode),
	}
}

func (m *MapNodeContainer) Get(c rune) *TrieNode {
	return m.nodes[c]
}

func (m *MapNodeContainer) Add(node *TrieNode) {
	m.nodes[node.val] = node
}

func (m *MapNodeContainer) Remove(c rune) *TrieNode {
	node := m.nodes[c]
	if node != nil {
		delete(m.nodes, c)
	}
	return node
}

func (m *MapNodeContainer) Size() int {
	return len(m.nodes)
}

func (m *MapNodeContainer) GetAll() []*TrieNode {
	nodes := make([]*TrieNode, 0, m.Size())
	for _, node := range m.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

func (m *MapNodeContainer) Transform() NodeContainer {
	arrayContainer := NewArrayNodeContainer()
	for _, node := range m.nodes {
		arrayContainer.Add(node)
	}
	return arrayContainer
}

// ArrayNodeContainer 基于数组的节点容器（用于子节点较少时节省内存）
type ArrayNodeContainer struct {
	nodes []*TrieNode
}

func NewArrayNodeContainer() *ArrayNodeContainer {
	return &ArrayNodeContainer{
		nodes: make([]*TrieNode, 0),
	}
}

func (a *ArrayNodeContainer) Get(c rune) *TrieNode {
	for _, node := range a.nodes {
		if node.val == c {
			return node
		}
	}
	return nil
}

func (a *ArrayNodeContainer) Add(node *TrieNode) {
	// 先检查是否已存在，避免重复添加
	for _, n := range a.nodes {
		if n.val == node.val {
			return
		}
	}
	a.nodes = append(a.nodes, node)
}

func (a *ArrayNodeContainer) Remove(c rune) *TrieNode {
	for i, node := range a.nodes {
		if node.val == c {
			// 从数组中删除
			a.nodes = append(a.nodes[:i], a.nodes[i+1:]...)
			return node
		}
	}
	return nil
}

func (a *ArrayNodeContainer) Size() int {
	return len(a.nodes)
}

func (a *ArrayNodeContainer) GetAll() []*TrieNode {
	// 返回副本，避免外部修改内部数组
	copyNodes := make([]*TrieNode, len(a.nodes))
	copy(copyNodes, a.nodes)
	return copyNodes
}

func (a *ArrayNodeContainer) Transform() NodeContainer {
	// 数组容器不需要转换，直接返回自身
	return a
}
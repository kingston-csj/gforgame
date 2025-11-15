package trie

// TrieNode Trie树节点
// 每一个节点代表一个字符，节点下面包含多个子节点
// @since 2.4.0
type TrieNode struct {
	val      rune           // 当前节点的字符值
	children NodeContainer  // 所有孩子子节点
	isLeaf   bool           // 是否是叶子节点（即是否是敏感词的最后一个字符）
}

// NewTrieNode 创建新的Trie节点
func NewTrieNode(val rune) *TrieNode {
	return &TrieNode{
		val:      val,
		children: NewMapNodeContainer(), // 默认使用map容器
		isLeaf:   false,
	}
}

// IsLeaf 判断是否为叶子节点
func (n *TrieNode) IsLeaf() bool {
	return n.isLeaf
}

// SetLeaf 设置是否为叶子节点
func (n *TrieNode) SetLeaf(leaf bool) {
	n.isLeaf = leaf
}

// AddChild 递归添加子节点
func (n *TrieNode) AddChild(cs string, index int) {
	if index >= len(cs) {
		n.SetLeaf(true)
		return
	}

	val := rune(cs[index])
	child := n.children.Get(val)
	if child == nil {
		child = NewTrieNode(val)
		n.children.Add(child)
	}
	child.AddChild(cs, index+1)
}

// RemoveChild 删除子节点
// @param character 要删除的字符
// @return 被删除的节点，如果不存在则返回nil
// @since 2.5.0
func (n *TrieNode) RemoveChild(character rune) *TrieNode {
	return n.children.Remove(character)
}

// HasPrefix 检查是否包含指定前缀
// @param cs 要检查的字符串
// @param idx 当前处理的字符索引
// @return 找到的前缀结束索引，未找到返回-1
func (n *TrieNode) HasPrefix(cs string, idx int) int {
	if idx >= len(cs) {
		if n.IsLeaf() {
			return idx
		}
		return -1
	}

	val := rune(cs[idx])
	child := n.children.Get(val)
	findIndex := -1

	if child != nil {
		findIndex = child.HasPrefix(cs, idx+1)
	}

	if findIndex != -1 {
		return findIndex
	} else if n.IsLeaf() {
		return idx
	} else {
		return -1
	}
}

// HasExactWord 检查是否精确匹配单词
// @param cs 要检查的字符串
// @param idx 当前处理的字符索引
// @return 是否精确匹配
// @since 2.5.0
func (n *TrieNode) HasExactWord(cs string, idx int) bool {
	if idx >= len(cs) {
		return n.IsLeaf()
	}

	val := rune(cs[idx])
	child := n.children.Get(val)
	if child == nil {
		return false
	}

	return child.HasExactWord(cs, idx+1)
}

// GetChild 获取指定字符的子节点
func (n *TrieNode) GetChild(c rune) *TrieNode {
	return n.children.Get(c)
}

// GetChildren 获取所有子节点
func (n *TrieNode) GetChildren() []*TrieNode {
	return n.children.GetAll()
}

// GetValue 获取当前节点的字符值
func (n *TrieNode) GetValue() rune {
	return n.val
}
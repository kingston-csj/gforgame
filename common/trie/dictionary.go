package trie

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// TrieDictionary Trie树（字典树/前缀树）
// 主要用于解决通过前缀来联想完整单词的问题
// 可用于脏词检测，好友模糊查询等场景
// @since 2.4.0
type TrieDictionary struct {
	threshold int       // 阈值，当节点的孩子节点数量小于等于阈值时，将map容器转化为数组
	root      *TrieNode // 前缀根节点
}

// NewTrieDictionary 创建Trie字典
func NewTrieDictionary() *TrieDictionary {
	return &TrieDictionary{
		threshold: 3, // 默认阈值为3
		root:      NewTrieNode(0), // 根节点字符值为0（空字符）
	}
}

// AddNode 添加单词节点
func (t *TrieDictionary) AddNode(word string) {
	word = t.normalize(word)
	if word == "" {
		return
	}
	t.root.AddChild(word, 0)
}

// DeleteNode 删除单词节点
// @param word 要删除的单词
// @return 是否成功删除
// @since 2.5.0
func (t *TrieDictionary) DeleteNode(word string) bool {
	word = t.normalize(word)
	if word == "" {
		return false
	}
	return t.deleteNodeRecursive(t.root, word, 0)
}

// deleteNodeRecursive 递归删除单词节点
func (t *TrieDictionary) deleteNodeRecursive(node *TrieNode, word string, index int) bool {
	// 如果已经处理完所有字符
	if index >= len(word) {
		// 如果当前节点是叶子节点，则删除叶子标记
		if node.IsLeaf() {
			node.SetLeaf(false)
			return true
		}
		return false
	}

	currentChar := rune(word[index])
	childNode := node.GetChild(currentChar)

	if childNode == nil {
		// 单词不存在
		return false
	}

	// 递归删除下一个字符
	deleted := t.deleteNodeRecursive(childNode, word, index+1)

	if deleted {
		// 如果子节点被删除且当前子节点没有其他子节点且不是叶子节点，则删除当前子节点
		if !childNode.IsLeaf() && childNode.children.Size() == 0 {
			node.RemoveChild(currentChar)
		}
	}

	return deleted
}

// ContainsWords 检查指定字符串是否包含敏感字
// @param word 要检查的字符串
// @return 是否包含敏感字
func (t *TrieDictionary) ContainsWords(word string) bool {
	word = t.normalize(word)
	for i := 0; i < len(word); {
		if end := t.root.HasPrefix(word, i); end != -1 {
			return true
		}
		// 移动到下一个rune（处理中文等多字节字符）
		_, step := utf8.DecodeRuneInString(word[i:])
		i += step
	}
	return false
}

// ContainsExactWord 检查字典是否精确匹配某个单词
// 例如"张无"是敏感字，但是"张无忌"不是
// @param word 要检查的单词
// @return 是否精确匹配
// @since 2.5.0
func (t *TrieDictionary) ContainsExactWord(word string) bool {
	word = t.normalize(word)
	if word == "" {
		return false
	}
	return t.root.HasExactWord(word, 0)
}

// ReplaceWords 将敏感字替换成字符'*'
// @param content 要处理的字符串
// @return 转换后的字符串
func (t *TrieDictionary) ReplaceWords(content string) string {
	normalizedString := t.normalize(content)
	if normalizedString == "" {
		return content
	}

	// 关键修复：先将字符串转为 rune 切片（字符序列），所有操作基于字符索引
	normalizedRunes := []rune(normalizedString)
	runeLen := len(normalizedRunes) // 字符长度（而非字节长度）
	if runeLen == 0 {
		return content
	}

	// 记录需要替换的字符索引区间 [start, end)（闭开区间）
	type replaceRange struct {
		start, end int
	}
	var ranges []replaceRange

	// 基于字符索引遍历，避免字节/字符索引混淆
	for i := 0; i < runeLen; {
		// 将当前字符及后续字符转回字符串，用于 HasPrefix 匹配
		// 注意：HasPrefix 接收的是字符串，这里需要正确截取字符片段
		subStr := string(normalizedRunes[i:])
		// HasPrefix 返回的是字节索引的结束位置，需要转换为字符索引
		byteEnd := t.root.HasPrefix(subStr, 0)
		if byteEnd != -1 {
			// 计算字节索引对应的字符索引：subStr[0:byteEnd] 包含的字符数
			charEnd := utf8.RuneCountInString(subStr[:byteEnd])
			// 记录字符索引区间
			ranges = append(ranges, replaceRange{
				start: i,
				end:   i + charEnd,
			})
			// 跳过已匹配的字符
			i += charEnd
		} else {
			// 移动到下一个字符（字符索引 +1）
			i++
		}
	}

	if len(ranges) == 0 {
		return content
	}

	// 直接操作 rune 切片替换（字符索引安全）
	for _, r := range ranges {
		// 确保区间不越界（防御性编程）
		if r.start >= runeLen {
			continue
		}
		// 修正 end 不超过总长度
		end := r.end
		if end > runeLen {
			end = runeLen
		}
		// 替换为 '*'
		for i := r.start; i < end; i++ {
			normalizedRunes[i] = '*'
		}
	}

	// 转回字符串返回
	return string(normalizedRunes)
}

// normalize 字符串预处理：英文统一转小写，只保留字母、数字、中文
func (t *TrieDictionary) normalize(dirtyWord string) string {
	if dirtyWord == "" {
		return ""
	}

	var sb strings.Builder
	for _, c := range dirtyWord {
		if unicode.IsLetter(c) || unicode.IsDigit(c) || t.isChineseCharacter(c) {
			// 转小写
			sb.WriteRune(unicode.ToLower(c))
		}
	}
	return sb.String()
}

// isChineseCharacter 判断是否为中文字符
func (t *TrieDictionary) isChineseCharacter(c rune) bool {
	// CJK统一汉字的Unicode范围
	return c >= 0x4E00 && c <= 0x9FFF ||
		c >= 0x3400 && c <= 0x4DBF || // CJK扩展A
		c >= 0xF900 && c <= 0xFAFF    // CJK兼容汉字
}

// GetRoot 获取根节点
func (t *TrieDictionary) GetRoot() *TrieNode {
	return t.root
}

// Rebuild 整颗树构建成功后，对孩子节点重新构造
// 如果某节点的孩子节点数量少于阈值，则将map容器转化为数组
func (t *TrieDictionary) Rebuild() {
	t.rebuildChildren(t.root)
}

func (t *TrieDictionary) rebuildChildren(node *TrieNode) {
	if mapContainer, ok := node.children.(*MapNodeContainer); ok {
		if mapContainer.Size() <= t.threshold {
			node.children = mapContainer.Transform()
		}
	}

	// 递归处理子节点
	for _, child := range node.GetChildren() {
		t.rebuildChildren(child)
	}
}
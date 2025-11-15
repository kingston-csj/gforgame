package trie

import (
	"fmt"
	"testing"
)

func TestLog(t *testing.T) {
	// 创建Trie字典
	dict :=  NewTrieDictionary()

	// 添加敏感词
	dict.AddNode("敏感词1")
	dict.AddNode("敏感词2")
	dict.AddNode("测试")

	// 构建优化（可选）
	dict.Rebuild()

	// 检查是否包含敏感词
	content := "这是一个敏感词1的测试内容"
	fmt.Println(dict.ContainsWords(content)) // 输出: true

	// 精确匹配检查
	fmt.Println(dict.ContainsExactWord("测试"))   // 输出: true
	fmt.Println(dict.ContainsExactWord("测试1")) // 输出: false

	// 替换敏感词
	result := dict.ReplaceWords(content)
	fmt.Println(result) // 输出: 这是一个*****的**内容

	// 删除敏感词
	dict.DeleteNode("敏感词1")
	fmt.Println(dict.ContainsWords("敏感词1")) // 输出: false
}
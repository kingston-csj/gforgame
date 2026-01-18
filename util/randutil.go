package util

import (
	"errors"
	"math/rand"
	"time"
)

// 初始化随机数种子（保证每次运行结果不同）
func init() {
	// Go 1.20+ math/rand 全局实例已并发安全，无需每个goroutine创建
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// NextInt 等价于 Java ThreadLocalRandom.current().nextInt()
// 返回一个随机的 int 类型值
func NextInt() int {
	return rand.Int()
}

// NextIntN 等价于 Java ThreadLocalRandom.current().nextInt(n)
// 返回 [0, n) 范围内的随机 int
// 注意：n <= 0 时会触发 panic（和 Java 行为一致）
func NextIntN(n int) int {
	return rand.Intn(n)
}

// RandomValue 等价于 Java randomValue(min, max)
// 返回 [min, max] 范围内的随机 int
// 参数校验：min > max 时返回错误
func RandomValue(min, max int) (int, error) {
	if min > max {
		return 0, errors.New("min > max")
	}
	if min == max {
		return min, nil
	}
	// 计算 [min, max] 区间长度，rand.Intn 返回 [0, len)，加 min 得到目标范围
	return min + rand.Intn(max-min+1), nil
}

// RandomIndex 等价于 Java randomIndex(probs)
// 根据概率切片计算随机索引：
// 1. 空切片返回 -1
// 2. 非空时计算概率总和，随机一个数，累加概率找到第一个命中的索引
// 3. 遍历完未命中返回错误（等价于 Java 的 IllegalArgumentException）
func RandomIndex(probs []int) (int, error) {
	if len(probs) == 0 {
		return -1, nil
	}

	// 计算概率总和（等价于 Java stream.reduce(0, Integer::sum)）
	sum := 0
	for _, p := range probs {
		sum += p
	}
	if sum <= 0 {
		return -1, errors.New("probability sum is zero")
	}

	// 生成 [0, sum) 范围内的随机数
	randomValue := rand.Intn(sum)
	accumulated := 0

	// 遍历累加概率，找到第一个命中的索引
	for i, p := range probs {
		accumulated += p
		if randomValue < accumulated {
			return i, nil
		}
	}

	// 遍历完未命中（理论上不会走到这里，除非概率总和计算错误）
	return -1, errors.New("randomIndex out of range")
}

// RandomIndexList 等价于 Java randomIndexList
// 返回 count 个随机索引，支持是否「去重选中」（remove=true 时选中后概率置0）
// 参数校验：
// 1. 概率切片为空 → 错误
// 2. 概率包含负数 → 错误
// 3. count <= 0 → 错误
// 4. remove=true 且 count > 切片长度 → 错误
func RandomIndexList(probabilityList []int, count int, remove bool) ([]int, error) {
	// 基础参数校验
	if len(probabilityList) == 0 {
		return nil, errors.New("probabilityList is empty")
	}
	// 检查概率是否包含负数
	for _, p := range probabilityList {
		if p < 0 {
			return nil, errors.New("probabilityList contains negative number")
		}
	}
	if count <= 0 {
		return nil, errors.New("count <= 0")
	}
	if remove && count > len(probabilityList) {
		return nil, errors.New("count > probabilityList size")
	}

	// 存储选中的索引
	hits := make([]int, 0, count)

	// 循环生成 count 个索引
	for i := 0; i < count; i++ {
		index, err := RandomIndex(probabilityList)
		if err != nil {
			return nil, err
		}
		hits = append(hits, index)

		// remove=true 时，将选中索引的概率置0（下次不会被选中）
		if remove {
			probabilityList[index] = 0
		}
	}

	return hits, nil
}
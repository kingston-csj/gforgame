package data

import (
	"fmt"
	"io/github/gforgame/logger"
	"testing"
)

func TestExcelReader(t *testing.T) {
	// 创建 ExcelDataReader 实例
	reader := NewExcelDataReader(true)

	type RewardDef struct {
		Type  string `json:"type" excel:"type"`
		Value string `json:"value" excel:"value"`
	}

	type ConsumeDef struct {
		Type  string `json:"type" excel:"type"`
		Value string `json:"value" excel:"value"`
	}

	type Name struct {
		Id       int64        `json:"id" excel:"id"`
		Name     string       `json:"type" excel:"name"`
		Rewards  []RewardDef  `json:"rewards" excel:"rewards"`
		Consumes []ConsumeDef `json:"consumes" excel:"consumes"`
	}

	// 读取 Excel 文件
	result, err := reader.Read("mall.xlsx", &Name{})
	if err != nil {
		logger.Error(fmt.Errorf("session.Send: %v", err))
	}

	// 打印结果
	for _, item := range result {
		fmt.Printf("%+v\n", item)
	}
}

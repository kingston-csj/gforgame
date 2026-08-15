package data_test

import (
	"fmt"
	"testing"

	"github.com/forfun/gforgame/data"
)

type RewardDef struct {
	Type  string `json:"type" excel:"type"`
	Value string `json:"value" excel:"value"`
}

type ConsumeDef struct {
	Type  string `json:"type" excel:"type"`
	Value string `json:"value" excel:"value"`
}

type Mall struct {
	Id       int32        `json:"id" excel:"id"`
	Type     int32        `json:"type" excel:"type"`
	Name     string       `json:"name" excel:"name"`
	Rewards  []RewardDef  `json:"rewards" excel:"rewards"`
	Consumes []ConsumeDef `json:"consumes" excel:"consumes"`
}

type Item struct {
	Id      int32  `json:"id" excel:"id"`
	Name    string `json:"name" excel:"name"`
	Quality int32  `json:"quality" excel:"quality"`
	Tips    string `json:"tips" excel:"tips"`
	Icon    string `json:"icon" excel:"icon"`
}

func TestDataContainer(t *testing.T) {
	// 创建 ExcelDataReader
	reader := data.NewExcelDataReader(true)

	// 读取 Excel 文件（强类型）
	records, err := data.ReadTyped[Mall](reader, "mall.xlsx")
	if err != nil {
		fmt.Println("Failed to read Excel file:", err)
		return
	}

	container, err := data.BuildContainer("mall", records,
		func(record *Mall) int32 { return record.Id },
		map[string]func(*Mall) any{
			"type": func(record *Mall) any { return record.Type },
		},
	)
	if err != nil {
		fmt.Printf("Failed to build container: %v\n", err)
		return
	}

	// 查询记录
	fmt.Println("All records:", container.GetAllRecords())
	target := container.GetRecord(1)
	fmt.Println("Record with ID 1:", target)
	fmt.Println("Records with type 2:", container.GetRecordsByIndex("type", 2))
}

func TestMultiDataContainer(t *testing.T) {
	// 创建 ExcelDataReader
	reader := data.NewExcelDataReader(true)
	// 处理商城表（强类型）
	mallContainer, err := data.ProcessTableTyped[Mall](
		reader,
		"mall",
		"mall.xlsx",
		func(record *Mall) int32 { return record.Id },
		map[string]func(*Mall) any{
			"type": func(record *Mall) any { return record.Type },
		},
	)
	if err != nil {
		fmt.Printf("Failed to process table mall: %v\n", err)
		return
	}

	// 处理道具表（强类型）
	itemContainer, err := data.ProcessTableTyped[Item](
		reader,
		"item",
		"item.xlsx",
		func(record *Item) int32 { return record.Id },
		nil,
	)
	if err != nil {
		fmt.Printf("Failed to process table item: %v\n", err)
		return
	}

	fmt.Println("All records in Mall table:", mallContainer.GetAllRecords())
	target := mallContainer.GetRecord(1)
	fmt.Println("Record with ID 1:", target)
	fmt.Println("Records with type 2 in Mall table:", mallContainer.GetRecordsByIndex("type", 2))

	fmt.Println("All records in Item table:", itemContainer.GetAllRecords())
	target2 := itemContainer.GetRecord(1)
	target3 := itemContainer.GetRecord(1)
	fmt.Println(target2 == target3)
	fmt.Println("Record with ID 1:", target2)
}

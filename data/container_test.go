package data

import (
	"fmt"
	"reflect"
	"testing"
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
	Id       int64        `json:"id" excel:"id"`
	Type     int64        `json:"type" excel:"type"`
	Name     string       `json:"name" excel:"name"`
	Rewards  []RewardDef  `json:"rewards" excel:"rewards"`
	Consumes []ConsumeDef `json:"consumes" excel:"consumes"`
}

type Item struct {
	Id      int64  `json:"id" excel:"id"`
	Name    string `json:"name" excel:"name"`
	Quality int64  `json:"quality" excel:"quality"`
	Tips    string `json:"tips" excel:"tips"`
	Icon    string `json:"icon" excel:"icon"`
}

func TestDataContainer(t *testing.T) {

	// 创建 ExcelDataReader
	reader := NewExcelDataReader(true)

	// 读取 Excel 文件
	records, err := reader.Read("mall.xlsx", &Mall{})
	if err != nil {
		fmt.Println("Failed to read Excel file:", err)
		return
	}

	// 创建 Container
	container := NewContainer[int64, Mall]()

	// 定义 ID 获取函数和索引函数
	getIdFunc := func(record Mall) int64 {
		return record.Id
	}
	indexFuncs := map[string]func(Mall) interface{}{
		"type": func(record Mall) interface{} {
			return record.Type
		},
	}

	// 将记录注入容器
	nameRecords := make([]Mall, len(records))
	for i, record := range records {
		nameRecords[i] = record.(Mall)
	}
	container.Inject(nameRecords, getIdFunc, indexFuncs)

	// 查询记录
	fmt.Println("All records:", container.GetAllRecords())
	target, _ := container.GetRecord(1)
	fmt.Println("Record with ID 1:", target)
	fmt.Println("Records with type 2:", container.GetRecordsBy("type", 2))
}

func TestMultiDataContainer(t *testing.T) {
	// 创建 ExcelDataReader
	reader := NewExcelDataReader(true)

	// 定义表配置
	tableConfigs := []TableMeta{
		// 商城表
		{
			TableName:  "mall",
			IDField:    "Id",
			IndexFuncs: map[string]string{"type": "Type"},
			RecordType: reflect.TypeOf(Mall{}),
		},
		// 道具表
		{
			TableName:  "item",
			IDField:    "Id",
			RecordType: reflect.TypeOf(Item{}),
		},
	}

	// 处理每张表
	containers := make(map[string]*Container[int64, interface{}])
	for _, config := range tableConfigs {
		container, err := ProcessTable(reader, config.TableName+".xlsx", config)
		if err != nil {
			fmt.Printf("Failed to process table %s: %v\n", config.TableName, err)
			continue
		}
		containers[config.TableName] = container
	}

	// 查询商城记录
	mallContainer := containers["mall"]
	fmt.Println("All records in Mall table:", mallContainer.GetAllRecords())
	target, _ := mallContainer.GetRecord(1)
	fmt.Println("Record with ID 1:", target)
	fmt.Println("Records with type 2 in Mall table:", mallContainer.GetRecordsBy("type", 2))

	// 查询商城记录
	itemContainer := containers["item"]
	fmt.Println("All records in Mall table:", itemContainer.GetAllRecords())
	target2, _ := itemContainer.GetRecord(1)
	fmt.Println("Record with ID 1:", target2)
}

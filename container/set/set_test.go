package set_test

import (
	"encoding/json"
	"io/github/gforgame/container/set"
	"testing"
)

func TestSet_Json(t *testing.T) {
	// 1. 测试序列化 (Set -> JSON Array)
	s := set.NewSet[int32]()
	s.Add(1)
	s.Add(2)
	s.Add(3)

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	jsonStr := string(data)
	// 顺序不固定，检查包含即可
	// Expected roughly: [1,2,3] (order may vary)
	t.Logf("Serialized JSON: %s", jsonStr)

	// 2. 测试反序列化 (JSON Array -> Set) 并去重
	// 输入包含重复元素 [1, 2, 2, 3]
	inputJson := `[1, 2, 2, 3]`
	var s2 set.Set[int32]
	err = json.Unmarshal([]byte(inputJson), &s2)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if s2.Len() != 3 {
		t.Errorf("Expected length 3, got %d", s2.Len())
	}
	if !s2.Contains(1) || !s2.Contains(2) || !s2.Contains(3) {
		t.Errorf("Set missing elements")
	}
	t.Logf("Deserialized Set size: %d", s2.Len())
}

type TestStruct struct {
	Ids *set.Set[int32] `json:"ids"`
}

func TestSet_StructField(t *testing.T) {
	// 测试作为结构体字段
	jsonStr := `{"ids": [10, 20, 20, 30]}`
	var ts TestStruct
	err := json.Unmarshal([]byte(jsonStr), &ts)
	if err != nil {
		t.Fatalf("Unmarshal struct failed: %v", err)
	}

	if ts.Ids == nil {
		t.Fatal("Ids field is nil")
	}
	if ts.Ids.Len() != 3 {
		t.Errorf("Expected length 3, got %d", ts.Ids.Len())
	}
	if !ts.Ids.Contains(10) || !ts.Ids.Contains(20) || !ts.Ids.Contains(30) {
		t.Errorf("Set missing elements")
	}
}

package jsonutil_test

import (
	"io/github/gforgame/util/jsonutil"
	"testing"
)

type TestUser struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Tags []string `json:"tags,omitempty"`
}

func TestJsonUtil(t *testing.T) {
	// 1. 测试 StructToJSON
	user := TestUser{Name: "Alice", Age: 25, Tags: []string{"golang", "game"}}
	jsonStr, err := jsonutil.StructToJSON(user)
	if err != nil {
		t.Fatalf("StructToJSON failed: %v", err)
	}
	t.Logf("StructToJSON: %s", jsonStr)

	// 2. 测试 JsonToStruct
	var user2 TestUser
	err = jsonutil.JsonToStruct(jsonStr, &user2)
	if err != nil {
		t.Fatalf("JsonToStruct failed: %v", err)
	}
	if user2.Name != user.Name || user2.Age != user.Age || len(user2.Tags) != len(user.Tags) {
		t.Errorf("JsonToStruct result mismatch: %+v", user2)
	}

	// 3. 测试 HTML Escape 禁用
	htmlStr := "<div>Hello & Welcome</div>"
	htmlJson, _ := jsonutil.StructToJSON(htmlStr)
	expectedHtmlJson := "\"" + htmlStr + "\""
	if htmlJson != expectedHtmlJson {
		t.Errorf("HTML escaping should be disabled. Expected %s, got %s", expectedHtmlJson, htmlJson)
	}
	t.Logf("HTML JSON (no escape): %s", htmlJson)

	// 4. 测试 Pretty JSON
	prettyJson, _ := jsonutil.StructToPrettyJSON(user)
	t.Logf("Pretty JSON:\n%s", prettyJson)

	// 5. 测试 omitempty
	userEmpty := TestUser{Name: "Bob", Age: 30}
	emptyJson, _ := jsonutil.StructToJSON(userEmpty)
	t.Logf("Empty Tags JSON: %s", emptyJson)
}

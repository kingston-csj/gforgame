package cache

import (
	"fmt"
	mysqldb "io/github/gforgame/db"
	"testing"
)

type Player struct {
	Id    int64
	Name  string
	Level uint
}

func TestCache(t *testing.T) {
	cm := NewCacheManager()
	dbLoader := func(key string) (interface{}, error) {
		var p Player
		mysqldb.Db.First(&p, "id=?", key)
		return &p, nil
	}
	cm.Register("player", dbLoader)

	cache, err := cm.GetCache("player")
	if err != nil {
		t.Error(err)
	}

	// 测试缓存
	key := "1001"
	p, err := cache.Get(key)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Printf("first query %s: %v\n", key, p)
	p2, ok := p.(*Player)
	if ok {
		p2.Name = "gforgam2"
		//// 使用 Set 方法更新缓存
		cache.Set(key, p2)
		p, err = cache.Get(key)
		fmt.Printf("second query %s: %v\n", key, p)
	}
}

package container

import (
	"fmt"
	"io/github/gforgame/examples/rank/model"

	"github.com/emirpasic/gods/maps/treemap"
)

// ConcurrentRankContainer 并发排行榜容器
// 只通过channel和内部goroutine并发安全
type ConcurrentRankContainer struct {
	ranks    *treemap.Map  // 红黑树数据结构
	capacity int          // 容量
	cmdChan  chan any     // 命令通道
}

type addCmd struct {
	key   any
	value any
	done  chan struct{}
}

type removeCmd struct {
	key  any
	done chan struct{}
}

type updateCmd struct {
	key   any
	value any
	done  chan struct{}
}

type getCmd struct {
	key  any
	resp chan any
}

type getItemsCmd struct {
	resp chan []RankEntry
}

type containsCmd struct {
	key  any
	resp chan bool
}

type rankSizeCmd struct {
	resp chan int
}

type closeCmd struct {
	done chan struct{}
}

// NewConcurrentRankContainer 创建一个新的并发排行榜容器
func NewConcurrentRankContainer(capacity int) *ConcurrentRankContainer {
	c := &ConcurrentRankContainer{
		ranks:    treemap.NewWith(model.CompareRank),
		capacity: capacity, 
		cmdChan:  make(chan any, 1000),
	}
	go c.run()
	return c
}

// run goroutine，串行处理所有命令
func (c *ConcurrentRankContainer) run() {
	for cmd := range c.cmdChan {
		switch v := cmd.(type) {
		case addCmd:
			c.ranks.Put(v.value, v.key)
			if c.ranks.Size() > c.capacity {
				// 移除最小的元素
				it := c.ranks.Iterator()
				it.Last()
				c.ranks.Remove(it.Key())
			}
			close(v.done)
		case removeCmd:
			// 需要遍历找到对应的key
			it := c.ranks.Iterator()
			for it.Next() {
				if it.Value() == v.key {
					c.ranks.Remove(it.Key())
					break
				}
			}
			close(v.done)
		case updateCmd:
			// 先删除旧值
			it := c.ranks.Iterator()
			for it.Next() {
				if it.Value() == v.key {
					c.ranks.Remove(it.Key())
					break
				}
			}
			// 添加新值
			c.ranks.Put(v.value, v.key)
			if c.ranks.Size() > c.capacity {
				it := c.ranks.Iterator()
				it.Last()
				c.ranks.Remove(it.Key())
			}
			close(v.done)
		case getCmd:
			it := c.ranks.Iterator()
			for it.Next() {
				if it.Value() == v.key {
					v.resp <- it.Key()
					break
				}
			}
			close(v.resp)
		case getItemsCmd:
			items := make([]RankEntry, 0, c.ranks.Size())
			it := c.ranks.Iterator()
			for it.Next() {
				items = append(items, RankEntry{
					Key:   it.Value(),
					Value: it.Key().(model.BaseRank),
				})
			}
			v.resp <- items
			close(v.resp)
		case containsCmd:
			exists := false
			it := c.ranks.Iterator()
			for it.Next() {
				if it.Value() == v.key {
					exists = true
					break
				}
			}
			v.resp <- exists
			close(v.resp)
		case rankSizeCmd:
			fmt.Println("rankSizeCmd", c.ranks.Size())
			v.resp <- c.ranks.Size()
			close(v.resp)
		case closeCmd:
			close(v.done)
			return
		}
	}
}

// 对外方法全部通过channel
func (c *ConcurrentRankContainer) Add(key, value any) {
	done := make(chan struct{})
	c.cmdChan <- addCmd{key, value, done}
	<-done
}

func (c *ConcurrentRankContainer) Remove(key any) {
	done := make(chan struct{})
	c.cmdChan <- removeCmd{key, done}
	<-done
}

func (c *ConcurrentRankContainer) Update(key, value any) {
	done := make(chan struct{})
	c.cmdChan <- updateCmd{key, value, done}
	<-done
}

func (c *ConcurrentRankContainer) Get(key any) any {
	resp := make(chan any, 1)
	c.cmdChan <- getCmd{key, resp}
	val, ok := <-resp
	if !ok {
		return nil
	}
	return val
}

func (c *ConcurrentRankContainer) GetItems() []RankEntry {
	resp := make(chan []RankEntry, 1)
	c.cmdChan <- getItemsCmd{resp}
	return <-resp
}

func (c *ConcurrentRankContainer) Contains(key any) bool {
	resp := make(chan bool, 1)
	c.cmdChan <- containsCmd{key, resp}
	return <-resp
}

func (c *ConcurrentRankContainer) RankSize() int {
	resp := make(chan int, 1)
	c.cmdChan <- rankSizeCmd{resp}
	return <-resp
}

// 停止容器, 停服时调用
func (c *ConcurrentRankContainer) Stop() {
	done := make(chan struct{})
	c.cmdChan <- closeCmd{done}
	<-done
	close(c.cmdChan) // 关闭通道，防止泄漏
}

package actor_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/forfun/gforgame/actor"
)

// ---------------- 业务自定义消息 ----------------
// MsgPlayerDonate 玩家发起捐赠，发给公会Actor
type MsgPlayerDonate struct {
	PlayerID int64
	AddVal   int64
}

// MsgGetContribution 查询玩家贡献
type MsgGetContribution struct {
}

// MsgDoDonate 触发玩家执行捐赠
type MsgDoDonate struct {
	Val int64
}

// ---------------- guildActor 业务实现 ----------------
type guildActor struct {
	guildID             int64
	contribution  int64
}

func newGuildActor(guildID int64) actor.Actor {
	return &guildActor{
		guildID:             guildID,
	}
}

func (g *guildActor) OnStart() {}

func (g *guildActor) OnMessage(msg actor.Message) actor.Message {
	switch m := msg.(type) {
	case MsgPlayerDonate:
		g.contribution += m.AddVal
		return nil
	case MsgGetContribution:
		return g.contribution
	case struct{}: // syncBarrier栅栏消息
		return nil
	default:
		return nil
	}
}

func (g *guildActor) OnStop() {}

// ---------------- playerActor 业务实现 ----------------
type playerActor struct {
	playerID int64
	guildRef *actor.ActorRef
}

func newPlayerActor(playerID int64, guildRef *actor.ActorRef) actor.Actor {
	return &playerActor{
		playerID: playerID,
		guildRef: guildRef,
	}
}

func (p *playerActor) OnStart() {}

func (p *playerActor) OnMessage(msg actor.Message) actor.Message {
	switch m := msg.(type) {
	case MsgDoDonate:
		_ = p.guildRef.Tell(MsgPlayerDonate{
			PlayerID: p.playerID,
			AddVal:   m.Val,
		})
		return nil
	default:
		return nil
	}
}

func (p *playerActor) OnStop() {}

// ---------------- 普通单元测试 ----------------
func TestActorSystem_TwoPlayerDonate(t *testing.T) {
	sys := actor.NewActorSystem()
	guildPath := actor.NewActorPath("game", "guild", "1")
	guildRef := sys.Spawn(guildPath, newGuildActor(1))

	p1Path := actor.NewActorPath("game", "player", "1001")
	player1Ref := sys.Spawn(p1Path, newPlayerActor(1001, guildRef))
	p2Path := actor.NewActorPath("game", "player", "1002")
	player2Ref := sys.Spawn(p2Path, newPlayerActor(1002, guildRef))

	for i := 0; i < 3; i++ {
		err := player1Ref.Tell(MsgDoDonate{Val: 100})
		if err != nil {
			t.Fatalf("player1 tell fail: %v", err)
		}
		err = player2Ref.Tell(MsgDoDonate{Val: 100})
		if err != nil {
			t.Fatalf("player2 tell fail: %v", err)
		}
	}

	// FIFO栅栏，等待前面全部消息处理完成
	type syncBarrier struct{}
	_, err := guildRef.Ask(syncBarrier{} )
	if err != nil {
		t.Fatalf("barrier wait fail %v", err)
	}

	resp1, err := guildRef.Ask(MsgGetContribution{})
	if err != nil {
		t.Fatalf("ask p1 fail: %v", err)
	}
	if v, ok := resp1.(int64); !ok || v != 600 {
		t.Errorf("player1 want 600, got %v", resp1)
	}


	sys.Stop(p1Path.String())
	sys.Stop(p2Path.String())
	sys.Stop(guildPath.String())
	time.Sleep(100 * time.Millisecond)
}

// ---------------- Example 文档测试，godoc会提取该示例 ----------------
func ExampleActorSystem_TwoPlayerDonate() {
	sys := actor.NewActorSystem()
	guildPath := actor.NewActorPath("game", "guild", "g1")
	guildRef := sys.Spawn(guildPath, newGuildActor(1))

	p1Path := actor.NewActorPath("game", "player", "p1001")
	player1Ref := sys.Spawn(p1Path, newPlayerActor(1001, guildRef))


	p2Path := actor.NewActorPath("game", "player", "p1002")
	player2Ref := sys.Spawn(p2Path, newPlayerActor(1002, guildRef))

	for i := 0; i < 3; i++ {
		_ = player1Ref.Tell(MsgDoDonate{Val: 100})
		_ = player2Ref.Tell(MsgDoDonate{Val: 100})
	}

	type syncBarrier struct{}
	_, _ = guildRef.Ask(syncBarrier{})

	resp1, _ := guildRef.Ask(MsgGetContribution{})

	fmt.Printf("guild contribution=%d\n", resp1.(int64))

	sys.Stop(p1Path.String())
	sys.Stop(p2Path.String())
	sys.Stop(guildPath.String())
	time.Sleep(100 * time.Millisecond)
	// Output:
	// guild contribution=600
}

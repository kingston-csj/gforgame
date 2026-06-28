package player

// 竞技场数据
type ArenaBox struct {

	DefenseTeam string 
	// 最近的战斗记录
	MatchRecords []*MatchRecord 
	// 最近匹配的对手id
	matchMemberIds []string
	// 竞技场门票
	Ticket int32
	// 竞技场战斗次数
	ChallengeTimes int32
}

func (b *ArenaBox) AfterLoad() {
	if b.MatchRecords == nil {
		b.MatchRecords = make([]*MatchRecord, 0)
	}
}

type MatchRecord struct {
	Id string
	OpponentId string
	Time int64
	Winner string
	Score int32
	// 是否主动进攻
	IsAttack bool
}
package buff

type Buff struct {
	Id      string
	ModelId int32
	// 回合数， 0表示永久
	Duration int32
}

func (b *Buff) GetId() string {
	return b.Id
}

func (b *Buff) TimeToDead() bool {
	if b.Duration != 0 {
		b.Duration--
		if b.Duration == 0 {
			return true
		}
	}
	return false
}

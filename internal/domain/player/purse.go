package player

type Purse struct {
	Diamond int32
	Gold    int32
}

func (p *Purse) AddDiamond(amount int32) {
	p.Diamond += amount
}

func (p *Purse) AddGold(amount int32) {
	p.Gold += amount
}

func (p *Purse) SubDiamond(amount int32) bool {
	if p.Diamond < amount {
		return false
	}
	p.Diamond -= amount
	return true
}

func (p *Purse) SubGold(amount int32) bool {
	if p.Gold < amount {
		return false
	}
	p.Gold -= amount
	return true
}

func (p *Purse) GetDiamond() int32 {
	return p.Diamond
}

func (p *Purse) IsEnoughGold(amount int32) bool {
	return p.Gold >= amount
}

func (p *Purse) IsEnoughDiamond(amount int32) bool {
	return p.Diamond >= amount
}

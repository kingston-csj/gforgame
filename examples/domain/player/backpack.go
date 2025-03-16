package player

type Backpack struct {
	Items map[int32]int32
}

func (b *Backpack) AddItem(itemId int32, count int32) {
	prevCount, ok := b.Items[itemId]
	if !ok {
		b.Items[itemId] = count
	} else {
		b.Items[itemId] = prevCount + count
	}
}

func (b *Backpack) RemoveItem(itemId int32, count int32) bool {
	if b.Items[itemId] < count {
		return false
	}
	b.Items[itemId] -= count
	return true
}

func (b *Backpack) GetItemCount(itemId int32) int32 {
	return b.Items[itemId]
}

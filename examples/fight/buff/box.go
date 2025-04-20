package buff

import (
	"io/github/gforgame/examples/context"
	"io/github/gforgame/examples/domain/config"
	"io/github/gforgame/examples/fight/attribute"
	"unsafe"
)

type BuffBox struct {
	buffs map[int32][]*Buff
	Attrs map[attribute.AttrType]int32
}

func NewBuffBox() *BuffBox {
	return &BuffBox{
		buffs: make(map[int32][]*Buff),
		Attrs: make(map[attribute.AttrType]int32),
	}
}

func (b *BuffBox) AddBuff(modelId int32) bool {
	buffData := context.GetConfigRecordAs[config.BuffData]("buff", int64(modelId))

	if buffData.Relation == 1 {
		// 可以叠加
		if len(b.buffs[modelId]) < int(buffData.Layer) {
			b.buffs[modelId] = append(b.buffs[modelId], NewBuff(modelId))
			return true
		}
	} else {
		// 直接替换
		b.buffs[modelId] = []*Buff{NewBuff(modelId)}
		return true
	}
	return false
}

func (b *BuffBox) RefreshAttrs() {
	for at := range b.Attrs {
		b.Attrs[at] = 0
	}

	for _, group := range b.buffs {
		for _, buff := range group {
			buffData := context.GetConfigRecordAs[config.BuffData]("buff", int64(buff.ModelId))
			if buffData.Type == 1 {
				// 将基类指针转换为派生类指针
				attrBuff := (*AttrBuff)(unsafe.Pointer(buff))
				for attrType, value := range attrBuff.Attrs {
					b.Attrs[attrType] += value
				}
			}
		}
	}
}

func NewBuff(modelId int32) *Buff {
	buffData := context.GetConfigRecordAs[config.BuffData]("buff", int64(modelId))
	if buffData.Type == 1 {
		// 属性buff
		attrBuff := NewAttrBuff(modelId)
		return (*Buff)(unsafe.Pointer(attrBuff))
	} else {
		// 状态buff
		stateBuff := NewStateBuff(modelId)
		return (*Buff)(unsafe.Pointer(stateBuff))
	}
}

func (b *BuffBox) GetAttrs() map[attribute.AttrType]int32 {
	return b.Attrs
}

func (b *BuffBox) GetAttrValue(attrType attribute.AttrType) int32 {
	v, ok := b.Attrs[attrType]
	if !ok {
		return 0
	}
	return v
}

func (b *BuffBox) CheckBuffLife() bool {
	changed := false
	// 遍历所有buff组
	for groupId, group := range b.buffs {
		// 创建一个新的切片来存储未过期的buff
		aliveBuffs := make([]*Buff, 0, len(group))
		// 遍历当前组中的所有buff
		for _, buff := range group {
			if buff.TimeToDead() {
				changed = true
				continue
			}
			// buff未过期，保留
			aliveBuffs = append(aliveBuffs, buff)
		}

		// 如果组中还有未过期的buff，更新组
		if len(aliveBuffs) > 0 {
			b.buffs[groupId] = aliveBuffs
		} else {
			// 如果组中所有buff都已过期，删除整个组
			delete(b.buffs, groupId)
		}
	}
	if changed {
		// 检查buff后刷新属性
		b.RefreshAttrs()
	}
	return changed
}

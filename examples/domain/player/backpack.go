package player

import (
	"io/github/gforgame/common"
	"io/github/gforgame/examples/config"
	configcontract "io/github/gforgame/examples/config/contracts"
	"io/github/gforgame/examples/constants"
	configdomain "io/github/gforgame/examples/domain/config"
	protos "io/github/gforgame/protos"
	"io/github/gforgame/util"
)

type Item struct {
	Type   int32
	Uid    string
	ItemId int32
	Count  int32
	Level  int32
	Extra string 
}

func (i *Item) ChangeAmount(change int32) int32 {
	i.Count += change
	return i.Count
}

func (i *Item) ToVo() protos.ItemInfo {
	vo :=  protos.ItemInfo{
		Cf_id: i.ItemId,
		Uid: i.Uid,
		Count: i.Count,
		Level: i.Level,
		Extra:  i.Extra,
	}
	return vo
}

var errorIllegalParams = common.NewBusinessRequestException(constants.I18N_COMMON_ILLEGAL_PARAMS)

/// 背包
type Backpack struct {
	Items map[string]*Item
	configProvider      configcontract.ItemConfigProvider `gorm:"-"`
	Capacity            int32 `gorm:"-"`
}

func (b *Backpack) AfterLoad() {
	if b.Items == nil {
		b.Items = make(map[string]*Item)
	}
}

/// 基础道具配置提供器
type BaseItemConfigProvider struct {
}

func (p *BaseItemConfigProvider) GetConfig(itemId int32) configcontract.ItemConfig {
	return config.QueryById[configdomain.PropData](itemId)
}

/// 符文配置提供器
type RuneConfigProvider struct {
}

func (p *RuneConfigProvider) GetConfig(itemId int32) configcontract.ItemConfig {
	return config.QueryById[configdomain.RuneData](itemId)
}

/// 场景道具配置提供器
type SceneItemConfigProvider struct {
}

func (p *SceneItemConfigProvider) GetConfig(itemId int32) configcontract.ItemConfig {
	return config.QueryById[configdomain.ScenePropData](itemId)
}

var (
	BaseItemConfigProviderInstance = &BaseItemConfigProvider{}
	RuneConfigProviderInstance = &RuneConfigProvider{}
	SceneItemConfigProviderInstance = &SceneItemConfigProvider{}
)

type ChangeItem struct {
	from int32
	to  int32
	change int32
	Item   *Item
}

/// 道具变更结果
type ChangeResult struct {
	Succ bool
	ChangedItems []*ChangeItem
}

func (r *ChangeResult) ToChangeInfos() []protos.ItemInfo {
	itemInfos := make([]protos.ItemInfo, 0, len(r.ChangedItems))
	for _, item := range r.ChangedItems {
		itemInfos = append(itemInfos, item.Item.ToVo())
	}
	return itemInfos
}

func (r *ChangeResult) addChanged(item *Item, from int32, to int32, change int32) {
	changeItem := &ChangeItem{
		from:   from,
		to:     to,
		change: change,	
		Item: item,
	}
	r.Succ = true
	r.ChangedItems = append(r.ChangedItems, changeItem)
}

/// 添加道具
func (b *Backpack) AddByModelId(itemId int32, count int32, initFunc func(*Item)) (*ChangeResult, error) {
	if itemId <= 0 || count <= 0 {
		return nil, errorIllegalParams
	}
	changeResult := &ChangeResult{
		Succ: false,
	}
	maxOverlap := b.configProvider.GetConfig(itemId).GetMaxOverlap()
	remaining := count
	// 先尝试往已有物品堆叠
	for _, item := range b.Items {
		if item.ItemId == itemId && (maxOverlap == 0 || item.Count < maxOverlap) {
			canAdd := remaining 
			if maxOverlap > 0 {
				canAdd = min(remaining, maxOverlap-item.Count)
			}
			fromNum := item.Count
			currNum := item.ChangeAmount(canAdd)
			remaining -= canAdd
			remaining -= canAdd
			changeResult.addChanged(item, fromNum, currNum, canAdd)
			if remaining == 0 {
				break
			}
		}
	}
	 // 若还有剩余，创建新物品
	 for remaining > 0 {
		newItemAmount := remaining
		if (maxOverlap == 0) {
			newItemAmount = remaining
		} else{
			newItemAmount = min(remaining, maxOverlap)
		}
		remaining -= newItemAmount
		newItem := &Item{
			Uid:    util.GetNextID(),
			ItemId: itemId,
			Count:  newItemAmount,
			Level:  0,
		}
		if initFunc != nil {
			initFunc(newItem)
		}
		b.Items[newItem.Uid] = newItem
		changeResult.addChanged(newItem, 0, newItemAmount, newItemAmount)
	 }
	return changeResult, nil
}

/// 通过配置模型ID减少道具
func (b *Backpack) ReduceByModelId(itemId int32, count int32) *ChangeResult {
	toRemove := count
	hasNum := b.GetItemCount(itemId)
	result :=  &ChangeResult{
		Succ: false,
	}
	if hasNum < toRemove {
		return result
	}
	
	for _, item := range b.Items {
		if item.ItemId == itemId {
			fromNum := item.Count
			if item.Count >= toRemove {
				item.ChangeAmount(-toRemove)
				curr := item.Count
				if (curr == 0) {
					delete(b.Items, item.Uid)
				}
				result.addChanged(item, fromNum, curr, -toRemove)
			} else {
				curr := item.Count
				toRemove -= item.Count
				item.ChangeAmount(-item.Count)
				delete(b.Items, item.Uid)
				result.addChanged(item, fromNum, 0, -curr)
			}
		}
	}

	return result
}

/// 通过道具UID减少道具
func (b *Backpack) ReduceByUid(uid string, count int32) (*ChangeResult, error) {
	if count <= 0 {
		return nil, errorIllegalParams
	}
	existed := b.GetItemByUid(uid)
	result := &ChangeResult{
		Succ: false,
	}
	if existed == nil || existed.Count < count {
		return result, nil
	}
	fromNum := existed.Count
	currNum := existed.ChangeAmount(-count)
	if currNum == 0 {
		delete(b.Items, uid)
	}
	result.addChanged(existed, fromNum, currNum, -count)
	return result, nil
}

/// 通过道具UID获取道具
func (b *Backpack) GetItemByUid(uid string) *Item {
	item, ok := b.Items[uid]
	if ok {
		return item
	}
	return nil
}

/// 获取道具数量
func (b *Backpack) GetItemCount(itemId int32) int32 {
	sum := int32(0)
	for _, item := range b.Items {
		if item.ItemId == itemId {
			sum += item.Count
		}
	}
	return sum
}

func (b *Backpack) IsEnough(itemId int32, count int32) bool {
	cost := map[int32]int32{
		itemId: count,
	}
	return b.IsEnough2(cost)
}

func (b *Backpack) IsEnough2(cost map[int32]int32) bool {
	owned := make(map[int32]int32)
	for _, item := range b.Items {
		itemId := item.ItemId
		if _,ok := cost[itemId];ok {
			prev := owned[itemId]
			owned[itemId] = prev + item.Count
		}
	}
	for itemId, count := range cost {
		if owned[itemId] < count {
			return false
		}
	}
	return true
}

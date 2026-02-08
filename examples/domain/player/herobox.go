package player

type HeroBox struct {
	// 累计招募次数
	RecruitSum int32
	// 英雄列表
	Heros map[int32]*Hero
	// 英雄升级次数
	UpLevelTimes int32
}

func (h *HeroBox) AfterLoad() {
	if h.Heros == nil {
		h.Heros = make(map[int32]*Hero)
	}
}

func (h *HeroBox) AddHero(hero *Hero) {
	h.Heros[hero.ModelId] = hero
}

func (h *HeroBox) GetHero(modelId int32) *Hero {
	return h.Heros[modelId]
}

func (h *HeroBox) GetAllHeros() []*Hero {
	heros := make([]*Hero, 0, len(h.Heros))
	for _, hero := range h.Heros {
		heros = append(heros, hero)
	}
	return heros
}

func (h *HeroBox) GetUpFightHeros() []*Hero {
	heros := make([]*Hero, 0, len(h.Heros))
	for _, hero := range h.Heros {
		if hero.Position > 0 {
			heros = append(heros, hero)
		}
	}
	return heros
}

func (h *HeroBox) HasHero(modelId int32) bool {
	_, ok := h.Heros[modelId]
	return ok
}

func (h *HeroBox) GetEmpostPos() []int32 {
	pos := make([]int32, 0, 5)
	used := make(map[int32]bool)
	for _, hero := range h.Heros {
		if hero.Position > 0 {
			used[hero.ModelId] = true
		}
	}
	for i := int32(1); i <= 5; i++ {
		if !used[i] {
			pos = append(pos, i)
		}
	}
	return pos
}

func (h *HeroBox) GetHeroByPosition(position int32) *Hero {
	for _, hero := range h.Heros {
		if hero.Position == position {
			return hero
		}
	}
	return nil
}

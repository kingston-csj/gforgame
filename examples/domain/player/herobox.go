package player

type HeroBox struct {
	// 累计招募次数
	RecruitSum int32
	// 英雄列表
	Heros map[int32]*Hero
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

func (h *HeroBox) HasHero(modelId int32) bool {
	_, ok := h.Heros[modelId]	
	return ok
}

package skill

type CommonSkill struct {
	baseSkill
}

func (s *CommonSkill) GetEffectType() int32 {
	return 1
}

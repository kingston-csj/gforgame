package report

// 同一个技能释放的效果，同时播放
type SkillReport struct {
	Attacker    string       `json:"attacker"`
	SkillId     int32        `json:"skill_id"`
	DamageUnits []DamageUnit `json:"damage_units"`
}

func (r *SkillReport) AddDamageUnit(damageUnit *DamageUnit) {
	r.DamageUnits = append(r.DamageUnits, *damageUnit)
}

func NewSkillReport(attacker string, skillId int32) *SkillReport {
	return &SkillReport{Attacker: attacker, SkillId: skillId, DamageUnits: make([]DamageUnit, 0)}
}

package report

type DamageUnit struct {
	attacker string
	defender string
	damage   int32
}

func NewDamageUnit(attacker string, defender string, damage int32) *DamageUnit {
	return &DamageUnit{attacker: attacker, defender: defender, damage: damage}
}

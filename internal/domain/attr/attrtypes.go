package attr

type AttrType string

const (
	Hp               AttrType = "hp"
	Attack           AttrType = "attack"
	Defense          AttrType = "defense"
	Speed            AttrType = "speed"
	Move_speed       AttrType = "move_speed"
	Attack_range     AttrType = "attack_range"
	Attack_speed_add AttrType = "attack_speed_add"
	Cd               AttrType = "cd"
	Hp_add           AttrType = "hp_add"
	Attack_add       AttrType = "attack_add"
	Defense_add      AttrType = "defense_add"
	Crit_probability AttrType = "crit_probability"
	Crit_hit         AttrType = "crit_hit"
	Penetration      AttrType = "penetration"
)

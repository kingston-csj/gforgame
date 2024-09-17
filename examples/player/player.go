package player

type Player struct {
	Id     string `json:"id"`
	Delete bool
	Name   string `gorm:"player's name"`
	Level  uint   `gorm:"player's' level"`
}

func (s *Player) GetId() string {
	return s.Id
}

func (s *Player) IsDeleted() bool {
	return s.Delete
}

func (s *Player) SetDeleted() {
	s.Delete = true
}

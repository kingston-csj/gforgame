package player

type Player struct {
	Id    int64  `gorm:"primarykey"`
	Name  string `gorm:"player's name"`
	Level uint   `gorm:"player's' level"`
}

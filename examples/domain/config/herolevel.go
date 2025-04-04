package config

type HeroLevelData struct {
	Id    int32 `json:"id" excel:"id"`
	Level int32 `json:"level" excel:"level"`
	Cost  int32 `json:"cost" excel:"cost"`
	Hp    int32 `json:"hp" excel:"hp"`
	Attack   int32 `json:"attack" excel:"attack"`
	Defense  int32 `json:"defense" excel:"defense"`
	Speed    int32 `json:"speed" excel:"speed"`
}

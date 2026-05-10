package config

type NameData struct {
	Id    int32  `json:"id" excel:"id"`
	First string `json:"first" excel:"first"`
	Last  string `json:"last" excel:"last"`
}
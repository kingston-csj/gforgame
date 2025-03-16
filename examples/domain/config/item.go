package config

type ItemData struct {
	Id      int64  `json:"id" excel:"id"`
	Name    string `json:"name" excel:"name"`
	Quality int64  `json:"quality" excel:"quality"`
	Tips    string `json:"tips" excel:"tips"`
	Icon    string `json:"icon" excel:"icon"`
}

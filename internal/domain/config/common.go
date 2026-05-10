package config

type CommonData struct {
	Id    int32  `json:"id" excel:"id"`
	Key   string `json:"key" excel:"key"`
	Value string `json:"value" excel:"value"`
}

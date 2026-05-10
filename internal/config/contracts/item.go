package config

type ItemConfig interface {
	GetMaxOverlap() int32
	GetId() int32
}

type ItemConfigProvider interface {
	GetConfig(itemId int32) ItemConfig
}

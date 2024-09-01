package network

// Module is the interface that represent a mudule, eg skill, shop, quest
type Module interface {
	Init()

	Shutdown()
}

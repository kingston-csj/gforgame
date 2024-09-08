package network

// Module is the interface that represent a module, eg skill, shop, quest
type Module interface {
	Init()

	Shutdown()
}

// Base implements for Module.
type Base struct{}

// Init will be invoked when app starts
func (c Base) Init() {

}

// Shutdown will be invoked before app closes
func (c Base) Shutdown() {}

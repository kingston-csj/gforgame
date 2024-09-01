package network

var router []Module

func RegisterModule(m Module) {
	router = append(router, m)
}

func ListModules() []Module {
	return router
}

package iface

type ServerInterface interface {
	Start()
	Stop()
	Serve()
}

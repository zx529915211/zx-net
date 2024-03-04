package iface

type ServerInterface interface {
	Start()
	Stop()
	Serve()

	//路由功能 给当前的服务注册一个路由 供客户端的链接使用
	AddRouter(router RouterInterface)
}

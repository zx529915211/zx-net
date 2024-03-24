package iface

type ServerInterface interface {
	Start()
	Stop()
	Serve()

	// AddRouter 路由功能 给当前的服务注册一个路由 供客户端的链接使用
	AddRouter(msgId uint32, router RouterInterface)
	GetConnManager() ConnManagerInterface
	SetOnConnStart(func(conn ConnectionInterface))
	SetOnConnStop(func(conn ConnectionInterface))
	CallOnConnStart(conn ConnectionInterface)
	CallOnConnStop(conn ConnectionInterface)
}

//AddMsgHandle(msgHandle MsgHandleInterface)

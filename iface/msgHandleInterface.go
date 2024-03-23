package iface

type MsgHandleInterface interface {
	//调用对应的Router消息处理方法
	DoMsgHandle(request RequestInterface)

	//给消息添加具体的处理逻辑
	AddRouter(msgId uint32, router RouterInterface)

	StartWorkerPool()

	SendMsgToTaskQueue(request RequestInterface)
}

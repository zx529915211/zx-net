package net

import (
	"strconv"
	"zx-net/iface"
	"zx-net/utils"
)

type MsgHandle struct {
	//存放每个msgId对应的处理方法
	Apis map[uint32]iface.RouterInterface
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{Apis: make(map[uint32]iface.RouterInterface)}
}

func (m *MsgHandle) DoMsgHandle(request iface.RequestInterface) {
	handle, ok := m.Apis[request.GetMsgId()]
	if !ok {
		request.GetConnection().SendMsg(0, []byte("msgId not found"))
		utils.LogErrorString("msgId ", strconv.Itoa(int(request.GetMsgId())), " not found")
		return
	}
	handle.BeforeHandle(request)
	handle.Handle(request)
	handle.AfterHandle(request)
}

// 为不同的消息添加具体的处理逻辑
func (m *MsgHandle) AddRouter(msgId uint32, router iface.RouterInterface) {
	if _, ok := m.Apis[msgId]; ok {
		utils.LogErrorString("repeat api,msgId = ", strconv.Itoa(int(msgId)))
		panic("add router panic")
	}
	m.Apis[msgId] = router

}

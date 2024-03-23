package net

import (
	"fmt"
	"strconv"
	"zx-net/iface"
	"zx-net/utils"
)

type MsgHandle struct {
	//存放每个msgId对应的处理方法
	Apis map[uint32]iface.RouterInterface
	//负责Worker取任务的消息队列
	TaskQueue []chan iface.RequestInterface
	//业务工作Worker池的worker数量
	WorkerPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]iface.RouterInterface),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan iface.RequestInterface, utils.GlobalObject.WorkerPoolSize),
	}
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

// StartWorkerPool 启动一个Worker工作池
func (m *MsgHandle) StartWorkerPool() {
	//根据workerPoolSize 分别开启Worker,每个worker用一个go来承载
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		//启动一个worker
		//1.给worker对应的channel消息队列开辟空间 第i个worker就用第i个channel
		m.TaskQueue[i] = make(chan iface.RequestInterface, utils.GlobalObject.MaxWorkerTaskLen)
		//2.启动当前的worker 阻塞等待消息从channel传递进来
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

// 启动一个worker工作流程
func (m MsgHandle) StartOneWorker(wokerId int, taskQueue chan iface.RequestInterface) {
	fmt.Println("Worker Id=", wokerId, "is started...")

	for {
		select {
		case request := <-taskQueue:
			m.DoMsgHandle(request)
			fmt.Println("workId = ", wokerId, " do msg handle")
			m.DoMsgHandle(request)
		}
	}
}

func (m *MsgHandle) SendMsgToTaskQueue(request iface.RequestInterface) {
	//将消息平均分配给不同的worker
	//分配的worker由连接id取余
	workerId := request.GetConnection().GetConnID() % m.WorkerPoolSize

	//将消息发送给对应的worker的TaskQueue
	m.TaskQueue[workerId] <- request
}

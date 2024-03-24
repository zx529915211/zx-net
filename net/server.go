package net

import (
	"fmt"
	"net"
	"zx-net/iface"
	"zx-net/utils"
)

type Server struct {
	Name        string
	IPVersion   string
	Ip          string
	Port        int
	MsgHandle   iface.MsgHandleInterface
	ConnManager iface.ConnManagerInterface
	OnConnStart func(conn iface.ConnectionInterface)
	OnConnStop  func(conn iface.ConnectionInterface)
}

//func (s *Server) AddMsgHandle(msgHandle iface.MsgHandleInterface) {
//	s.MsgHandle = msgHandle
//}

func (s *Server) AddRouter(msgId uint32, router iface.RouterInterface) {
	s.MsgHandle.AddRouter(msgId, router)
}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("try to start %s:%d server\n", s.Ip, s.Port)
	go func() {
		//开启工作池
		s.MsgHandle.StartWorkerPool()
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
		if err != nil {
			fmt.Println("create tcp addr err:", err)
			return
		}
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listenTcp err:", err)
			return
		}
		var cid uint32
		cid = 0
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("acceptTcp err:", err)
				continue
			}
			//最大连接个数判断
			if s.ConnManager.Len() >= utils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}

			connection := NewConnection(s, conn, cid, s.MsgHandle)
			cid++
			go connection.Start()
		}
	}()

}

// 停止服务器
func (s *Server) Stop() {
	fmt.Println("server stop")
	s.ConnManager.Clear()

}

// 运行服务器
func (s *Server) Serve() {
	s.Start()

	//阻塞
	select {}
}

func (s *Server) GetConnManager() iface.ConnManagerInterface {
	return s.ConnManager
}

func NewServer() iface.ServerInterface {
	return &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		Ip:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		MsgHandle:   NewMsgHandle(),
		ConnManager: NewConnManager(),
	}
}

func (s *Server) SetOnConnStart(hookFunc func(conn iface.ConnectionInterface)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(conn iface.ConnectionInterface)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(conn iface.ConnectionInterface) {
	if s.OnConnStart != nil {
		fmt.Println("CallOnConnStart.. ")
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn iface.ConnectionInterface) {
	if s.OnConnStop != nil {
		fmt.Println("CallOnConnStop.. ")
		s.OnConnStop(conn)
	}
}

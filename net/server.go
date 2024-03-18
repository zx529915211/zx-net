package net

import (
	"fmt"
	"net"
	"zx-net/iface"
	"zx-net/utils"
)

type Server struct {
	Name      string
	IPVersion string
	Ip        string
	Port      int
	MsgHandle iface.MsgHandleInterface
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
			connection := NewConnection(conn, cid, s.MsgHandle)
			cid++
			go connection.Start()
		}
	}()

}

// 停止服务器
func (s *Server) Stop() {

}

// 运行服务器
func (s *Server) Serve() {
	s.Start()

	//阻塞
	select {}
}

func NewServer() iface.ServerInterface {
	return &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		Ip:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		MsgHandle: NewMsgHandle(),
	}
}

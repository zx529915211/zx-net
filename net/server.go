package net

import (
	"errors"
	"fmt"
	"net"
	"zx-net/iface"
)

type Server struct {
	Name      string
	IPVersion string
	Ip        string
	Port      int
}

// 定义当前客户端链接绑定的handle api 目前写死了回显，后续可以由用固话增加协议解析等回调
func CallBackToClient(conn *net.TCPConn, data []byte, msgLen int) error {
	fmt.Println("callback to client")
	_, err := conn.Write(data[:msgLen])
	if err != nil {
		fmt.Println("write back err:", err)
		return errors.New("callback to client error")
	}
	return nil
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

		fmt.Println("server start success")

		var cid uint32
		cid = 0

		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("acceptTcp err:", err)
				continue
			}

			connection := NewConnection(conn, cid, CallBackToClient)

			cid++

			go connection.Start()
			//go func() {
			//	for {
			//		buf := make([]byte, 512)
			//		readLen, err := conn.Read(buf)
			//		if err != nil {
			//			fmt.Println("read buf err:", err)
			//			continue
			//		}
			//
			//		if _, err := conn.Write(buf[:readLen]); err != nil {
			//			fmt.Println("write buf err:", err)
			//			continue
			//		}
			//	}
			//}()
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

func NewServer(name string) iface.ServerInterface {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		Ip:        "0.0.0.0",
		Port:      8888,
	}
}

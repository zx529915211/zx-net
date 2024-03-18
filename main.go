package main

import (
	"fmt"
	"zx-net/iface"
	"zx-net/net"
)

type PingRouter struct {
	net.BaseRouter
}

//func (p *PingRouter) BeforeHandle(request iface.RequestInterface) {
//	fmt.Println("before Router")
//	_, err := request.GetConnection().GetTcpConnection().Write([]byte("before ping..."))
//	if err != nil {
//		fmt.Println("call back before ping error")
//	}
//}

func (p *PingRouter) Handle(request iface.RequestInterface) {
	fmt.Println("recv from client: msgID = ", request.GetMsgId(), ", data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(0, []byte("ping..ping.ping"))
	if err != nil {
		fmt.Println(err)
	}
}

//	func (p *PingRouter) AfterHandle(request iface.RequestInterface) {
//		fmt.Println("call after Router")
//		_, err := request.GetConnection().GetTcpConnection().Write([]byte("after ping..."))
//		if err != nil {
//			fmt.Println("call back after ping error")
//		}
//	}
type HelloRouter struct {
	net.BaseRouter
}

func (p *HelloRouter) Handle(request iface.RequestInterface) {
	fmt.Println("recv from client: msgID = ", request.GetMsgId(), ", data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(1, []byte("hello..hello.hello"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	s := net.NewServer()

	//msgHandle := net.NewMsgHandle()
	//msgHandle.AddRouter(1, &PingRouter{})
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	s.Serve()
}

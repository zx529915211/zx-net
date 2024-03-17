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
	fmt.Println("call Router")
	//_, err := request.GetConnection().GetTcpConnection().Write([]byte("ping..."))
	//if err != nil {
	//	fmt.Println("call back ping error")
	//}
	fmt.Println("recv from client: msgID = ", request.GetMsgId(), ", data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("ping..ping.ping"))
	if err != nil {
		fmt.Println(err)
	}
}

//func (p *PingRouter) AfterHandle(request iface.RequestInterface) {
//	fmt.Println("call after Router")
//	_, err := request.GetConnection().GetTcpConnection().Write([]byte("after ping..."))
//	if err != nil {
//		fmt.Println("call back after ping error")
//	}
//}

func main() {
	s := net.NewServer()

	s.AddRouter(&PingRouter{})

	s.Serve()
}

package iface

import "net"

type ConnectionInterface interface {
	Start()

	Stop()

	GetTcpConnection() *net.TCPConn

	GetConnID() uint32

	RemoteAddr() net.Addr

	SendMsg(msgId uint32, data []byte) error
}

// 处理链接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error

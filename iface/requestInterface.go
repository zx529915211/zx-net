package iface

type RequestInterface interface {
	GetConnection() ConnectionInterface

	GetData() []byte

	GetMsgId() uint32
}

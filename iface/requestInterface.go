package iface

type RequestInterface interface {
	GetConnection() ConnectionInterface

	GetData() []byte
}

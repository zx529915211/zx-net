package iface

type ConnManagerInterface interface {
	Add(connectionInterface ConnectionInterface)

	Remove(connectionInterface ConnectionInterface)

	Get(connId uint32) (ConnectionInterface, error)

	Len() int

	Clear()
}

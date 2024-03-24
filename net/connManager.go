package net

import (
	"errors"
	"fmt"
	"sync"
	"zx-net/iface"
)

type ConnManager struct {
	Connections map[uint32]iface.ConnectionInterface
	connLock    sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		Connections: make(map[uint32]iface.ConnectionInterface),
	}
}

func (c *ConnManager) Add(conn iface.ConnectionInterface) {
	c.connLock.Lock()
	c.Connections[conn.GetConnID()] = conn

	c.connLock.Unlock()
}

func (c *ConnManager) Remove(conn iface.ConnectionInterface) {
	c.connLock.Lock()
	delete(c.Connections, conn.GetConnID())
	c.connLock.Unlock()
}

func (c *ConnManager) Get(connId uint32) (iface.ConnectionInterface, error) {
	c.connLock.RLock()
	defer c.connLock.RUnlock()
	if conn, ok := c.Connections[connId]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

func (c *ConnManager) Len() int {
	return len(c.Connections)
}

func (c *ConnManager) Clear() {
	c.connLock.RLock()
	defer c.connLock.RUnlock()
	for connId, conn := range c.Connections {
		conn.Stop()
		delete(c.Connections, connId)
	}
	fmt.Println("clear all connection")
}

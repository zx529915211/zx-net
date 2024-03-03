package net

import (
	"fmt"
	"net"
	"zx-net/iface"
)

// 链接类
type Connection struct {
	Conn      *net.TCPConn
	ConnID    uint32
	isClosed  bool
	handleAPI iface.HandleFunc
	ExitChan  chan bool
}

func NewConnection(conn *net.TCPConn, connId uint32, callbackApi iface.HandleFunc) *Connection {
	return &Connection{
		Conn:      conn,
		ConnID:    connId,
		handleAPI: callbackApi,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
	}
}

// 链接的读业务
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")

	defer fmt.Println("connId = ", c.ConnID, "reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		buf := make([]byte, 512)

		msgLength, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err", err)
			continue
		}

		//调用当前链接所绑定的HandleAPI
		if err := c.handleAPI(c.Conn, buf, msgLength); err != nil {
			fmt.Println("connId", c.ConnID, "handle is error", err)
			break
		}
	}
}

// 启动链接
func (c *Connection) Start() {
	fmt.Println("Conn start ,connId = ", c.ConnID)

	//启动从当前链接读数据的业务
	go c.StartReader()
}

func (c *Connection) Stop() {
	fmt.Println("conn stop().. connid=", c.ConnID)

	if c.isClosed {
		return
	}

	c.isClosed = true

	c.Conn.Close()

	close(c.ExitChan)
}

func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	c.Conn.Write(data[:])
	return nil
}

package net

import (
	"fmt"
	"net"
	"zx-net/iface"
	"zx-net/utils"
)

// 链接类
type Connection struct {
	Conn     *net.TCPConn
	ConnID   uint32
	isClosed bool
	ExitChan chan bool
	Router   iface.RouterInterface
}

func NewConnection(conn *net.TCPConn, connId uint32, router iface.RouterInterface) *Connection {
	return &Connection{
		Conn:     conn,
		ConnID:   connId,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}
}

// 链接的读业务
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")

	defer fmt.Println("connId = ", c.ConnID, "reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		buf := make([]byte, utils.GlobalObject.MaxPackageSize)

		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err", err)
			continue
		}

		request := Request{
			conn: c,
			data: buf,
		}

		//执行注册的路由方法
		go func(r iface.RequestInterface) {
			c.Router.BeforeHandle(r)
			c.Router.Handle(r)
			c.Router.AfterHandle(r)
		}(&request)

		//调用当前链接所绑定的HandleAPI
		//if err := c.handleAPI(c.Conn, buf, msgLength); err != nil {
		//	fmt.Println("connId", c.ConnID, "handle is error", err)
		//	break
		//}
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

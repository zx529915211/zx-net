package net

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zx-net/iface"
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
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf err", err)
		//	continue
		//}
		//创建拆包解包对象
		dp := NewDataPack()

		//读取客户端的Msg Head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		n, err := io.ReadFull(c.GetTcpConnection(), headData)
		if (uint32(n)) != dp.GetHeadLen() || err != nil {
			fmt.Println("read msg head error", err)
			break
		}
		//拆包 得到dataLen和msgId
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}

		//根据dataLen 再次读取Data
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			_, err := io.ReadFull(c.GetTcpConnection(), data)
			if err != nil {
				fmt.Println("read msg error", err)
			}
		}
		msg.SetMsgData(data)

		request := Request{
			conn: c,
			msg:  msg,
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

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	//data封包
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))

	if err != nil {
		fmt.Println("Pack error msg id =", msgId)
		return errors.New("Pack error msg")
	}

	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("Write msg id ", msgId, " error :", err)
		return errors.New("conn Write error")
	}

	return nil
}

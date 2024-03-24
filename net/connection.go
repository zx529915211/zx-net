package net

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zx-net/iface"
	"zx-net/utils"
)

// 链接类
type Connection struct {
	TcpServer iface.ServerInterface
	Conn      *net.TCPConn
	ConnID    uint32
	isClosed  bool
	ExitChan  chan bool
	msgChan   chan []byte
	MsgHandle iface.MsgHandleInterface
}

func NewConnection(server iface.ServerInterface, conn *net.TCPConn, connId uint32, msgHandle iface.MsgHandleInterface) *Connection {
	c := &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnID:    connId,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
		msgChan:   make(chan []byte),
		MsgHandle: msgHandle,
	}
	c.TcpServer.GetConnManager().Add(c)
	return c
}

// 链接的读业务
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")

	defer fmt.Println("connId = ", c.ConnID, "【reader is exit】, remote addr is ", c.RemoteAddr().String())
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

		if utils.GlobalObject.WorkerPoolSize > 0 {
			//开启了工作池，交给工作池处理
			c.MsgHandle.SendMsgToTaskQueue(&request)
		} else {
			go c.MsgHandle.DoMsgHandle(&request)
		}

	}
}

func (c *Connection) StartWriter() {
	fmt.Println("Writer Goroutine is running...")
	defer fmt.Println("【writer is exit】 remote addr is", c.RemoteAddr().String())
	for {
		select {
		case msg := <-c.msgChan:
			if _, err := c.Conn.Write(msg); err != nil {
				fmt.Println("send data error,", err)
				return
			}
		case <-c.ExitChan:
			return

		}
	}
}

// 启动链接
func (c *Connection) Start() {
	fmt.Println("Conn start ,connId = ", c.ConnID)

	//启动从当前链接读数据的业务
	go c.StartReader()
	//启动写协程，处理读业务最后发送给客户端的消息
	go c.StartWriter()
}

func (c *Connection) Stop() {
	fmt.Println("a conn stop().. connid=", c.ConnID)

	if c.isClosed {
		return
	}

	c.isClosed = true

	c.Conn.Close()
	c.ExitChan <- true
	c.TcpServer.GetConnManager().Remove(c)
	close(c.ExitChan)
	close(c.msgChan)
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

	//if _, err := c.Conn.Write(binaryMsg); err != nil {
	//	fmt.Println("Write msg id ", msgId, " error :", err)
	//	return errors.New("conn Write error")
	//}
	c.msgChan <- binaryMsg

	return nil
}

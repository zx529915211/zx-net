package net

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {
	//模拟一个服务器
	listen, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}
	go func() {
		//从客户端读数据
		for {
			conn, err := listen.Accept()
			if err != nil {
				fmt.Println("Accept error", err)
			}
			go func(conn net.Conn) {
				pack := NewDataPack()
				for {
					//把包的head读出来
					headData := make([]byte, pack.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error")
						return
					}
					msgRes, err := pack.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err", err)
					}
					if msgRes.GetMsgLen() > 0 {
						//有数据，第二次读取
						msg := msgRes.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err:", err)
							return
						}
						//完整消息读取完毕
						fmt.Printf("成功接收一个完整的消息，消息长度：%d,消息类型：%d,消息内容:%s\n", msg.DataLen, msg.Id, msg.Data)
					}
				}
			}(conn)
		}
	}()

	//模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}

	dp := NewDataPack()

	//模拟粘包，封装两个msg一起发送
	msg1 := &Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte{'g', 'z', 'f', 'l', 'f'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error", err)
		return
	}

	msg2 := &Message{
		Id:      1,
		DataLen: 7,
		Data:    []byte{'g', 'z', 'a', 'i', 'f', 'l', 'f'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg1 error", err)
		return
	}

	//两个包黏在一起
	sendData1 = append(sendData1, sendData2...)

	//一次性发给服务端
	conn.Write(sendData1)

	select {}
}

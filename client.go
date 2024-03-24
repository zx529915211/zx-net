package main

import (
	"fmt"
	"io"
	"net"
	"zx-net/iface"
	net2 "zx-net/net"
	"zx-net/utils"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}
	dp := net2.NewDataPack()
	data := ReadAndParse(conn, dp)
	if data == nil {
		return
	}
	for {
		//发送封包的Message
		binaryMsg, err := dp.Pack(net2.NewMsgPackage(1, []byte("gzf566")))
		if err != nil {
			utils.LogError("dp pack", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			utils.LogError("conn write", err)
			return
		}

		data := ReadAndParse(conn, dp)
		if data == nil {
			break
		}

	}
}

func ReadAndParse(conn net.Conn, dp iface.DataPackInterface) []byte {
	headData := make([]byte, dp.GetHeadLen())
	n, err := io.ReadFull(conn, headData)
	if (uint32(n)) != dp.GetHeadLen() || err != nil {
		utils.LogError("read msg head", err)
		return nil
	}
	msg, _ := dp.Unpack(headData)
	headLen := msg.GetMsgLen()
	data := make([]byte, headLen)
	if headLen > 0 {
		io.ReadFull(conn, data)
		fmt.Println("data:", string(data))
	}
	return data
}

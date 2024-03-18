package main

import (
	"fmt"
	"io"
	"net"
	net2 "zx-net/net"
	"zx-net/utils"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		//发送封包的Message
		dp := net2.NewDataPack()
		binaryMsg, err := dp.Pack(net2.NewMsgPackage(1, []byte("gzf566")))
		if err != nil {
			utils.LogError("dp pack", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			utils.LogError("conn write", err)
			return
		}

		headData := make([]byte, dp.GetHeadLen())
		n, err := io.ReadFull(conn, headData)
		if (uint32(n)) != dp.GetHeadLen() || err != nil {
			utils.LogError("read msg head", err)
			break
		}
		msg, _ := dp.Unpack(headData)
		headLen := msg.GetMsgLen()

		if headLen > 0 {
			data := make([]byte, headLen)
			io.ReadFull(conn, data)
			fmt.Println("data:", string(data))
		}

	}
}

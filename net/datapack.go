package net

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zx-net/iface"
	"zx-net/utils"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包头长度
func (d *DataPack) GetHeadLen() uint32 {
	//Datalen uint32(4字节)+ Id uint31(4字节)
	return 8
}

// 封包
func (d *DataPack) Pack(msg iface.MessageInterface) (data []byte, err error) {
	//创建一个存放bytes字节的缓冲
	dataBuf := bytes.NewBuffer([]byte{})

	//将dataLen写金databuf
	err = binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgLen())
	if err != nil {
		return
	}

	err = binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		return
	}

	err = binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgData())
	if err != nil {
		return
	}
	data = dataBuf.Bytes()
	return
}

// 拆包
func (d *DataPack) Unpack(binaryData []byte) (iface.MessageInterface, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuf := bytes.NewReader(binaryData)

	msg := &Message{}

	//读取dataLen
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读msgId
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断datalen是否已经超过了允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("包的数据长度超过最大值")
	}
	return msg, nil
}

package iface

/*
封包拆包模块
直接处理TCP的数据流，用于处理TCP粘包的问题
*/
type DataPackInterface interface {
	// GetHeadLen 获取包头长度
	GetHeadLen() uint32
	// Pack 封包
	Pack(msg MessageInterface) ([]byte, error)
	// Unpack 拆包
	Unpack([]byte) (MessageInterface, error)
}

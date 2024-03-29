package iface

type MessageInterface interface {
	GetMsgId() uint32
	GetMsgLen() uint32
	GetMsgData() []byte

	SetMsgId(uint32)
	SetMsgLen(uint32)
	SetMsgData([]byte)
}

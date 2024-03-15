package net

type Message struct {
	DataLen uint32
	Id      uint32
	Data    []byte
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}
func (m *Message) GetMsgData() []byte {
	return m.Data
}

func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}
func (m *Message) SetMsgLen(len uint32) {
	m.DataLen = len
}
func (m *Message) SetMsgData(data []byte) {
	m.Data = data
}

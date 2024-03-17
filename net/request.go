package net

import "zx-net/iface"

type Request struct {
	//已经和客户端简历好的链接
	conn iface.ConnectionInterface

	//客户端请求的数据
	msg iface.MessageInterface
}

func (r *Request) GetConnection() iface.ConnectionInterface {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetMsgData()
}

func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}

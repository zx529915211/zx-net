package net

import "zx-net/iface"

type Request struct {
	//已经和客户端简历好的链接
	conn iface.ConnectionInterface

	//客户端请求的数据
	data []byte
}

func (r *Request) GetConnection() iface.ConnectionInterface {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.data
}

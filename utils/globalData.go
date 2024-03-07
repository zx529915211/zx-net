package utils

import (
	"encoding/json"
	"os"
	"zx-net/iface"
)

type GlobalData struct {
	TcpServer iface.ServerInterface //全局Server对象
	Host      string
	TcpPort   int
	Name      string

	Version        string
	MaxConn        int    //主机允许的最大链接数
	MaxPackageSize uint32 //数据包最大值
}

var GlobalObject *GlobalData

func (g *GlobalData) Reload() {
	data, err := os.ReadFile("conf/config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	GlobalObject = &GlobalData{
		TcpServer:      nil,
		Host:           "0.0.0.0",
		TcpPort:        8888,
		Name:           "zx-net",
		Version:        "v0.1",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}
	GlobalObject.Reload()
}

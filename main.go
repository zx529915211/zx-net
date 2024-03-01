package main

import (
	"zx-net/net"
)

func main() {
	s := net.NewServer("zx-net")
	s.Serve()
}

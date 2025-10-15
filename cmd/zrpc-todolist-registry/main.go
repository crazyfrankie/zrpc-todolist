package main

import (
	"flag"
	"log"

	"github.com/crazyfrankie/zrpc/registry"
)

var (
	addr      string
	keepAlive int
)

func main() {
	flag.StringVar(&addr, "addr", "127.0.0.1:9000", "service registry")
	flag.IntVar(&keepAlive, "keepalive", 120, "service registry keepalive time")

	flag.Parse()

	tcpRegistry, err := registry.NewTcpRegistry(addr, keepAlive)
	if err != nil {
		panic(err)
	}
	log.Printf("tcp registry start at %s\n", addr)

	tcpRegistry.Serve()
}

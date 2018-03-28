package main

import (
	"flag"
	"fmt"
	"net"
	"redis_go/log"
	"redis_go/server"
)

var flags struct {
	addr string
}

func init() {
	flag.StringVar(&flags.addr, "addr", ":9736", "The TCP address to bing to ")
}

func run() error {
	server := server.NewServer(nil)
	lis, err := net.Listen("tcp", flags.addr)
	if err != nil {
		return err
	}
	defer lis.Close()
	log.Info("waiting for connections on %s", lis.Addr().String())

	return server.Serve(lis)
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatal(fmt.Sprintf("start server error %+v", err))
	}
}

package main

import (
	"flag"
	"fmt"
	"net"
	"redis_go/conf"
	"redis_go/loggers"
	"redis_go/server"
	"runtime"
)

var flags struct {
	addr     string
	logLevel int64
}

func init() {
	flag.StringVar(&flags.addr, "addr", fmt.Sprintf("%s:%d", conf.RedisServerAddr, conf.RedisServerPort), "The TCP address to bing to ")
	//flag.Int64Var(&flags.logLevel, "log-level", conf.RedisLogLevel, "System Log Level")
}

func run() error {
	// TODO lmj 解决这个问题。console输入的参数优先级最高->然后是配置文件里的->然后是系统默认的
	server := server.NewServer(nil)
	lis, err := net.Listen("tcp", flags.addr)
	if err != nil {
		return err
	}
	defer lis.Close()
	loggers.Info("waiting for connections on %s", lis.Addr().String())
	return server.Serve(lis)
}

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(6)
	if err := run(); err != nil {
		loggers.Fatal(fmt.Sprintf("start server error %+v", err))
	}
}

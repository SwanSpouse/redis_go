package server

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/mreiferson/go-options"
	"net"
	"os"
	"redis_go/conf"
	"redis_go/loggers"
)

func redisGoFlagSet(opts *conf.ServerConfig) *flag.FlagSet {
	flagSet := flag.NewFlagSet("redis_go", flag.ExitOnError)

	flagSet.String("config", "", "path to config file")

	flagSet.Int("sentinel-mode", opts.SentinelMode, "is sentinel mode")
	flagSet.Int64("log-level", opts.LogLevel, "System Log Level")
	flagSet.Duration("timeout", opts.Timeout, "timeout")

	flagSet.Int("port", opts.Port, "port")
	flagSet.String("addr", opts.BindAddr, "addr")
	flagSet.Int("reader-pool-size", opts.ReaderPoolSize, "reader-pool-size")
	flagSet.Int("writer-pool-size", opts.WriterPoolSize, "writer-pool-size")

	flagSet.Int("verbosity", opts.Verbosity, "verbosity")
	flagSet.Int64("max-idle-time", opts.MaxIdleTime, "max-idle-time")
	flagSet.Int("db-num", opts.DBNum, "db-num")

	flagSet.Int64("client-max-query-buf-len", opts.ClientMaxQueryBufLen, "client-max-query-buf-len")

	flagSet.Int("aof-state", opts.AofState, "aof switch default off")
	flagSet.String("aof-fsync", opts.AofFSync, "")
	flagSet.String("aof-filename", opts.AofFilename, "")

	flagSet.Int64("dirty", opts.Dirty, "")
	flagSet.Int64("dirty-before-bg-save", opts.DirtyBeforeBgSave, "")
	flagSet.String("rdb-filename", opts.RdbFilename, "")
	return flagSet
}

type Program struct {
	*Server
}

func (p *Program) Init() {
	/**
	配置优先级
		1. Command line flag
		2. Deprecated command line flag
		3. Config file value
		4. Get() value (if Getter)
		5. Options struct default value
	*/

	// 用默认值初始化config
	opts := conf.NewServerConfig()

	flagSet := redisGoFlagSet(opts)
	flagSet.Parse(os.Args[1:])

	// 如果有配置文件则从配置文件中读取
	var cfg map[string]interface{}
	configFile := flagSet.Lookup("config").Value.String()
	if configFile != "" {
		_, err := toml.DecodeFile(configFile, &cfg)
		if err != nil {
			loggers.Fatal("failed to load config file")
		}
	}
	// TODO @lmj 在这里需要对配置文件中的配置进行校验，判断输入是否合法
	options.Resolve(opts, flagSet, cfg)
	loggers.Info("input opts :%+v", opts)

	server := NewServer(opts)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Config.BindAddr, server.Config.Port))
	if err != nil || listener == nil {
		os.Exit(1)
	}
	server.TcpListener = listener

	p.Server = server
}

func (p *Program) InitForMock(inputConfig *conf.ServerConfig) {
	var opts *conf.ServerConfig
	if inputConfig == nil {
		// 用默认值初始化config
		opts = conf.NewServerConfig()
	} else {
		opts = inputConfig
	}
	server := NewServer(opts)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Config.BindAddr, server.Config.Port))
	if err != nil || listener == nil {
		os.Exit(1)
	}
	server.TcpListener = listener

	p.Server = server
}

func (p *Program) Start() {
	// 依次启动后台服务
	p.WaitGroup.Wrap(func() {
		p.processTimeEvents()
	})

	p.WaitGroup.Wrap(func() {
		p.TCPServe()
	})
}

func (p *Program) Stop() {
	if p.TcpListener != nil {
		p.TcpListener.Close()
	}

	// TODO @Lmj 这里要对数据进行持久化

	loggers.Info("REDIS GO: stopping subsystems")
	close(p.ExitChan)
	p.WaitGroup.Wait()
	loggers.Info("REDIS GO: bye")
}

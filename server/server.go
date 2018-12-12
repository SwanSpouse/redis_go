package server

import (
	"net"
	"redis_go/client"
	"redis_go/conf"
	"redis_go/database"
	"redis_go/loggers"
	"redis_go/tcp"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const (
	RedisServerStatusNormal             = 0
	RedisServerStatusRdbSaveInProcess   = 1
	RedisServerStatusRdbBgSaveInProcess = 2
)

// Redis server
type Server struct {
	Config        *conf.ServerConfig
	Databases     []*database.Database /* database*/
	dbIndex       int                  // rdb process current db
	clients       []*client.Client
	FakeClient    *client.Client // used in rdb and aof
	password      string         /* Pass for AUTH command, or NULL */
	commandTable  map[string]*client.Command
	mu            sync.RWMutex
	Status        atomic.Value
	Dirty         int64
	rdbLastSave   time.Time
	aofSelectDBId int
	aofBuf        []byte
	aofLastSave   time.Time
}

func NewServer(config *conf.ServerConfig) *Server {
	if config == nil {
		config = conf.InitServerConfig()
	}
	server := &Server{
		Config:       config,
		commandTable: make(map[string]*client.Command),
	}
	// init general parameters
	server.initServer()
	// init Reader & Writer Sync Pool
	server.initIOPool()
	// init databases
	server.initDB()
	// init commandTable table
	server.populateCommandTable()
	// init time events
	go server.initTimeEvents()
	// load data
	server.loadDataFromDisk()
	loggers.Info("redis server: %+v", server)
	return server
}

// 判断server 此时是否可以对外提供服务
func (srv *Server) isServiceAvailable() bool {
	return srv.Status.Load() == RedisServerStatusNormal || srv.Status.Load() == RedisServerStatusRdbBgSaveInProcess
}

// 处理来自客户端的请求
func (srv *Server) IoLoop(conn net.Conn) {

}

// 启动redis server 并开始监听TCP连接
func (srv *Server) Serve(listener net.Listener) {
	loggers.Errorf("TCP: listening on %s", listener.Addr())
	// start to scan clients
	go srv.scanClients()
	// loop for accepting tcp connection from redis client
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			// ignore temporary errors
			if netError, ok := err.(net.Error); ok && netError.Temporary() {
				loggers.Warn("temporary Accept() failure %s", err)
				runtime.Gosched()
				continue
			}
			break
		}
		go srv.IOLoop(clientConn)
	}
	loggers.Errorf("TCP: closing %s", listener.Addr())
}

func (srv *Server) addClientToServer(c *client.Client) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	// add client to server's clients list
	for i := 0; i < len(srv.clients); i++ {
		if srv.clients[i] == nil || srv.clients[i].Closed {
			srv.clients[i] = c
			return
		}
	}
	// TODO lmj srv.clients 有个数限制
	srv.clients = append(srv.clients, c)
}

func (srv *Server) initServer() {
	srv.Status.Store(RedisServerStatusNormal)
	srv.aofSelectDBId = -1
	if srv.Config.AofState == conf.RedisAofOn {
		srv.aofBuf = make([]byte, 0)
	}
	loggers.Level = srv.Config.LogLevel
}

func (srv *Server) initDB() {
	// add default database
	srv.Databases = make([]*database.Database, srv.Config.DBNum)
	for i := 0; i < srv.Config.DBNum; i++ {
		srv.Databases[i] = database.NewDatabase(i)
	}
}

func (srv *Server) initIOPool() {
	for i := 0; i < srv.Config.ReaderPoolNum; i++ {
		tcp.ReaderPool.Put(tcp.NewBufIoReaderWithoutConn())
	}
	for i := 0; i < srv.Config.WriterPoolNum; i++ {
		tcp.WriterPool.Put(tcp.NewBufIoWriterWithoutConn())
	}
	loggers.Debug("Successful init reader and writer pool. ReaderPoolSize:%d, WriterPoolSize:%d", srv.Config.ReaderPoolNum, srv.Config.WriterPoolNum)
}

func (srv *Server) initTimeEvents() {
	//ticker := time.NewTicker(time.Second)
	//for _ = range ticker.C {
	//	log.Info("TICKER INFO client list length %d", len(srv.clients))
	//	for i, item := range srv.clients {
	//		log.Info("TICKER current client list index %d, info:%+v", i, item)
	//	}
	//}
}

func (srv *Server) getDefaultDB() *database.Database {
	if srv.Databases == nil {
		srv.initDB()
	}
	return srv.Databases[0]
}

func (srv *Server) loadDataFromDisk() {
	startTime := time.Now()
	if srv.Config.AofState == conf.RedisAofOn {
		loggers.Info("redis aof start to load data from disk at %s", startTime.Format("20060102 15:04:05"))
		srv.loadAppendOnlyFile()
	} else {
		loggers.Info("redis rdb start to load data from disk at %s", startTime.Format("20060102 15:04:05"))
		srv.rdbLoad()
	}
}

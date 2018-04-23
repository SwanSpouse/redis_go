package server

import (
	"net"
	"redis_go/client"
	"redis_go/conf"
	"redis_go/database"
	"redis_go/handlers"
	"redis_go/loggers"
	"redis_go/tcp"
	"sync"
)

// Redis server
type Server struct {
	Config    *conf.ServerConfig
	Databases []*database.Database /* database*/
	clients   []*client.Client
	password  string /* Pass for AUTH command, or NULL */
	commands  map[string]handlers.BaseHandler
	mu        sync.RWMutex
}

func NewServer(config *conf.ServerConfig) *Server {
	if config == nil {
		config = conf.InitServerConfig()
	}
	server := &Server{
		Config:   config,
		commands: make(map[string]handlers.BaseHandler),
	}
	loggers.Info("redis server config: %+v", config)
	// init general parameters
	server.initServer()
	// init Reader & Writer Sync Pool
	server.initIOPool()
	// init databases
	server.initDB()
	// init commands table
	server.populateCommandTable()
	// init time events
	go server.initTimeEvents()
	loggers.Info("redis server: %+v", server)
	return server
}

func (srv *Server) Serve(lis net.Listener) error {
	// start to scan clients
	go srv.scanClients()
	// loop for accept tcp client
	for {
		cn, err := lis.Accept()
		if err != nil {
			return err
		}
		c := client.NewClient(cn, srv.getDefaultDB())
		srv.addClientToServer(c)
		loggers.Info("new client %d come in ! from %+v and has been added in server's client list.", c.ID(), cn.RemoteAddr().String())
	}
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
	loggers.Level = srv.Config.LogLevel
}

func (srv *Server) initDB() {
	// add default database
	srv.Databases = make([]*database.Database, srv.Config.DBNum)
	for i := 0; i < srv.Config.DBNum; i++ {
		srv.Databases[i] = database.NewDatabase()
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

// register all command handlers
func (srv *Server) populateCommandTable() {
	connectionHandler := new(handlers.ConnectionHandler)
	stringHandler := new(handlers.StringHandler)
	keyHandler := new(handlers.KeyHandler)

	// key command
	srv.commands[handlers.RedisKeyCommandDel] = keyHandler
	srv.commands[handlers.RedisKeyCommandObject] = keyHandler
	srv.commands[handlers.RedisKeyCommandType] = keyHandler
	srv.commands[handlers.RedisKeyCommandExists] = keyHandler

	// connection command
	srv.commands[handlers.RedisConnectionCommandPing] = connectionHandler
	srv.commands[handlers.RedisConnectionCommandAuth] = connectionHandler

	// string command
	srv.commands[handlers.RedisStringCommandAppend] = stringHandler
	srv.commands[handlers.RedisStringCommandSet] = stringHandler
	srv.commands[handlers.RedisStringCommandGet] = stringHandler
	srv.commands[handlers.RedisStringCommandIncr] = stringHandler
	srv.commands[handlers.RedisStringCommandDecr] = stringHandler
	srv.commands[handlers.RedisStringCommandStrLen] = stringHandler
}

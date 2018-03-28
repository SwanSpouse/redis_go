package server

import (
	"net"
	"redis_go/client"
	"redis_go/conf"
	"redis_go/handlers"
	"redis_go/log"
	"redis_go/redis_database"
	"sync"
)

// Redis server
type Server struct {
	config    *conf.ServerConfig
	commands  map[string]handlers.BaseHandler
	mu        sync.RWMutex
	Databases []*redis_database.Database // database
}

func NewServer(config *conf.ServerConfig) *Server {
	if config == nil {
		config = new(conf.ServerConfig)
	}
	server := &Server{
		config:   config,
		commands: make(map[string]handlers.BaseHandler),
	}
	server.initDB()
	// init commands table
	server.populateCommandTable()
	return server
}

func (srv *Server) serveClient(c *client.Client) {
	defer c.Release()
	// TODO lmj 增加Timeout的判断
	// TODO lmj 除了Timeout的方式，还有什么好的办法能够判断client端是否已经断开连接
	// loop to handle redis command
	for more := true; more; more = c.RequestReader.Buffered() != 0 {
		cmd, err := c.RequestReader.ReadCmd(nil)
		if err != nil {
			c.ResponseWriter.AppendErrorf("read command error %+v", err)
			continue
		}
		/**
		首先判断是否在command table中,
			如果不在command table中,则返回command not found
			如果在command table中，则获取到相应的command handler来进行处理。
		*/
		log.Debug("get command from client %+v", cmd)
		if handler, ok := srv.commands[cmd.GetName()]; ok {
			handler.Process(srv.Databases, c, cmd)
		} else {
			log.Errorf("command not found %s", cmd.GetOriginName())
			c.ResponseWriter.AppendError("command not found")
		}
		if err := c.ResponseWriter.Flush(); err != nil {
			log.Info("response writer flush data error %+v", err)
			return
		}
	}
	log.Debug("No more data for current connection")
}

func (srv *Server) Serve(lis net.Listener) error {
	for {
		cn, err := lis.Accept()
		if err != nil {
			return err
		}
		go srv.serveClient(client.NewClient(cn))
		log.Info("new client come in ! from %+v", cn.RemoteAddr().String())
	}
}

func (srv *Server) initDB() {
	srv.Databases = make([]*redis_database.Database, 0)
	// add default database
	srv.Databases = append(srv.Databases, redis_database.NewDatabase())
}

// register all command handlers
func (srv *Server) populateCommandTable() {
	connectionHandler := new(handlers.ConnectionHandler)
	stringHandler := new(handlers.StringHandler)
	srv.commands["PING"] = connectionHandler
	srv.commands["TEST"] = connectionHandler
	srv.commands["SET"] = stringHandler
	srv.commands["GET"] = stringHandler
}

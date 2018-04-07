package server

import (
	"fmt"
	"net"
	"redis_go/client"
	"redis_go/conf"
	"redis_go/database"
	"redis_go/handlers"
	"redis_go/log"
	"sync"
)

// Redis server
type Server struct {
	Config    *conf.ServerConfig
	Databases []*database.Database /* database*/
	password  string               /* Pass for AUTH command, or NULL */
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
	log.Info("redis server config: %+v", config)
	server.initDB()
	// init commands table
	server.populateCommandTable()
	log.Info("redis server: %+v", server)
	return server
}

func (srv *Server) serveClient(c *client.Client) {
	defer c.Release()
	// TODO lmj 增加Timeout的判断
	// TODO lmj 除了Timeout的方式，还有什么好的办法能够判断client端是否已经断开连接
	// loop to handle redis command
	for more := true; more; more = c.RequestReader.Buffered() != 0 {
		cmd, err := c.RequestReader.ReadCmd()
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
			handler.Process(c, cmd)
		} else {
			log.Errorf("command not found %s", cmd.GetOriginName())
			c.ResponseWriter.AppendError(fmt.Sprintf("command not found %s", cmd.GetOriginName()))
		}
		if err := c.ResponseWriter.Flush(); err != nil {
			log.Errorf("response writer flush data error %+v", err)
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
		go srv.serveClient(client.NewClient(cn, srv.getDefaultDB()))
		log.Info("new client come in ! from %+v", cn.RemoteAddr().String())
	}
}

func (srv *Server) initDB() {
	srv.Databases = make([]*database.Database, srv.Config.DBNum)
	// add default database
	for i := 0; i < srv.Config.DBNum; i++ {
		srv.Databases[i] = database.NewDatabase()
	}
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
	srv.commands["PING"] = connectionHandler
	srv.commands["TEST"] = connectionHandler
	srv.commands["SET"] = stringHandler
	srv.commands["GET"] = stringHandler
}

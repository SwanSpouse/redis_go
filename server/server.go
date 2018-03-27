package server

import (
	"net"
	"redis_go/client"
	"redis_go/conf"
	"redis_go/handlers"
	"redis_go/log"
	"sync"
)

// Redis server
type Server struct {
	config   *conf.ServerConfig
	commands map[string]handlers.BaseHandler
	mu       sync.RWMutex
}

func NewServer(config *conf.ServerConfig) *Server {
	if config == nil {
		config = new(conf.ServerConfig)
	}
	server := &Server{
		config:   config,
		commands: make(map[string]handlers.BaseHandler),
	}
	// init commands table
	server.populateCommandTable()
	return server
}

func (srv *Server) serveClient(c *client.Client) {
	defer c.Release()
	for !c.Closed {
		// TODO lmj 增加Timeout的判断

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
				handler.Process(c, cmd)
			} else {
				log.Info("command not found %s", cmd.GetName())
				c.ResponseWriter.AppendError("command not found")
			}
			if err := c.ResponseWriter.Flush(); err != nil {
				log.Info("response writer flush data error %+v", err)
				return
			}
		}
		log.Debug("No more data for current connection")
	}
	log.Debug("connection closed")
}

func (srv *Server) Serve(lis net.Listener) error {
	for {
		cn, err := lis.Accept()
		if err != nil {
			return err
		}
		log.Info("new client come in ! from %+v", cn.RemoteAddr().String())
		go srv.serveClient(client.NewClient(cn))
	}
}

// register all command handlers
func (srv *Server) populateCommandTable() {
	connectionHandler := new(handlers.ConnectionHandler)
	srv.commands["PING"] = connectionHandler
	srv.commands["TEST"] = connectionHandler
}
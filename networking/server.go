package networking

import (
	"sync"
	"net"
	"redis_go/conf"
	"github.com/labstack/gommon/log"
)

// Redis server
type Server struct {
	config   *conf.ServerConfig
	commands map[string]func(client *Client, command *Command) error
	mu       sync.RWMutex
}

func NewServer(config *conf.ServerConfig) *Server {
	if config == nil {
		config = new(conf.ServerConfig)
	}
	server := &Server{
		config:   config,
		commands: make(map[string]func(client *Client, command *Command) error),
	}
	// init commands table
	server.populateCommandTable()
	return server
}

func (srv *Server) serveClient(c *Client) {
	defer c.release()
	for !c.closed {
		// TODO lmj 增加Timeout的判断

		// loop to handle redis command
		for more := true; more; more = c.requestReader.Buffered() != 0 {
			cmd, err := c.requestReader.ReadCmd(nil)
			if err != nil {
				c.responseWriter.AppendErrorf("read command error %+v", err)
				continue
			}
			// TODO lmj handle redis commands
			/**
				首先判断是否在command table中,
					如果不在command table中,则返回command not found
					如果在command table中，则获取到相应的command handler来进行处理。
			 */
			log.Printf("get command from client %+v", cmd)
			if handler, ok := srv.commands[cmd.GetName()]; !ok {
				log.Printf("command not found %s", cmd.GetName())
				c.responseWriter.AppendError("command not found")
			} else {
				err := handler(c, cmd)
				if err != nil {
					log.Printf("hand command error %+v", err)
				}
			}
			if err := c.responseWriter.Flush(); err != nil {
				log.Printf("response writer flush data error %+v", err)
				return
			}
		}
		log.Printf("No more data for current connection")
	}
	log.Printf("connection closed")
}

func (srv *Server) Serve(lis net.Listener) error {
	for {
		cn, err := lis.Accept()
		if err != nil {
			return err
		}
		log.Printf("new client come in ! from %+v", cn.RemoteAddr().String())
		go srv.serveClient(newClient(cn))
	}
}

func (srv *Server) populateCommandTable() {
	srv.commands["ping"] = PING
}

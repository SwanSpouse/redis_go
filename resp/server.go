package resp

import (
	"sync"
	"net"
	"redis_go/conf"
	"github.com/labstack/gommon/log"
)

// Redis server
type Server struct {
	config   *conf.ServerConfig
	commands map[string]interface{}
	mu       sync.RWMutex
}

func NewServer(config *conf.ServerConfig) *Server {
	if config == nil {
		config = new(conf.ServerConfig)
	}
	return &Server{
		config:   config,
		commands: make(map[string]interface{}),
	}
}

func (srv *Server) serveClient(c *Client) {
	defer c.release()
	for !c.closed {
		for more := true; more; more = c.requestReader.Buffered() != 0 {
			cmd, err := c.requestReader.ReadCmd(nil)
			if err != nil {
				c.responseWriter.AppendErrorf("read command error %+v", err)
				continue
			}
			log.Printf("get command from client %+v", cmd)
			if _, ok := srv.commands[cmd.GetName()]; !ok {
				log.Printf("command not found %s", cmd.GetName())
				c.responseWriter.AppendError("command not found")
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

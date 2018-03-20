package server

import (
	"sync"
	"net"
	"redis_go/conf"
	"log"
	"strings"
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
			name, err := c.requestReader.PeekCmd()
			if err != nil {
				c.requestReader.SkipCmd()
			}
			norm := strings.ToLower(name)
			if _, ok := srv.commands[norm]; !ok {
				log.Printf("command not found %s", name)
				c.responseWriter.AppendError("command not found")
			}
			if err := c.responseWriter.Flush(); err != nil {
				log.Printf("reponse writer flush data error %+v", err)
				return
			}
		}
		log.Printf("No more data for current connection")
	}
	log.Printf("connection clonsed")
}

func (srv *Server) Serve(lis net.Listener) error {
	for {
		cn, err := lis.Accept()
		if err != nil {
			return err
		}
		log.Printf("new client come in !")
		go srv.serveClient(newClient(cn))
	}
}

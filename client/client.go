package client

import (
	"net"
	"redis_go/database"
	"redis_go/protocol"
	"sync"
	"sync/atomic"
)

var (
	clientInc  = uint64(0)
	readerPool sync.Pool
	writerPool sync.Pool
)

type Client struct {
	id             uint64
	cn             net.Conn
	db             *database.Database // chosen database
	Closed         bool
	RequestReader  *protocol.RequestReader
	ResponseWriter protocol.ResponseWriter
	args           uint64            // args number of command
	cmd            *protocol.Command // current command
	lastCmd        *protocol.Command // last command
}

func (c *Client) reset(cn net.Conn) {
	*c = Client{
		id: atomic.AddUint64(&clientInc, 1),
		cn: cn,
	}
	if v := readerPool.Get(); v != nil {
		rd := v.(*protocol.RequestReader)
		rd.Reset(cn)
		c.RequestReader = rd
	} else {
		c.RequestReader = protocol.NewRequestReader(cn)
	}
	if v := writerPool.Get(); v != nil {
		wr := v.(protocol.ResponseWriter)
		wr.Reset(cn)
		c.ResponseWriter = wr
	} else {
		c.ResponseWriter = protocol.NewResponseWriter(cn)
	}
}

func (c *Client) Release() {
	_ = c.cn.Close()
}

func NewClient(cn net.Conn, defaultDB *database.Database) *Client {
	c := new(Client)
	c.reset(cn)
	c.db = defaultDB
	return c
}

// unique client id
func (c *Client) ID() uint64 { return c.id }

// return the remote client address
func (c *Client) RemoteAddr() net.Addr {
	return c.cn.RemoteAddr()
}

func (c *Client) GetChosenDB() *database.Database {
	return c.db
}

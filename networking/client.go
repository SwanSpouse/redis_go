package networking

import (
	"net"
	"sync/atomic"
	"sync"
)

var (
	clientInc  = uint64(0)
	readerPool sync.Pool
	writerPool sync.Pool
)

type Client struct {
	id             uint64
	cn             net.Conn
	closed         bool
	requestReader  *RequestReader
	responseWriter ResponseWriter
	args           uint64   // args number of command
	cmd            *Command // current command
	lastCmd        *Command // last command
}

func (c *Client) reset(cn net.Conn) {
	*c = Client{
		id: atomic.AddUint64(&clientInc, 1),
		cn: cn,
	}
	if v := readerPool.Get(); v != nil {
		rd := v.(*RequestReader)
		rd.Reset(cn)
		c.requestReader = rd
	} else {
		c.requestReader = NewRequestReader(cn)
	}
	if v := writerPool.Get(); v != nil {
		wr := v.(ResponseWriter)
		wr.Reset(cn)
		c.responseWriter = wr
	} else {
		c.responseWriter = NewResponseWriter(cn)
	}
}

func (c *Client) Close() {
	c.closed = true
}

func (c *Client) release() {
	_ = c.cn.Close()
}

func newClient(cn net.Conn) *Client {
	c := new(Client)
	c.reset(cn)
	return c
}

// unique client id
func (c *Client) ID() uint64 { return c.id }

// return the remote client address
func (c *Client) RemoteAddr() net.Addr {
	return c.cn.RemoteAddr()
}

package client

import (
	"net"
	"redis_go/database"
	re "redis_go/error"
	"redis_go/log"
	"redis_go/protocol"
	"redis_go/tcp"
	"sync/atomic"
	"time"
)

var (
	clientInc = uint64(0)
)

type Client struct {
	id          uint64
	cn          net.Conn
	db          *database.Database // chosen database
	Closed      bool
	reader      *tcp.BufIoReader  // request reader
	writer      *tcp.BufIoWriter  // response writer
	args        uint64            // args number of command
	Cmd         *protocol.Command // current command
	lastCmd     *protocol.Command // last command
	timeoutTime time.Time         // timeout time
}

func (c *Client) reset(cn net.Conn) {
	*c = Client{
		id: atomic.AddUint64(&clientInc, 1),
		cn: cn,
	}
	c.reader = tcp.NewBufIoReader(cn)
	c.writer = tcp.NewBufIoWriter(cn)
}

func (c *Client) Release() {
	_ = c.cn.Close()
	c.reader.ReturnBufIoReader()
	c.writer.ReturnBufIoWriter()
	c.reader = nil
	c.writer = nil
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

func (c *Client) SelectedDatabase() *database.Database {
	return c.db
}

func (c *Client) Buffered() int {
	return c.reader.Buffered()
}

func (c *Client) Close() {
	c.Closed = true
}

func (c *Client) SetTimeout(duration time.Duration) {
	c.timeoutTime = time.Now().Add(duration)
}

func (c *Client) IsTimeout() bool {
	if c.timeoutTime.IsZero() {
		return false
	}
	return c.timeoutTime.Before(time.Now())
}

/**
construct a command from bufIoReader

command format:
	1. status reply     : +OK\r\n
	2. error reply      : -ERROR\r\n
	3. integer replay   : :1\r\n
	4. bulk reply       : $4\r\nPING\r\n
	5. multi bulk reply : *3\r\n$3\r\nSET\r\n$5\r\nMyKey\r\n$7\r\nMyValue\r\n
*/
func (c *Client) ReadCmd() (*protocol.Command, error) {
	// read one line from buffer
	line, err := c.reader.PeekLine(0)
	if err != nil || len(line) == 0 {
		return nil, err
	}
	cmd := protocol.NewCommand()
	switch line[0] {
	case '+', '-', ':':
		cmd.SetName(line.FirstWord())
	case '$':
		cmdName, err := c.reader.ReadBulkString()
		if err != nil || cmdName == "" {
			return nil, err
		}
		cmd.SetName(cmdName)
	case '*':
		arrayLen, err := c.reader.ReadArrayLen()
		if err != nil || arrayLen == 0 {
			return nil, err
		}
		for i := 0; i < arrayLen; i++ {
			arg, err := c.reader.ReadBulkString()
			if err != nil || arg == "" {
				return nil, err
			}
			if i == 0 {
				cmd.SetName(arg)
			} else {
				cmd.AddArgs(tcp.CommandArgument(arg))
			}
		}
	}
	// 更新client端记录的本次命令和上次命令
	c.lastCmd = c.Cmd
	c.Cmd = cmd
	log.Info("current command we received is %+v", cmd)
	return cmd, nil
}

func (c *Client) peekCmd(offset int) (string, error) {
	line, err := c.reader.PeekLine(offset)
	if err != nil {
		return "", err
	}
	offset += len(line)

	if len(line) == 0 {
		return "", nil
	} else if line[0] != '*' {
		return line.FirstWord(), nil
	}

	n, err := line.ParseSize('*', re.ErrInvalidMultiBulkLength)
	if err != nil {
		return "", err
	}

	if n < 1 {
		return c.peekCmd(offset)
	}

	line, err = c.reader.PeekLine(offset)
	if err != nil {
		return "", err
	}
	offset += len(line)

	n, err = line.ParseSize('$', re.ErrInvalidBulkLength)
	if err != nil {
		return "", err
	}

	data, err := c.reader.PeekN(offset, int(n))
	return string(data), err
}

func (c *Client) Response(value interface{}) {
	c.writer.Append(value)
}

func (c *Client) ResponseOK() {
	c.writer.AppendOK()
}

func (c *Client) ResponseError(msg string, args ...interface{}) {
	if len(args) == 0 {
		c.writer.AppendError(msg)
	} else {
		c.writer.AppendErrorf(msg, args...)
	}
	c.Flush()
}

func (c *Client) ResponseReError(err error, args ...interface{}) {
	if re.IsProtocolError(err) {
		switch err {
		case re.ErrNilValue:
			c.Response(nil)
		default:
			c.ResponseError(err.Error(), args...)
		}
	} else {
		c.ResponseError(err.Error(), args...)
	}
}

func (c *Client) Flush() error {
	return c.writer.Flush()
}

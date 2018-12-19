package client

import (
	"errors"
	"net"
	"redis_go/database"
	re "redis_go/error"
	"redis_go/loggers"
	"redis_go/tcp"
	"sync"
	"time"
)

var clientPool = &sync.Pool{
	New: func() interface{} {
		return new(Client)
	},
}

type Client struct {
	id          int64              // Client ID
	cn          net.Conn           // TCP connection
	db          *database.Database // chosen database
	Closed      bool               // isClientClosed
	reader      *tcp.BufIoReader   // request reader
	writer      *tcp.BufIoWriter   // response writer
	Argv        []string           // arguments vector
	Argc        int                // arguments counter
	Cmd         *Command           // current command
	LastCmd     *Command           // last command
	Dirty       int64
	execTimeout time.Time
	idleTimeout time.Time // timeout
}

func (c *Client) reset(clientId int64, cn net.Conn, defaultDB *database.Database) {
	c.id = clientId
	c.cn = cn
	c.db = defaultDB
	c.Closed = false
	c.reader = tcp.NewBufIoReader(cn)
	c.writer = tcp.NewBufIoWriter(cn)
	c.Argv = nil
	c.Argc = 0
	c.Cmd = nil
	c.LastCmd = nil
	c.Dirty = 0
	c.execTimeout = time.Time{}
	c.idleTimeout = time.Time{}
}

func (c *Client) release() {
	if c.IsFakeClient() {
		return
	}
	tcp.ReturnBufIoReader(c.reader)
	tcp.ReturnBufIoWriter(c.writer)
	c.cn.Close()
}

func NewClient(clientId int64, cn net.Conn, defaultDB *database.Database) *Client {
	c := clientPool.Get().(*Client)
	c.reset(clientId, cn, defaultDB)
	return c
}

func ReturnClient(c *Client) {
	clientPool.Put(c)
}

func NewFakeClient() *Client {
	return &Client{
		id: -1,
	}
}

func (c *Client) IsFakeClient() bool {
	return c.id == -1
}

// unique client id
func (c *Client) ID() int64 { return c.id }

// return the remote client address
func (c *Client) RemoteAddr() net.Addr {
	if c.IsFakeClient() {
		return nil
	}
	return c.cn.RemoteAddr()
}

func (c *Client) SetDatabase(db *database.Database) {
	c.db = db
}

func (c *Client) SelectedDatabase() *database.Database {
	return c.db
}

func (c *Client) Buffered() int {
	if c.IsFakeClient() {
		return 0
	}
	return c.reader.Buffered()
}

func (c *Client) Close() {
	c.Closed = true
	c.release()
}

func (c *Client) SetIdleTimeout(duration time.Duration) {
	c.idleTimeout = time.Now().Add(duration)
}

func (c *Client) IsIdleTimeout() bool {
	if c.idleTimeout.IsZero() {
		return false
	}
	return c.idleTimeout.Before(time.Now())
}

func (c *Client) GetIdleTimeoutAt() time.Time {
	return c.idleTimeout
}

func (c *Client) SetExecTimeout(duration time.Duration) {
	c.execTimeout = time.Now().Add(duration)
}

func (c *Client) IsExecTimeout() bool {
	if c.execTimeout.IsZero() {
		return false
	}
	return c.execTimeout.Before(time.Now())
}

func (c *Client) GetExecTimeoutAt() time.Time {
	return c.execTimeout
}

/**
command format:
	1. status reply     : +OK\r\n
	2. error reply      : -ERROR\r\n
	3. integer replay   : :1\r\n
	4. bulk reply       : $4\r\nPING\r\n
	5. multi bulk reply : *3\r\n$3\r\nSET\r\n$5\r\nMyKey\r\n$7\r\nMyValue\r\n
*/
func (c *Client) ProcessInputBuffer() error {
	// read one line from buffer
	line, err := c.reader.PeekLine(0)
	if err != nil || len(line) == 0 {
		return err
	}
	c.Argv = make([]string, 0)
	c.Argc = 0
	switch line[0] {
	case '+', '-', ':':
		c.Argv = append(c.Argv, line.FirstWord())
		c.Argc += 1
	case '$':
		cmdName, err := c.reader.ReadBulkString()
		if err != nil || cmdName == "" {
			return err
		}
		c.Argv = append(c.Argv, cmdName)
	case '*':
		arrayLen, err := c.reader.ReadArrayLen()
		if err != nil || arrayLen == 0 {
			return err
		}
		for i := 0; i < arrayLen; i++ {
			arg, err := c.reader.ReadBulkString()
			if err != nil || arg == "" {
				return err
			}
			c.Argv = append(c.Argv, arg)
		}
	}
	c.Argc = len(c.Argv)
	loggers.Info("server receive:%d args, argv:%+v", c.Argc, c.Argv)
	return nil
}

func (c *Client) peekCmd(offset int) (string, error) {
	if c.IsFakeClient() {
		return "", errors.New("this client is a fake client")
	}
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
	if c.IsFakeClient() {
		return
	}
	loggers.Debug("server response:%+v", value)
	c.writer.Append(value)
}

func (c *Client) ResponseOK() {
	if c.IsFakeClient() {
		return
	}
	c.writer.AppendOK()
}

func (c *Client) ResponseError(msg string, args ...interface{}) {
	if c.IsFakeClient() {
		return
	}
	if len(args) == 0 {
		c.writer.AppendError(msg)
	} else {
		c.writer.AppendErrorf(msg, args...)
	}
	c.Flush()
}

func (c *Client) ResponseReError(err error, args ...interface{}) {
	if c.IsFakeClient() {
		loggers.Errorf("fake client receive a error:%+v", err)
		return
	}
	loggers.Debug("server response error:%+v args:%+v", err, args)
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
	if c.IsFakeClient() {
		return nil
	}
	return c.writer.Flush()
}

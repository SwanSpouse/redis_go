package protocol

import (
	"io"
	re "redis_go/error"
	"redis_go/log"
	"redis_go/tcp"
)

type RequestReader struct {
	reader *tcp.BufIoReader
}

// NewRequestReader wraps any reader interface
func NewRequestReader(rd io.Reader) *RequestReader {
	r := new(tcp.BufIoReader)
	r.Reset(rd)
	return &RequestReader{reader: r}
}

func (r *RequestReader) Buffered() int {
	return r.reader.Buffered()
}

func (r *RequestReader) Reset(rd io.Reader) {
	r.reader.Reset(rd)
}

// peek next command name
func (r *RequestReader) PeekCmdName() (string, error) {
	return r.peekCmd(0)
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
func (r *RequestReader) ReadCmd() (*Command, error) {
	// read one line from buffer
	line, err := r.reader.PeekLine(0)
	if err != nil || len(line) == 0 {
		return nil, err
	}
	cmd := NewCommand()
	switch line[0] {
	case '+', '-', ':':
		cmd.SetName(line.FirstWord())
	case '$':
		cmdName, err := r.reader.ReadBulkString()
		if err != nil || cmdName == "" {
			return nil, err
		}
		cmd.SetName(cmdName)
	case '*':
		arrayLen, err := r.reader.ReadArrayLen()
		if err != nil || arrayLen == 0 {
			return nil, err
		}
		for i := 0; i < arrayLen; i++ {
			arg, err := r.reader.ReadBulkString()
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
	log.Info("current command we received is %+v", cmd)
	return cmd, nil
}

func (r *RequestReader) peekCmd(offset int) (string, error) {
	line, err := r.reader.PeekLine(offset)
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
		return r.peekCmd(offset)
	}

	line, err = r.reader.PeekLine(offset)
	if err != nil {
		return "", err
	}
	offset += len(line)

	n, err = line.ParseSize('$', re.ErrInvalidBulkLength)
	if err != nil {
		return "", err
	}

	data, err := r.reader.PeekN(offset, int(n))
	return string(data), err
}

// skips the next command
func (r *RequestReader) SkipCmd() error {
	c, err := r.reader.PeekByte()
	if err != nil {
		return err
	}
	if c != '*' {
		_, err = r.reader.ReadLine()
		return err
	}
	n, err := r.reader.ReadArrayLen()
	if err != nil {
		return err
	}
	if n < 1 {
		return r.SkipCmd()
	}
	for i := 0; i < n; i++ {
		if err := r.reader.SkipBulk(); err != nil {
			return err
		}
	}
	return nil
}

// --------------------------------------------------------------------

type RequestWriter struct {
	w *tcp.BufIoWriter
}

// NewRequestWriter wraps any writer interface
func NewRequestWriter(wr io.Writer) *RequestWriter {
	w := new(tcp.BufIoWriter)
	w.Reset(wr)
	return &RequestWriter{w: w}
}

func (w *RequestWriter) Reset(wr io.Writer) {
	w.w.Reset(wr)
}

func (w *RequestWriter) Buffered() int {
	return w.w.Buffered()
}

func (w *RequestWriter) Flush() error {
	return w.w.Flush()
}

func (w *RequestWriter) WriteCmd(cmd string, args ...[]byte) {
	w.w.AppendArrayLen(len(args) + 1)
	w.w.AppendBulkString(cmd)
	for _, arg := range args {
		w.w.AppendBulk(arg)
	}
}

func (w *RequestWriter) WriteCmdString(cmd string, args ...string) {
	w.w.AppendArrayLen(len(args) + 1)
	w.w.AppendBulkString(cmd)
	for _, arg := range args {
		w.w.AppendBulkString(arg)
	}
}

func (w *RequestWriter) WriteMultiBulkSize(n int) error {
	if n < 0 {
		return re.ErrInvalidMultiBulkLength
	}
	w.w.AppendArrayLen(n)
	return nil
}

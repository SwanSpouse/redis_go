package resp

import (
	"io"
	"github.com/labstack/gommon/log"
)

type RequestReader struct {
	reader *bufIoReader
}

// NewRequestReader wraps any reader interface
func NewRequestReader(rd io.Reader) *RequestReader {
	r := new(bufIoReader)
	r.reset(mkStdBuffer(), rd)
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

// construct a command from bufIoReader
func (r *RequestReader) ReadCmd(cmd *Command) (*Command, error) {
	if cmd == nil {
		cmd = NewCommand()
	}
	len, err := r.reader.ReadArrayLen()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len; i ++ {
		arg, err := r.reader.ReadBulkString()
		if err != nil {
			return nil, err
		}
		// first string is command name
		if i == 0 {
			cmd.SetName(arg)
		} else {
			cmd.AddArgs(CommandArgument([]byte(arg)))
		}
	}
	log.Printf("current cmd %+v", cmd)
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

	n, err := line.ParseSize('*', errInvalidMultiBulkLength)
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

	n, err = line.ParseSize('$', errInvalidBulkLength)
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
	for i := 0; i < n; i ++ {
		if err := r.reader.SkipBulk(); err != nil {
			return err
		}
	}
	return nil
}

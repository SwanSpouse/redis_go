package resp

import "io"

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

func (r *RequestReader) PeekCmd() (string, error) {
	return r.peekCmd(0)
}

func (r *RequestReader) ReadCmd(cmd *Command) (*Command, error) {
	if cmd == nil {
		cmd = new(Command)
	} else {
	}
	//return cmd, readCommand(cmd, r.reader)
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

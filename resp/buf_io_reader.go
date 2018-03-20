package resp

import (
	"bytes"
	"io"
)

type bufIoReader struct {
	rd  io.Reader
	buf []byte
	r   int // reader index
	w   int // writer index
}

// reset buffer & rd
func (b *bufIoReader) reset(buf []byte, rd io.Reader) {
	*b = bufIoReader{buf: buf, rd: rd}
}

// compact moves the unread chunk to the beginning of the buffer
func (b *bufIoReader) compact() {
	if b.r > 0 {
		copy(b.buf, b.buf[b.r:b.w])
		b.w = b.w - b.r
		b.r = 0
	}
}

// returns the number of buffered bytes unread
func (b *bufIoReader) Buffered() int {
	return b.w - b.r
}

// make sure that sz bytes can be buffered
func (b *bufIoReader) require(sz int) error {
	extra := sz - b.Buffered()
	if extra < 1 {
		return nil
	}
	// compact first
	b.compact()

	// grow the buffer if necessary
	if n := b.w + extra; n > len(b.buf) {
		buf := make([]byte, n)
		copy(buf, buf[:b.w])
		b.buf = buf
	}

	// read data into buffer
	n, err := io.ReadAtLeast(b.rd, b.buf[b.w:], extra)
	b.w += n
	return err
}

// tries to read more data into the buffer
func (b *bufIoReader) fill() error {
	b.compact()

	if b.w < len(b.buf) {
		n, err := b.rd.Read(b.buf[b.w:])
		b.w += n
		return err
	}
	return nil
}

// peek byte of the buffer
func (b *bufIoReader) PeekByte() (byte, error) {
	if err := b.require(1); err != nil {
		return 0, err
	}
	return b.buf[b.r], nil
}

// PeekLine returns the next line until CRLF without reading it
func (b *bufIoReader) PeekLine(offset int) (buffer, error) {
	index := -1

	// try to find the end of the line
	start := b.r + offset
	if start < b.w {
		index = bytes.IndexByte(b.buf[start:b.w], '\n')
	}

	// try to read more data into the buffer if not in the buffer
	if index < 0 {
		if err := b.fill(); err != nil {
			return nil, err
		}
		start = b.r + offset
		if start < b.w {
			index = bytes.IndexByte(b.buf[start:b.w], '\n')
		}
	}

	// fail if still nothing found
	if index < 0 {
		return nil, errInlineRequestTooLong
	}
	return buffer(b.buf[start : start+index+1]), nil
}

/*
	状态回复(status reply)  的第一个字节是        +
	错误回复(error reply)   的第一个字节是        -
	整数回复(integer reply) 的第一个字节是        :
	批量回复(bulk reply)    的第一个字节是        $
	多条批量回复(multi bulk reply)的第一个字节是   *
*/
func (b *bufIoReader) PeekType() (t ResponseType, err error) {
	if err = b.require(1); err != nil {
		return
	}
	switch b.buf[b.r] {
	case '*':
		t = TypeArray
	case '$':
		if err = b.require(2); err != nil {
			return
		}
		if b.buf[b.r+1] == '-' {
			t = TypeNil
		} else {
			t = TypeBulk
		}
	case '+':
		t = TypeInline
	case '-':
		t = TypeError
	case ':':
		t = TypeInt
	}
	return
}

func (b *bufIoReader) PeekN(offset, n int) ([]byte, error) {
	if err := b.require(offset + n); err != nil {
		return nil, err
	}
	return b.buf[b.r+offset : b.r+offset+n], nil
}

// return the next line until CRLF
func (b *bufIoReader) ReadLine() (buffer, error) {
	line, err := b.PeekLine(0)
	b.r += len(line)
	return line, err
}

func (b *bufIoReader) ReadNil() error {
	line, err := b.ReadLine()
	if err != nil {
		return err
	}
	if len(line) < 3 || !bytes.Equal(line[:3], binNIL[:3]) {
		return errNotANilMessage
	}
	return nil
}

func (b *bufIoReader) ReadInt() (int64, error) {
	line, err := b.ReadLine()
	if err != nil {
		return 0, err
	}
	return line.ParseInt()
}

func (b *bufIoReader) ReadError() (string, error) {
	line, err := b.ReadLine()
	if err != nil {
		return "", err
	}
	return line.ParseMessage('-')
}

func (b *bufIoReader) ReadInlineString() (string, error) {
	line, err := b.ReadLine()
	if err != nil {
		return "", err
	}
	return line.ParseMessage('+')
}

func (b *bufIoReader) ReadArrayLen() (int, error) {
	line, err := b.ReadLine()
	if err != nil {
		return 0, err
	}
	sz, err := line.ParseSize('*', errInvalidMultiBulkLength)
	if err != nil {
		return 0, err
	}
	return int(sz), nil
}

func (b *bufIoReader) ReadBulkLen() (int64, error) {
	line, err := b.ReadLine()
	if err != nil {
		return 0, err
	}
	return line.ParseSize('$', errInvalidBulkLength)
}

func (b *bufIoReader) ReadBulk(p []byte) ([]byte, error) {
	sz, err := b.ReadBulkLen()
	if err != nil {
		return p, err
	}
	//TODO lmj sz + 2 means size + \r\n ?
	if err := b.require(int(sz + 2)); err != nil {
		return p, err
	}
	p = append(p, b.buf[b.r:b.r+int(sz)]...)
	b.r += int(sz + 2)
	return p, nil
}

func (b *bufIoReader) ReadBulkString() (string, error) {
	sz, err := b.ReadBulkLen()
	if err != nil {
		return "", err
	}
	if err := b.require(int(sz + 2)); err != nil {
		return "", err
	}
	s := string(b.buf[b.r : b.r+int(sz)])
	b.r += int(sz + 2)
	return s, nil
}

func (b *bufIoReader) Scan(vv ...interface{}) error {
	//TODO lmj
	return nil
}

func (b *bufIoReader) Reset(r io.Reader) {
	b.reset(b.buf, r)
}

func (b *bufIoReader) skip(sz int) {
	if b.Buffered() >= sz {
		b.r += sz
	}
	//TODO lmj need first compact ?
}

func (b *bufIoReader) SkipBulk() error {
	sz, err := b.ReadBulkLen()
	if err != nil {
		return err
	}
	return b.skipN(sz + 2)
}

func (b *bufIoReader) skipN(sz int64) error {
	// if bulk doesn't overflow buffer
	extra := sz - int64(b.Buffered())
	if extra < 1 {
		b.r += int(sz)
		return nil
	}
	// otherwise, reset buffer
	b.r = 0
	b.w = 0

	// ... and discard the extra bytes
	x := extra
	r := io.LimitReader(b.rd, x)
	for {
		n, err := r.Read(b.buf)
		x -= int64(n)

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}

	if x != 0 {
		return io.EOF
	}
	return nil
}

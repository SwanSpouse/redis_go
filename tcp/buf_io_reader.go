package tcp

import (
	"bytes"
	"io"
	"net"
	re "redis_go/error"
	"redis_go/loggers"
	"sync"
)

type BufIoReader struct {
	rd  io.Reader
	buf []byte
	r   int // reader index
	w   int // writer index
}

var (
	ReaderPool sync.Pool // Reader连接池
)

func NewBufIoReader(cn net.Conn) *BufIoReader {
	var r *BufIoReader
	if v := ReaderPool.Get(); v != nil {
		loggers.Debug("Get BufIoReader from ReaderPool")
		r = v.(*BufIoReader)
	} else {
		loggers.Debug("Can not get BufIoReader from ReaderPool, return a New BufIOReader")
		r = new(BufIoReader)
	}
	r.Reset(cn)
	return r
}

func (r *BufIoReader) ReturnBufIoReader() {
	ReaderPool.Put(r)
}

func NewBufIoReaderWithoutConn() *BufIoReader {
	return new(BufIoReader)
}

// reset buffer & rd
func (r *BufIoReader) reset(buf []byte, rd io.Reader) {
	*r = BufIoReader{buf: buf, rd: rd}
}

// compact moves the unread chunk to the beginning of the buffer
func (r *BufIoReader) compact() {
	if r.r > 0 {
		copy(r.buf, r.buf[r.r:r.w])
		r.w = r.w - r.r
		r.r = 0
	}
}

// returns the number of buffered bytes unread
func (r *BufIoReader) Buffered() int {
	//log.Debug("[BUFFERED DATA]:%s", string(r.buf[r.r:r.w]))
	return r.w - r.r
}

// make sure that sz bytes can be buffered
func (r *BufIoReader) require(sz int) error {
	extra := sz - r.Buffered()
	if extra < 1 {
		return nil
	}
	// compact first
	r.compact()

	// grow the buffer if necessary
	if n := r.w + extra; n > len(r.buf) {
		buf := make([]byte, n)
		copy(buf, buf[:r.w])
		r.buf = buf
	}

	// read data into buffer
	n, err := io.ReadAtLeast(r.rd, r.buf[r.w:], extra)
	r.w += n

	return err
}

// tries to read more data into the buffer
func (r *BufIoReader) fill() error {
	r.compact()

	if r.w < len(r.buf) {
		n, err := r.rd.Read(r.buf[r.w:])
		r.w += n
		//log.Info("current io reader buffer %s", string(r.buf[r.r:r.w]))
		return err
	}
	return nil
}

// peek byte of the buffer
func (r *BufIoReader) PeekByte() (byte, error) {
	if err := r.require(1); err != nil {
		return 0, err
	}
	return r.buf[r.r], nil
}

// PeekLine returns the next line until CRLF without reading it
func (r *BufIoReader) PeekLine(offset int) (buffer, error) {
	index := -1

	// try to find the end of the line
	start := r.r + offset
	if start < r.w {
		index = bytes.IndexByte(r.buf[start:r.w], '\n')
	}

	// try to read more data into the buffer if not in the buffer
	if index < 0 {
		if err := r.fill(); err != nil {
			return nil, err
		}
		start = r.r + offset
		if start < r.w {
			index = bytes.IndexByte(r.buf[start:r.w], '\n')
		}
	}

	// fail if still nothing found
	if index < 0 {
		return nil, re.ErrInlineRequestTooLong
	}
	return buffer(r.buf[start : start+index+1]), nil
}

/*
	状态回复(status reply)  的第一个字节是        +
	错误回复(error reply)   的第一个字节是        -
	整数回复(integer reply) 的第一个字节是        :
	批量回复(bulk reply)    的第一个字节是        $
	多条批量回复(multi bulk reply)的第一个字节是   *
*/
func (r *BufIoReader) PeekType() (t ResponseType, err error) {
	loggers.Info("peek type start")
	if err = r.require(1); err != nil {
		loggers.Info("peek type err %+v", err)
		return
	}
	switch r.buf[r.r] {
	case '*':
		t = TypeArray
	case '$':
		if err = r.require(2); err != nil {
			return
		}
		if r.buf[r.r+1] == '-' {
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

func (r *BufIoReader) PeekN(offset, n int) ([]byte, error) {
	if err := r.require(offset + n); err != nil {
		return nil, err
	}
	return r.buf[r.r+offset : r.r+offset+n], nil
}

// return the next line until CRLF
func (r *BufIoReader) ReadLine() (buffer, error) {
	line, err := r.PeekLine(0)
	r.r += len(line)
	return line, err
}

func (r *BufIoReader) ReadNil() error {
	line, err := r.ReadLine()
	if err != nil {
		return err
	}
	if len(line) < 3 || !bytes.Equal(line[:3], BinNIL[:3]) {
		return re.ErrNotANilMessage
	}
	return nil
}

func (r *BufIoReader) ReadInt() (int64, error) {
	line, err := r.ReadLine()
	if err != nil {
		return 0, err
	}
	return line.ParseInt()
}

func (r *BufIoReader) ReadError() (string, error) {
	line, err := r.ReadLine()
	if err != nil {
		return "", err
	}
	return line.ParseMessage('-')
}

func (r *BufIoReader) ReadInlineString() (string, error) {
	line, err := r.ReadLine()
	if err != nil {
		return "", err
	}
	return line.ParseMessage('+')
}

func (r *BufIoReader) ReadArrayLen() (int, error) {
	line, err := r.ReadLine()
	if err != nil {
		return 0, err
	}
	sz, err := line.ParseSize('*', re.ErrInvalidMultiBulkLength)
	if err != nil {
		return 0, err
	}
	return int(sz), nil
}

func (r *BufIoReader) ReadBulkLen() (int64, error) {
	line, err := r.ReadLine()
	if err != nil {
		return 0, err
	}
	return line.ParseSize('$', re.ErrInvalidBulkLength)
}

func (r *BufIoReader) ReadBulk(p []byte) ([]byte, error) {
	sz, err := r.ReadBulkLen()
	if err != nil {
		return p, err
	}
	if err := r.require(int(sz + 2)); err != nil {
		return p, err
	}
	p = append(p, r.buf[r.r:r.r+int(sz)]...)
	r.r += int(sz + 2)
	return p, nil
}

func (r *BufIoReader) ReadBulkString() (string, error) {
	sz, err := r.ReadBulkLen()
	if err != nil {
		return "", err
	}
	if err := r.require(int(sz + 2)); err != nil {
		return "", err
	}
	s := string(r.buf[r.r : r.r+int(sz)])
	r.r += int(sz + 2)
	return s, nil
}

func (r *BufIoReader) Scan(vv ...interface{}) error {
	//TODO lmj
	return nil
}

func (r *BufIoReader) Reset(ioReader io.Reader) {
	r.reset(MkStdBuffer(), ioReader)
}

func (r *BufIoReader) skip(sz int) {
	if r.Buffered() >= sz {
		r.r += sz
	}
	//TODO lmj need first compact ?
}

func (r *BufIoReader) SkipBulk() error {
	sz, err := r.ReadBulkLen()
	if err != nil {
		return err
	}
	return r.skipN(sz + 2)
}

func (r *BufIoReader) skipN(sz int64) error {
	// if bulk doesn't overflow buffer
	extra := sz - int64(r.Buffered())
	if extra < 1 {
		r.r += int(sz)
		return nil
	}
	// otherwise, reset buffer
	r.r = 0
	r.w = 0

	// ... and discard the extra bytes
	x := extra
	reader := io.LimitReader(r.rd, x)
	for {
		n, err := reader.Read(r.buf)
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

package tcp

import (
	"fmt"
	"io"
	"net"
	"redis_go/log"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// response type iota
const (
	TypeUnknown ResponseType = iota
	TypeArray
	TypeBulk
	TypeInline
	TypeError
	TypeInt
	TypeNil
)

type ResponseType uint8

func (t ResponseType) String() string {
	switch t {
	case TypeArray:
		return "Array"
	case TypeBulk:
		return "Bulk"
	case TypeInline:
		return "Inline"
	case TypeError:
		return "Error"
	case TypeInt:
		return "Int"
	case TypeNil:
		return "Nil"
	}
	return "Unknown"
}

var (
	BinCRLF = []byte("\r\n")
	BinOK   = []byte("+OK\r\n")
	BinZERO = []byte(":0\r\n")
	BinONE  = []byte(":1\r\n")
	BinNIL  = []byte("$-1\r\n")
)

// MaxBufferSize is the max request/response buffer size
const MaxBufferSize = 64 * 1024

// len of \r\n
const CRLFLen = 2

func MkStdBuffer() []byte { return make([]byte, MaxBufferSize) }

// CommandArgument is an argument of a command
type CommandArgument []byte

//------------------------------------------------------------------------------------------------

var (
	WriterPool sync.Pool // Writer连接池
)

type BufIoWriter struct {
	io.Writer
	buf []byte
	mu  sync.Mutex
}

func NewBufIoWriter(cn net.Conn) *BufIoWriter {
	var w *BufIoWriter
	if v := WriterPool.Get(); v != nil {
		log.Debug("Get BufIoWriter from WriterPool")
		w = v.(*BufIoWriter)
	} else {
		log.Debug("Can not get BufIoWriter from WriterPool, return a New BufIoWriter")
		w = new(BufIoWriter)
	}
	w.Reset(cn)
	return w
}

func (b *BufIoWriter) ReturnBufIoWriter() {
	WriterPool.Put(b)
}

func NewBufIoWriterWithoutConn() *BufIoWriter {
	return new(BufIoWriter)
}

// returns the number of buffered bytes
func (b *BufIoWriter) Buffered() int {
	b.mu.Lock()
	n := len(b.buf)
	b.mu.Unlock()
	return n
}

func (b *BufIoWriter) appendSize(c byte, n int64) {
	b.buf = append(b.buf, c)
	b.buf = append(b.buf, strconv.FormatInt(n, 10)...)
	b.buf = append(b.buf, BinCRLF...)
}

func (b *BufIoWriter) reset(buf []byte, wr io.Writer) {
	*b = BufIoWriter{buf: buf[:0], Writer: wr}
}

func (b *BufIoWriter) flush() error {
	if len(b.buf) == 0 {
		return nil
	}

	if _, err := b.Write(b.buf); err != nil {
		return err
	}
	b.buf = b.buf[:0]
	return nil
}

// appends an array header  to the output buffer
func (b *BufIoWriter) AppendArrayLen(n int) {
	b.mu.Lock()
	b.appendSize('*', int64(n))
	b.mu.Unlock()
}

// appends bulk bytes to the output buffer
func (b *BufIoWriter) AppendBulk(p []byte) {
	b.mu.Lock()
	b.appendSize('$', int64(len(p)))
	b.buf = append(b.buf, p...)
	b.buf = append(b.buf, BinCRLF...)
	b.mu.Unlock()
}

// appends a bulk string to the output buffer
func (b *BufIoWriter) AppendBulkString(s string) {
	b.mu.Lock()
	b.appendSize('$', int64(len(s)))
	b.buf = append(b.buf, s...)
	b.buf = append(b.buf, BinCRLF...)
	b.mu.Unlock()
}

// appends inline bytes to the output buffer
func (b *BufIoWriter) AppendInline(p []byte) {
	b.mu.Lock()
	b.buf = append(b.buf, '+')
	b.buf = append(b.buf, p...)
	b.buf = append(b.buf, BinCRLF...)
	b.mu.Unlock()
}

// appends an inline string to the output buffer
func (b *BufIoWriter) AppendInlineString(s string) {
	b.mu.Lock()
	b.buf = append(b.buf, '+')
	b.buf = append(b.buf, s...)
	b.buf = append(b.buf, BinCRLF...)
	b.mu.Unlock()
}

// appends an error message to the output buffer
func (b *BufIoWriter) AppendError(msg string) {
	b.mu.Lock()
	b.buf = append(b.buf, '-')
	b.buf = append(b.buf, msg...)
	b.buf = append(b.buf, BinCRLF...)
	b.mu.Unlock()
}

func (b *BufIoWriter) AppendErrorf(pattern string, args ...interface{}) {
	b.AppendError(fmt.Sprintf(pattern, args...))
}

// appends a numeric response to the output buffer
func (b *BufIoWriter) AppendInt(n int64) {
	b.mu.Lock()
	switch n {
	case 0:
		b.buf = append(b.buf, BinZERO...)
	case 1:
		b.buf = append(b.buf, BinONE...)
	default:
		b.buf = append(b.buf, ':')
		b.buf = append(b.buf, strconv.FormatInt(n, 10)...)
		b.buf = append(b.buf, BinCRLF...)
	}
	b.mu.Unlock()
}

// appends a nil-value to the output
func (b *BufIoWriter) AppendNil() {
	b.mu.Lock()
	b.buf = append(b.buf, BinNIL...)
	b.mu.Unlock()
}

// appends OK to the output buffer
func (b *BufIoWriter) AppendOK() {
	b.mu.Lock()
	b.buf = append(b.buf, BinOK...)
	b.mu.Unlock()
}

// flush pending buffer
func (b *BufIoWriter) Flush() error {
	b.mu.Lock()
	err := b.flush()
	defer b.mu.Unlock()
	return err
}

// resets the writer with an new interface
func (b *BufIoWriter) Reset(w io.Writer) {
	b.reset(MkStdBuffer(), w)
}

// Append implements ResponseWriter
func (w *BufIoWriter) Append(v interface{}) error {
	switch v := v.(type) {
	case nil:
		w.AppendNil()
	case error:
		msg := v.Error()
		if !strings.HasPrefix(msg, "ERR ") {
			msg = "ERR " + msg
		}
		w.AppendError(msg)
	case bool:
		if v {
			w.AppendInt(1)
		} else {
			w.AppendInt(0)
		}
	case int:
		w.AppendInt(int64(v))
	case int8:
		w.AppendInt(int64(v))
	case int16:
		w.AppendInt(int64(v))
	case int32:
		w.AppendInt(int64(v))
	case int64:
		w.AppendInt(v)
	case uint:
		w.AppendInt(int64(v))
	case uint8:
		w.AppendInt(int64(v))
	case uint16:
		w.AppendInt(int64(v))
	case uint32:
		w.AppendInt(int64(v))
	case uint64:
		w.AppendInt(int64(v))
	case string:
		w.AppendBulkString(v)
	case []byte:
		w.AppendBulk(v)
	case CommandArgument:
		w.AppendBulk(v)
	case float32:
		w.AppendInlineString(strconv.FormatFloat(float64(v), 'f', -1, 32))
	case float64:
		w.AppendInlineString(strconv.FormatFloat(v, 'f', -1, 64))
	default:
		switch reflect.TypeOf(v).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(v)

			w.AppendArrayLen(s.Len())
			for i := 0; i < s.Len(); i++ {
				w.Append(s.Index(i).Interface())
			}
		case reflect.Map:
			s := reflect.ValueOf(v)

			w.AppendArrayLen(s.Len() * 2)
			for _, key := range s.MapKeys() {
				w.Append(key.Interface())
				w.Append(s.MapIndex(key).Interface())
			}
		default:
			return fmt.Errorf("resp: unsupported type %T", v)
		}
	}
	return nil
}

// 来啥写啥，一点儿不变
func (b *BufIoWriter) AppendRawString(rawInput []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()
	log.Debug("raw string %s", string(rawInput))
	b.buf = append(b.buf, rawInput...)
}

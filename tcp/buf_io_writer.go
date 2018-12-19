package tcp

import (
	"fmt"
	"io"
	"net"
	re "redis_go/error"
	"redis_go/loggers"
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

// Writer连接池
var WriterPool = &sync.Pool{
	New: func() interface{} {
		return new(BufIoWriter)
	},
}

type BufIoWriter struct {
	io.Writer
	buf []byte
	mu  sync.Mutex
}

func NewBufIoWriter(cn net.Conn) *BufIoWriter {
	w := WriterPool.Get().(*BufIoWriter)
	w.Reset(cn)
	return w
}

func ReturnBufIoWriter(w *BufIoWriter) {
	WriterPool.Put(w)
}

func InitBufIoWriterPool(size int) {
	for i := 0; i < size; i ++ {
		WriterPool.Put(new(BufIoWriter))
	}
}

// returns the number of buffered bytes
func (w *BufIoWriter) Buffered() int {
	w.mu.Lock()
	n := len(w.buf)
	w.mu.Unlock()
	return n
}

func (w *BufIoWriter) appendSize(c byte, n int64) {
	w.buf = append(w.buf, c)
	w.buf = append(w.buf, strconv.FormatInt(n, 10)...)
	w.buf = append(w.buf, BinCRLF...)
}

func (w *BufIoWriter) reset(buf []byte, wr io.Writer) {
	*w = BufIoWriter{buf: buf[:0], Writer: wr}
}

func (w *BufIoWriter) flush() error {
	if len(w.buf) == 0 {
		return nil
	}

	if _, err := w.Write(w.buf); err != nil {
		return err
	}
	w.buf = w.buf[:0]
	return nil
}

// appends an array header  to the output buffer
func (w *BufIoWriter) AppendArrayLen(n int) {
	w.mu.Lock()
	w.appendSize('*', int64(n))
	w.mu.Unlock()
}

// appends bulk bytes to the output buffer
func (w *BufIoWriter) AppendBulk(p []byte) {
	w.mu.Lock()
	w.appendSize('$', int64(len(p)))
	w.buf = append(w.buf, p...)
	w.buf = append(w.buf, BinCRLF...)
	w.mu.Unlock()
}

// appends a bulk string to the output buffer
func (w *BufIoWriter) AppendBulkString(s string) {
	w.mu.Lock()
	w.appendSize('$', int64(len(s)))
	w.buf = append(w.buf, s...)
	w.buf = append(w.buf, BinCRLF...)
	w.mu.Unlock()
}

// appends inline bytes to the output buffer
func (w *BufIoWriter) AppendInline(p []byte) {
	w.mu.Lock()
	w.buf = append(w.buf, '+')
	w.buf = append(w.buf, p...)
	w.buf = append(w.buf, BinCRLF...)
	w.mu.Unlock()
}

// appends an inline string to the output buffer
func (w *BufIoWriter) AppendInlineString(s string) {
	w.mu.Lock()
	w.buf = append(w.buf, '+')
	w.buf = append(w.buf, s...)
	w.buf = append(w.buf, BinCRLF...)
	w.mu.Unlock()
}

// appends an error message to the output buffer
func (w *BufIoWriter) AppendError(msg string) {
	w.mu.Lock()
	w.buf = append(w.buf, '-')
	w.buf = append(w.buf, msg...)
	w.buf = append(w.buf, BinCRLF...)
	w.mu.Unlock()
}

func (w *BufIoWriter) AppendErrorf(pattern string, args ...interface{}) {
	w.AppendError(fmt.Sprintf(pattern, args...))
}

// appends a numeric response to the output buffer
func (w *BufIoWriter) AppendInt(n int64) {
	w.mu.Lock()
	switch n {
	case 0:
		w.buf = append(w.buf, BinZERO...)
	case 1:
		w.buf = append(w.buf, BinONE...)
	default:
		w.buf = append(w.buf, ':')
		w.buf = append(w.buf, strconv.FormatInt(n, 10)...)
		w.buf = append(w.buf, BinCRLF...)
	}
	w.mu.Unlock()
}

// appends a nil-value to the output
func (w *BufIoWriter) AppendNil() {
	w.mu.Lock()
	w.buf = append(w.buf, BinNIL...)
	w.mu.Unlock()
}

// appends OK to the output buffer
func (w *BufIoWriter) AppendOK() {
	w.mu.Lock()
	w.buf = append(w.buf, BinOK...)
	w.mu.Unlock()
}

// flush pending buffer
func (w *BufIoWriter) Flush() error {
	w.mu.Lock()
	err := w.flush()
	defer w.mu.Unlock()
	return err
}

// resets the writer with an new interface
func (w *BufIoWriter) Reset(writer io.Writer) {
	w.reset(MkStdBuffer(), writer)
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
	case re.ProtoError:
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
			if s.Len() == 0 {
				w.AppendError(re.ErrEmptyListOrSet.Error())
				return nil
			}
			w.AppendArrayLen(s.Len())
			for i := 0; i < s.Len(); i++ {
				w.Append(s.Index(i).Interface())
			}
		case reflect.Map:
			s := reflect.ValueOf(v)
			if s.Len() == 0 {
				w.AppendError(re.ErrEmptyListOrSet.Error())
				return nil
			}
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
func (w *BufIoWriter) AppendRawString(rawInput []byte) {
	w.mu.Lock()
	defer w.mu.Unlock()
	loggers.Debug("raw string %s", string(rawInput))
	w.buf = append(w.buf, rawInput...)
}

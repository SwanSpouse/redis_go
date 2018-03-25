package networking

import (
	"io"
	"sync"
	"strconv"
	"fmt"
	"strings"
	"reflect"
)

type bufIoWriter struct {
	io.Writer
	buf []byte
	mu  sync.Mutex
}

// returns the number of buffered bytes
func (b *bufIoWriter) Buffered() int {
	b.mu.Lock()
	n := len(b.buf)
	b.mu.Unlock()
	return n
}

func (b *bufIoWriter) appendSize(c byte, n int64) {
	b.buf = append(b.buf, c)
	b.buf = append(b.buf, strconv.FormatInt(n, 10)...)
	b.buf = append(b.buf, binCRLF...)
}

func (b *bufIoWriter) reset(buf []byte, wr io.Writer) {
	*b = bufIoWriter{buf: buf[:0], Writer: wr}
}

func (b *bufIoWriter) flush() error {
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
func (b *bufIoWriter) AppendArrayLen(n int) {
	b.mu.Lock()
	b.appendSize('*', int64(n))
	b.mu.Unlock()
}

// appends bulk bytes to the output buffer
func (b *bufIoWriter) AppendBulk(p []byte) {
	b.mu.Lock()
	b.appendSize('$', int64(len(p)))
	b.buf = append(b.buf, p...)
	b.buf = append(b.buf, binCRLF...)
	b.mu.Unlock()
}

// appends a bulk string to the output buffer
func (b *bufIoWriter) AppendBulkString(s string) {
	b.mu.Lock()
	b.appendSize('$', int64(len(s)))
	b.buf = append(b.buf, s ...)
	b.buf = append(b.buf, binCRLF...)
	b.mu.Unlock()
}

// appends inline bytes to the output buffer
func (b *bufIoWriter) AppendInline(p []byte) {
	b.mu.Lock()
	b.buf = append(b.buf, '+')
	b.buf = append(b.buf, p...)
	b.buf = append(b.buf, binCRLF...)
	b.mu.Unlock()
}

// appends an inline string to the output buffer
func (b *bufIoWriter) AppendInlineString(s string) {
	b.mu.Lock()
	b.buf = append(b.buf, '+')
	b.buf = append(b.buf, s...)
	b.buf = append(b.buf, binCRLF...)
	b.mu.Unlock()
}

// appends an error message to the output buffer
func (b *bufIoWriter) AppendError(msg string) {
	b.mu.Lock()
	b.buf = append(b.buf, '-')
	b.buf = append(b.buf, msg...)
	b.buf = append(b.buf, binCRLF...)
	b.mu.Unlock()
}

func (b *bufIoWriter) AppendErrorf(pattern string, args ...interface{}) {
	b.AppendError(fmt.Sprintf(pattern, args...))
}

// appends a numeric response to the output buffer
func (b *bufIoWriter) AppendInt(n int64) {
	b.mu.Lock()
	switch n {
	case 0:
		b.buf = append(b.buf, binZERO...)
	case 1:
		b.buf = append(b.buf, binONE...)
	default:
		b.buf = append(b.buf, ':')
		b.buf = append(b.buf, strconv.FormatInt(n, 10)...)
		b.buf = append(b.buf, binCRLF...)
	}
	b.mu.Unlock()
}

// appends a nil-value to the output
func (b *bufIoWriter) AppendNil() {
	b.mu.Lock()
	b.buf = append(b.buf, binNIL...)
	b.mu.Unlock()
}

// appends OK to the output buffer
func (b *bufIoWriter) AppendOK() {
	b.mu.Lock()
	b.buf = append(b.buf, binNIL...)
	b.mu.Unlock()
}

// flush pending buffer
func (b *bufIoWriter) Flush() error {
	b.mu.Lock()
	err := b.flush()
	defer b.mu.Unlock()
	return err
}

// resets the writer with an new interface
func (b *bufIoWriter) Reset(w io.Writer) {
	b.reset(b.buf, w)
}

// Append implements ResponseWriter
func (w *bufIoWriter) Append(v interface{}) error {
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

package resp

import (
	"fmt"
	"io"
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

type protoError string

func (p protoError) Error() string { return string(p) }

func protoErrorf(m string, args ...interface{}) error {
	return protoError(fmt.Sprintf(m, args...))
}

// IsProtocolError returns true if the error is a protocol error
func IsProtocolError(err error) bool {
	_, ok := err.(protoError)
	return ok
}

var (
	binCRLF = []byte("\r\n")
	binOK   = []byte("+OK\r\n")
	binZERO = []byte(":0\r\n")
	binONE  = []byte(":1\r\n")
	binNIL  = []byte("$-1\r\n")
)

const (
	errInvalidMultiBulkLength = protoError("Protocol error: invalid multibulk length")
	errInvalidBulkLength      = protoError("Protocol error: invalid bulk length")
	errBlankBulkLength        = protoError("Protocol error: expected '$', got ' '")
	errInlineRequestTooLong   = protoError("Protocol error: too big inline request")
	errNotANumber             = protoError("Protocol error: expected a number")
	errNotANilMessage         = protoError("Protocol error: expected a nil")
	errBadResponseType        = protoError("Protocol error: bad response type")
)

// MaxBufferSize is the max request/response buffer size
const MaxBufferSize = 64 * 1024

// len of \r\n
const CRLFLen = 2

func mkStdBuffer() []byte { return make([]byte, MaxBufferSize) }

type ResponseWriter interface {
	io.Writer
	// AppendArrayLen appends an array header to the output buffer.
	AppendArrayLen(n int)
	// AppendBulk appends bulk bytes to the output buffer.
	AppendBulk(p []byte)
	// AppendBulkString appends a bulk string to the output buffer.
	AppendBulkString(s string)
	// AppendInline appends inline bytes to the output buffer.
	AppendInline(p []byte)
	// AppendInlineString appends an inline string to the output buffer.
	AppendInlineString(s string)
	// AppendError appends an error message to the output buffer.
	AppendError(msg string)
	// AppendErrorf appends an error message to the output buffer.
	AppendErrorf(pattern string, args ...interface{})
	// AppendInt appends a numeric response to the output buffer.
	AppendInt(n int64)
	// AppendNil appends a nil-value to the output buffer.
	AppendNil()
	// AppendOK appends "OK" to the output buffer.
	AppendOK()
	// Append automatically serialized given values and appends them to the output buffer.
	// Supported values include:
	//   * nil
	//   * error
	//   * string
	//   * []byte
	//   * bool
	//   * float32, float64
	//   * int, int8, int16, int32, int64
	//   * uint, uint8, uint16, uint32, uint64
	//   * CustomResponse instances
	//   * slices and maps of any of the above
	Append(v interface{}) error

	// CopyBulk copies n bytes from a reader.
	// This call may flush pending buffer to prevent overflows.
	//CopyBulk(src io.Reader, n int64) error

	// Buffered returns the number of pending bytes.
	Buffered() int
	// Flush flushes pending buffer.
	Flush() error
	// Reset resets the writer to a new writer and recycles internal buffers.
	Reset(w io.Writer)
}

// NewResponseWriter wraps any writer interface, but
// normally a net.Conn.
func NewResponseWriter(wr io.Writer) ResponseWriter {
	w := new(bufIoWriter)
	w.reset(mkStdBuffer(), wr)
	return w
}

// ResponseParser is a basic response parser
type ResponseParser interface {
	// PeekType returns the type of the next response block
	PeekType() (ResponseType, error)
	// ReadNil reads a nil value
	ReadNil() error
	// ReadBulkString reads a bulk and returns a string
	ReadBulkString() (string, error)
	// ReadBulk reads a bulk and returns bytes (optionally appending to a passed p buffer)
	ReadBulk(p []byte) ([]byte, error)
	// ReadInt reads an int value
	ReadInt() (int64, error)
	// ReadArrayLen reads the array length
	ReadArrayLen() (int, error)
	// ReadError reads an error string
	ReadError() (string, error)
	// ReadInlineString reads a status string
	ReadInlineString() (string, error)
	// Scan scans results into the given values.
	Scan(vv ...interface{}) error
}

// ResponseReader is used by clients to wrap a server connection and
// parse responses.
type ResponseReader interface {
	ResponseParser

	// Buffered returns the number of buffered (unread) bytes.
	Buffered() int
	// Reset resets the reader to a new reader and recycles internal buffers.
	Reset(r io.Reader)
}

// NewResponseReader returns ResponseReader, which wraps any reader interface, but
// normally a net.Conn.
func NewResponseReader(rd io.Reader) ResponseReader {
	r := new(bufIoReader)
	r.reset(mkStdBuffer(), rd)
	return r
}

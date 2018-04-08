package protocol

import (
	"io"
	"redis_go/tcp"
)

/**
    @author lmj
    ResponseReader 主要用于mock测试，mock client利用ResponseReader读取server发送的请求。
并对结果进行验证。
*/
type ResponseReader struct {
	r *tcp.BufIoReader
}

func NewResponseReader(rd io.Reader) *ResponseReader {
	r := new(tcp.BufIoReader)
	r.Reset(rd)
	return &ResponseReader{r: r}
}

func (r *ResponseReader) PeekType() (tcp.ResponseType, error) {
	return r.r.PeekType()
}

func (r *ResponseReader) ReadNil() error {
	return r.r.ReadNil()
}

func (r *ResponseReader) ReadBulkString() (string, error) {
	return r.r.ReadBulkString()
}

func (r *ResponseReader) ReadBulk(p []byte) ([]byte, error) {
	return r.r.ReadBulk(p)
}

func (r *ResponseReader) ReadInt() (int64, error) {
	return r.r.ReadInt()
}

func (r *ResponseReader) ReadArrayLen() (int, error) {
	return r.r.ReadArrayLen()
}

func (r *ResponseReader) ReadError() (string, error) {
	return r.r.ReadError()
}

func (r *ResponseReader) ReadInlineString() (string, error) {
	return r.r.ReadInlineString()
}

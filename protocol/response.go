package protocol

import (
	"fmt"
	"io"
	"redis_go/loggers"
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

func (r *ResponseReader) Reset(rd io.Reader) {
	r.r.Reset(rd)
	r.r = nil
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

func (r *ResponseReader) Read() ([]string, error) {
	ret := make([]string, 0)
	responseType, err := r.PeekType()
	if err != nil {
		return nil, err
	}

	switch responseType {
	case tcp.TypeInt:
		val, err := r.ReadInt()
		if err != nil {
			return ret, err
		}
		ret = append(ret, fmt.Sprintf("%d", val))
	case tcp.TypeInline:
		val, err := r.ReadInlineString()
		if err != nil {
			return ret, err
		}
		ret = append(ret, val)
	case tcp.TypeError:
		val, err := r.ReadError()
		if err != nil {
			return ret, err
		}
		ret = append(ret, val)
	case tcp.TypeBulk:
		val, err := r.ReadBulkString()
		if err != nil {
			return ret, err
		}
		ret = append(ret, string(val))
	case tcp.TypeNil:
		err := r.ReadNil()
		if err != nil {
			return ret, err
		}
		ret = append(ret, "NIL")
	case tcp.TypeArray:
		arrayLen, err := r.ReadArrayLen()
		if err != nil {
			return ret, err
		}
		for i := 0; i < arrayLen; i++ {
			tempRet, err := r.Read()
			if err != nil {
				return ret, err
			}
			ret = append(ret, tempRet...)
		}
	case tcp.TypeUnknown:
		loggers.Errorf("unknown type %v", responseType)
	}
	return ret, nil
}

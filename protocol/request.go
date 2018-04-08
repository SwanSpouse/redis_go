package protocol

import (
	"io"
	re "redis_go/error"
	"redis_go/tcp"
)

/**
    @author lmj
    RequestWriter 主要用于mock测试，mock client利用RequestWriter构造请求
 */
type RequestWriter struct {
	w *tcp.BufIoWriter
}

// NewRequestWriter wraps any writer interface
func NewRequestWriter(wr io.Writer) *RequestWriter {
	w := new(tcp.BufIoWriter)
	w.Reset(wr)
	return &RequestWriter{w: w}
}

func (w *RequestWriter) Reset(wr io.Writer) {
	w.w.Reset(wr)
}

func (w *RequestWriter) Buffered() int {
	return w.w.Buffered()
}

func (w *RequestWriter) Flush() error {
	return w.w.Flush()
}

func (w *RequestWriter) WriteCmd(cmd string, args ...[]byte) {
	w.w.AppendArrayLen(len(args) + 1)
	w.w.AppendBulkString(cmd)
	for _, arg := range args {
		w.w.AppendBulk(arg)
	}
}

func (w *RequestWriter) WriteCmdString(cmd string, args ...string) {
	w.w.AppendArrayLen(len(args) + 1)
	w.w.AppendBulkString(cmd)
	for _, arg := range args {
		w.w.AppendBulkString(arg)
	}
}

func (w *RequestWriter) WriteMultiBulkSize(n int) error {
	if n < 0 {
		return re.ErrInvalidMultiBulkLength
	}
	w.w.AppendArrayLen(n)
	return nil
}

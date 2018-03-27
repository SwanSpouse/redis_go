package log

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"sync"
	"time"
	"strings"
)

type Logger struct {
	mu              sync.Mutex
	prefix          string
	out             io.Writer
	showGoRoutineId bool
}

func New(out io.Writer, prefix string) *Logger {
	return &Logger{out: out, prefix: prefix}
}

func (l *Logger) header(tm time.Time, file string, line int, s string) string {
	ms := tm.Nanosecond() / int(time.Millisecond)
	goRoutineId := ""
	if l.showGoRoutineId {
		goRoutineId = "goRoutineId" + GetGID()
	}
	return fmt.Sprintf("%s.%03d %s %s file %s line %d ", tm.Format("2006-01-02 15:04:05"), ms, l.prefix, goRoutineId, file, line)
}

func (l *Logger) Output(callDepth int, s string) error {
	s = strings.Replace(s, "\r\n", "\\r\\n", -1)
	now := time.Now() // get this early.
	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	var ok bool
	_, file, line, ok = runtime.Caller(callDepth)
	if !ok {
		file = "???"
		line = 0
	}

	head := l.header(now, file, line, s)
	buf := make([]byte, 0, len(head))
	buf = append(buf, head...)
	for _, c := range []byte(s) {
		if c != '\n' {
			buf = append(buf, c)
		} else {
			buf = append(buf, '\n')
			_, err := l.out.Write(buf)
			if err != nil {
				return err
			}
			buf = buf[:0]
			buf = append(buf, head...)
		}
	}
	if len(buf) > len(head) {
		buf = append(buf, '\n')
		_, err := l.out.Write(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetGID() string {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	endIndex := bytes.IndexByte(b, ' ')
	if endIndex >= 0 {
		return string(b[:endIndex])
	}
	return "-1"
}

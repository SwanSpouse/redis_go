package aof

import (
	"bufio"
	"errors"
	"os"
	re "redis_go/error"
	"redis_go/util"
	"strconv"
	"strings"
)

type Decoder struct {
	rd *bufio.Reader
}

func NewDecoder(filename string) *Decoder {
	if !util.FileExists(filename) {
		return nil
	}
	f, err := os.Open(filename)
	if err != nil || f == nil {
		return nil
	}
	return &Decoder{rd: bufio.NewReader(f)}
}

type CmdOutput struct {
	Argc int
	Argv []string
}

func (decoder *Decoder) DecodeAppendOnlyFile() (*CmdOutput, error) {
	out := &CmdOutput{
		Argc: 0,
		Argv: make([]string, 0),
	}
	argc, err := decoder.readMultiBulkLength()
	if err != nil {
		return nil, err
	}
	for i := 0; i < argc; i++ {
		argvLength, err := decoder.readBulkLength()
		if err != nil {
			return nil, err
		}
		argv, err := decoder.readCmdArgv(argvLength)
		if err != nil {
			return nil, err
		}
		out.Argv = append(out.Argv, string(argv))
	}
	out.Argc = argc
	return out, nil
}

func (decoder *Decoder) readMultiBulkLength() (int, error) {
	line, err := decoder.rd.ReadString('\n')
	if err != nil {
		return 0, err
	}
	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r")
	if len(line) < 2 || line[0] != '*' {
		return 0, re.ErrInvalidMultiBulkLength
	}
	count, err := strconv.Atoi(line[1:])
	if err != nil {
		return 0, re.ErrInvalidMultiBulkLength
	}
	return count, nil
}

func (decoder *Decoder) readBulkLength() (int, error) {
	line, err := decoder.rd.ReadString('\n')
	if err != nil {
		return 0, err
	}
	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r")
	if len(line) < 2 || line[0] != '$' {
		return 0, re.ErrInvalidMultiBulkLength
	}
	length, err := strconv.Atoi(line[1:])
	if err != nil {
		return 0, re.ErrInvalidMultiBulkLength
	}
	return length, nil
}

func (decoder *Decoder) readCmdArgv(length int) ([]byte, error) {
	if length <= 0 {
		return nil, errors.New("can not read length 0")
	}
	outBuf := make([]byte, length+2)
	readLength, err := decoder.rd.Read(outBuf)
	if err != nil {
		return nil, err
	}
	if length+2 != readLength {
		return nil, re.ErrInvalidBulkLength
	}
	// trim tail \r\n
	return outBuf[:len(outBuf)-2], nil
}

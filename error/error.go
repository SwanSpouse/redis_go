package error

import (
	"fmt"
)

type protoError string

func (p protoError) Error() string { return string(p) }

func ProtoErrorf(m string, args ...interface{}) error {
	return protoError(fmt.Sprintf(m, args...))
}

// IsProtocolError returns true if the error is a protocol error
func IsProtocolError(err error) bool {
	_, ok := err.(protoError)
	return ok
}

const (
	ErrInvalidMultiBulkLength = protoError("ERR Protocol error: invalid multibulk length")
	ErrInvalidBulkLength      = protoError("ERR Protocol error: invalid bulk length")
	ErrBlankBulkLength        = protoError("ERR Protocol error: expected '$', got ' '")
	ErrInlineRequestTooLong   = protoError("ERR Protocol error: too big inline request")
	ErrNotANumber             = protoError("ERR Protocol error: expected a number")
	ErrNotANilMessage         = protoError("ERR Protocol error: expected a nil")
	ErrBadResponseType        = protoError("ERR Protocol error: bad response type")
	ErrUnknown                = protoError("ERR Protocol error: unknown")
	ErrImpossible             = protoError("ERR Protocol error: impossible")

	ErrNotIntegerOrOutOfRange = protoError("value is not an integer or out of range")
	ErrWrongNumberOfArgs      = protoError("wrong number of arguments for '%s' command")
	ErrUnknownCommand         = protoError("ERR unknown command '%s'")
	ErrNilCommand             = protoError("ERR nil command")
	ErrFunctionNotImplement   = protoError("This command has not been implement.")
	ErrWrongType              = protoError("WRONGTYPE Operation against a key holding the wrong kind of value")
	ErrWrongTypeOrEncoding    = protoError("error object type or encoding. type:%s, encoding:%s")
	ErrConvertToTargetType    = protoError("ERR cannot convert tbase to target type")
	ErrConvertEncoding        = protoError("Err convert encoding")
	ErrNilValue               = protoError("Err nil")
	ErrNoSuchKey              = protoError("ERR no such key")
	ErrIncrOrDecrOverflow     = protoError("ERR increment or decrement would overflow")
	ErrValueIsNotFloat        = protoError("ERR value is not a valid float")
	ErrEmptyListOrSet         = protoError("(empty list or set)")
	ErrSyntaxError            = protoError("ERR syntax error")
	ErrRedisRdbSaveInProcess  = protoError("ERR redis rdb save is in process")
	ErrAofFormat              = protoError("Bad file format reading the append only file: make a backup of your AOF file, then use ./redis-check-aof --fix <filename>")
)

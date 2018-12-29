package error

import (
	"fmt"
)

type ProtoError string

func (p ProtoError) Error() string {
	return string(p)
}

func ProtoErrorf(m string, args ...interface{}) error {
	return ProtoError(fmt.Sprintf(m, args...))
}

// IsProtocolError returns true if the error is a protocol error
func IsProtocolError(err error) bool {
	_, ok := err.(ProtoError)
	return ok
}

const (
	ErrInvalidMultiBulkLength = ProtoError("ERR Protocol error: invalid multibulk length")
	ErrInvalidBulkLength      = ProtoError("ERR Protocol error: invalid bulk length")
	ErrBlankBulkLength        = ProtoError("ERR Protocol error: expected '$', got ' '")
	ErrInlineRequestTooLong   = ProtoError("ERR Protocol error: too big inline request")
	ErrNotANumber             = ProtoError("ERR Protocol error: expected a number")
	ErrNotANilMessage         = ProtoError("ERR Protocol error: expected a nil")
	ErrBadResponseType        = ProtoError("ERR Protocol error: bad response type")
	ErrUnknown                = ProtoError("ERR Protocol error: unknown")
	ErrImpossible             = ProtoError("ERR Protocol error: impossible")

	ErrNotIntegerOrOutOfRange = ProtoError("value is not an integer or out of range")
	ErrWrongNumberOfArgs      = ProtoError("wrong number of arguments for '%s' command")
	ErrUnknownCommand         = ProtoError("ERR unknown command '%s'")
	ErrNilCommand             = ProtoError("ERR nil command")
	ErrFunctionNotImplement   = ProtoError("This command has not been implement.")
	ErrWrongType              = ProtoError("WRONGTYPE Operation against a key holding the wrong kind of value")
	ErrWrongTypeOrEncoding    = ProtoError("error object type or encoding. type:%s, encoding:%s")
	ErrConvertToTargetType    = ProtoError("ERR cannot convert tbase to target type")
	ErrConvertEncoding        = ProtoError("Err convert encoding")
	ErrNilValue               = ProtoError("Err nil")
	ErrNoSuchKey              = ProtoError("ERR no such key")
	ErrIncrOrDecrOverflow     = ProtoError("ERR increment or decrement would overflow")
	ErrValueIsNotFloat        = ProtoError("ERR value is not a valid float")
	ErrEmptyListOrSet         = ProtoError("(empty list or set)")
	ErrSyntaxError            = ProtoError("ERR syntax error")
	ErrRedisRdbSaveInProcess  = ProtoError("ERR redis rdb save is in process")
	ErrAofFormat              = ProtoError("Bad file format reading the append only file: make a backup of your AOF file, then use ./redis-check-aof --fix <filename>")
	ErrPubSubCommand          = ProtoError("ERR Unknown PUBSUB subcommand or wrong number of arguments for %s")
)

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
	ErrInvalidMultiBulkLength = protoError("Protocol error: invalid multibulk length")
	ErrInvalidBulkLength      = protoError("Protocol error: invalid bulk length")
	ErrBlankBulkLength        = protoError("Protocol error: expected '$', got ' '")
	ErrInlineRequestTooLong   = protoError("Protocol error: too big inline request")
	ErrNotANumber             = protoError("Protocol error: expected a number")
	ErrNotANilMessage         = protoError("Protocol error: expected a nil")
	ErrBadResponseType        = protoError("Protocol error: bad response type")
	ErrUnknown                = protoError("Protocol error: unknown")

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
	ErrIncrOrDecrOverflow     = protoError("ERR increment or decrement would overflow")
	ErrValueIsNotFloat        = protoError("ERR value is not a valid float")
)

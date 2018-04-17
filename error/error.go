package error

import "fmt"

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
)

const (
	ErrWrongNumberOfArgs      = "wrong number of arguments for '%s' command"
	ErrFunctionNotImplement   = "This command has not been implement."
	ErrNotIntegerOrOutOfRange = "value is not an integer or out of range"
)

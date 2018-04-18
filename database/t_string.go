package database

import (
	"redis_go/encodings"
	"strconv"
)

var (
	_ TString = (*encodings.StringRaw)(nil)
	_ TString = (*encodings.StringInt)(nil)
	_ TString = (*encodings.StringEmb)(nil)
)

type TString interface {
	// common operation
	GetObjectType() string
	SetObjectType(string)
	GetEncoding() string
	SetEncoding(string)
	GetLRU() int
	SetLRU(int)
	GetRefCount() int
	IncrRefCount() int
	DecrRefCount() int
	GetTTL() int
	SetTTL(int)
	GetValue() interface{}
	SetValue(interface{})
	IsExpired() bool

	// string command operation
	Append(string) int
	Incr() (int, error)
	Decr() (int, error)
	IncrBy(int) (int, error)
	DecrBy(int) (int, error)
	Strlen() int
}

// 创建一个新的redis string object
func NewRedisStringObject(value string) TBase {
	return NewRedisStringObjectWithTTL(value, -1)
}

// 创建一个新的带有ttl的redis string object
func NewRedisStringObjectWithTTL(value string, ttl int) TBase {
	if valueInt, err := strconv.Atoi(value); err == nil {
		return encodings.NewRedisStringWithEncodingRawInt(valueInt, ttl)
	}
	return encodings.NewRedisStringWithEncodingRawString(value, ttl)
}

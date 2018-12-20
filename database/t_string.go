package database

import (
	"github.com/SwanSpouse/redis_go/encodings"
)

var (
	// string对象的三种实现方式
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
	Incr() (int64, error)
	Decr() (int64, error)
	IncrBy(string) (int64, error)
	DecrBy(string) (int64, error)
	IncrByFloat(string) (string, error)
	Strlen() int
	String() string
}

func NewRedisStringWithEncodingRawString(value string, ttl int) TString {
	return encodings.NewStringRaw(ttl, value)
}

func NewRedisStringWithEncodingStringInt(value int64, ttl int) TString {
	return encodings.NewStringInt(ttl, value)
}

func NewRedisStringWithEncodingStringEmb(value int, ttl int) TString {
	return encodings.NewStringEmb(ttl, value)
}

// 创建一个新的redis string object
func NewRedisStringObject(value string) TBase {
	return NewRedisStringObjectWithTTL(value, -1)
}

// 创建一个新的带有ttl的redis string object
func NewRedisStringObjectWithTTL(value string, ttl int) TBase {
	// 先默认都是raw string
	//if valueInt, err := strconv.ParseInt(value, 10, 64); err == nil {
	//	return NewRedisStringWithEncodingStringInt(valueInt, ttl)
	//}
	return NewRedisStringWithEncodingRawString(value, ttl)
}

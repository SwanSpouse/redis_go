package database

import (
	"github.com/SwanSpouse/redis_go/encodings"
)

var (
	// hash对象的实现方式
	_ THash = (*encodings.HashDict)(nil)
)

type THash interface {
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

	// hash command operation
	HSet(string, string) int
	HGet(string) (string, error)
	HExists(string) int
	HDel([]string) int
	HLen() int
	HGetAll() []string
	HIncrBy(string, string) (string, error)
	HIncrByFloat(string, string) (string, error)
	HDebug()
}

func NewRedisHashObject() TBase {
	return NewRedisHashObjectWithTTL(-1)
}

func NewRedisHashObjectWithTTL(ttl int) TBase {
	return encodings.NewHashDict(ttl)
}

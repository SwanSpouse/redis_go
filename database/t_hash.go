package database

import (
	"redis_go/encodings"
)

var (
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
	HDebug()
}

func NewRedisHashObject() TBase {
	return NewRedisHashObjectWithTTL(-1)
}

func NewRedisHashObjectWithTTL(ttl int) TBase {
	return encodings.NewHashDict(ttl)
}

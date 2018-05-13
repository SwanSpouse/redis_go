package database

import "redis_go/encodings"

var (
	_ TSet = (*encodings.HashSet)(nil)
)

type TSet interface {
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
	String() string

	// set command operation
	SAdd([]string) int
	SCard() int
	SIsMember(string) int
	SMembers() []string
	SPop() (string, error)
	SRandMember() (string, error)
	SRem([]string) int
	SDebug()
}

func NewRedisSetObject() TBase {
	return NewRedisSetObjectWithTTL(-1)
}

func NewRedisSetObjectWithTTL(ttl int) TBase {
	return encodings.NewHasSet(ttl)
}

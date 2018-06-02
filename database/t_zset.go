package database

import "redis_go/encodings"

var (
	_ TZSet = (*encodings.SortedSet)(nil)
)

type TZSet interface {
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

	// sorted set command operation
	ZAdd([]string) (int, error)
	ZCard() int
	ZCount(string, string) (int, error)
	ZIncrBy(string, string) (float64, error)
	ZRange(string, string) ([]string, error)
	ZRangeByScore(string, string) ([]string, error)
	ZRevRange(string, string) ([]string, error)
	ZRevRangeByScore(string, string) ([]string, error)
	ZRank(string) (int, error)
	ZRevRank(string) (int, error)
	ZRem([]string) int
	ZRemRangeByRank(string, string) (int, error)
	ZRemRangeByScore(string, string) (int, error)
	ZScore(string) (float64, error)
	//ZUnionStore()
	//ZInterStore()
	//ZScan()
}

func NewRedisSortedSetObject() TZSet {
	return NewRedisSortedSetObjectWithTTL(-1)
}

func NewRedisSortedSetObjectWithTTL(ttl int) TZSet {
	return encodings.NewSortedSet(ttl)
}

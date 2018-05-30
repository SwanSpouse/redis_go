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
	ZAdd([]string) int
	ZCard() int
	ZCount() int
	ZIncrBy(string, float64) float64
	ZRange(int, int) []string
	ZRangeByScore(float64, float64) []string
	ZRevRange(int, int) []string
	ZRevRangeByScore(float64, float64) []string
	ZRank(string) int
	ZRevRank(string) int
	ZRem([]string) int
	ZRemRangeByRyRank(int, int) int
	ZRemRangeByScore(float64, float64) int
	ZScore(string) float64
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

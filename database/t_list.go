package database

import "redis_go/encodings"

var (
	// list对象的实现方式
	_ TList = (*encodings.ListLinkedList)(nil)
)

type TList interface {
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

	// list command operation
	LPush([]string) int
	RPush([]string) int
	LPop() string
	RPop() string
	LRange(int, int) []string
	LIndex(int) (string, error)
	LLen() int
	LInsert(int, string, ...string) (int, error)
	LRem(int, string) int
	LTrim(int, int) error
	LSet(int, string) error
	GetAllMembers() []string
	Debug()
}

func NewRedisListWithEncodingLinkedList(ttl int) TBase {
	return encodings.NewListLinkedList(ttl)
}

// 创建一个新的redis list object
func NewRedisListObject() TBase {
	return NewRedisListObjectWithTTL(-1)
}

// 创建一个新的带有ttl的redis list object
func NewRedisListObjectWithTTL(ttl int) TBase {
	return NewRedisListWithEncodingLinkedList(ttl)
}

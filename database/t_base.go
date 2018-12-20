package database

import "github.com/SwanSpouse/redis_go/encodings"

var (
	_ TBase = (*encodings.StringRaw)(nil)
	_ TBase = (*encodings.StringInt)(nil)
	_ TBase = (*encodings.StringEmb)(nil)

	_ TBase = (*encodings.ListLinkedList)(nil)

	_ TBase = (*encodings.HashDict)(nil)

	_ TBase = (*encodings.HashSet)(nil)

	_ TBase = (*encodings.SortedSet)(nil)
)

type TBase interface {
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
}

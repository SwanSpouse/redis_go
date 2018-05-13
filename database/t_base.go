package database

import "redis_go/encodings"

var (
	_ TBase = (*encodings.StringRaw)(nil)
	_ TBase = (*encodings.StringInt)(nil)
	_ TBase = (*encodings.StringEmb)(nil)

	_ TList = (*encodings.ListLinkedList)(nil)

	_ THash = (*encodings.HashDict)(nil)

	_ TSet = (*encodings.HashSet)(nil)
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

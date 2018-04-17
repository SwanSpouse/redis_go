package database

import "time"

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
	Incr() int
	Decr() int
	IncrBy(int) int
	DecrBy(int) int
	Strlen() int
}

// 创建一个新的redis string object
func NewRedisStringObject(value string) (TBase, error) {
	return NewRedisStringObjectWithTTL(value, -1)
}

/*
 *	创建一个新的带有ttl的redis string object
 */
func NewRedisStringObjectWithTTL(value string, ttl int) (TBase, error) {
	obj := &RedisObject{
		objectType: RedisTypeString,
		encoding:   RedisEncodingRaw, // 暂时都默认为raw吧。
		ttl:        ttl,
		value:      value,
	}
	if ttl > 0 {
		obj.expireTime = time.Now().Add(time.Duration(ttl) * time.Second)
	}
	return obj, nil
}

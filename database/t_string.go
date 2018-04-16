package database

import "time"

type TString interface {
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
}

// 创建一个新的redis string object
func NewRedisStringObject(value string) (TBase, error) {
	return NewRedisStringObjectWithTTL(value, -1)
}

/*
 *	创建一个新的带有ttl的redis string object
 *      //TODO lmj 这里需要根据变量的值来进行判断，看是创建什么样encoding的redis object，
 *      //TODO lmj 这里的value应该传进来一个interface? 还是看根据是否能够转换成数字来判断?
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

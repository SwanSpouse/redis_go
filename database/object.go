package database

import (
	"time"
)

type RedisObject struct {
	objectType string      // 类型
	encoding   string      // 编码
	lru        int         // LRU时间
	refCount   int         // 引用计数
	ttl        int         // ttl
	expireTime time.Time   // 过期时间
	value      interface{} // 指向的对象
}

// @deprecated
func NewRedisObject(value interface{}) TBase {
	return &RedisObject{
		value: value,
	}
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

func (obj *RedisObject) GetObjectType() string {
	return obj.objectType
}

func (obj *RedisObject) SetObjectType(objType string) {
	obj.objectType = objType
}

func (obj *RedisObject) GetEncoding() string {
	return obj.encoding
}

func (obj *RedisObject) SetEncoding(encoding string) {
	obj.encoding = encoding
}

func (obj *RedisObject) GetLRU() int {
	return obj.lru
}

func (obj *RedisObject) SetLRU(lru int) {
	obj.lru = lru
}

func (obj *RedisObject) GetRefCount() int {
	return obj.refCount
}

func (obj *RedisObject) IncrRefCount() int {
	obj.refCount += 1
	return obj.refCount
}

func (obj *RedisObject) DecrRefCount() int {
	obj.refCount -= 1
	return obj.refCount
}

func (obj *RedisObject) GetTTL() int {
	return obj.ttl
}

func (obj *RedisObject) SetTTL(ttl int) {
	obj.ttl = ttl
}

func (obj *RedisObject) GetValue() interface{} {
	return obj.value
}

func (obj *RedisObject) SetValue(value interface{}) {
	obj.value = value
}

func (obj *RedisObject) IsExpired() bool {
	// 如果过期时间是有效值，并且当前时间在过期时间之后，说明已经过期。
	if !obj.expireTime.IsZero() && time.Now().After(obj.expireTime) {
		return true
	}
	return false
}

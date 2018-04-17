package raw_type

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

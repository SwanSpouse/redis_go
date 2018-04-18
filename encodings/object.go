package encodings

import (
	"time"
)

const (
	/**
	RedisTypeString  ->  RedisEncodingInt	 	: 使用整数值实现的字符串对象
	RedisTypeString  ->  RedisEncodingEmbStr        : 使用embstr编码的简单动态字符串实现的字符串对象
	RedisTypeString  ->  RedisEncodingRaw		: 使用简单动态字符串实现的字符串对象
	RedisTypeList    ->  RedisEncodingZipList	: 使用压缩列表实现的列表对象
	RedisTypeList    ->  RedisEncodingLinkedList	: 使用双端链表实现的列表对象
	RedisTypeHash    ->  RedisEncodingZipList	: 使用压缩链表实现的列表对象
	RedisTypeHash    ->  RedisEncodingHT		: 使用字典实现的哈希对象
	RedisTypeSet     ->  RedisEncodingIntSet	: 使用整数集合实现的集合对象
	RedisTypeSet     ->  RedisEncodingHT		: 使用字典实现的集合对象
	RedisTypeZSet    ->  RedisEncodingZipList	: 使用压缩链表实现的有序集合对象
	RedisTypeZSet    ->  RedisEncodingSkipList	: 使用跳跃表和字典实现的有序集合对象
	*/

	/* object type */
	RedisTypeString = "string"
	RedisTypeList   = "list"
	RedisTypeHash   = "hash"
	RedisTypeSet    = "set"
	RedisTypeZSet   = "zset"

	/* redis encoding type */
	RedisEncodingInt        = "int"
	RedisEncodingEmbStr     = "embstr"
	RedisEncodingRaw        = "raw"
	RedisEncodingHT         = "hashtable"
	RedisEncodingLinkedList = "linkedlist"
	RedisEncodingZipList    = "ziplist"
	RedisEncodingIntSet     = "intset"
	RedisEncodingSkipList   = "skiplist"
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

package database

type RedisObject struct {
	objectType string      // 类型
	encoding   int         // 编码
	lru        int         // LRU时间
	refCount   int         // 引用计数
	ttl        int         // ttl
	value      interface{} // 指向的对象
}

func NewRedisObject(value interface{}) *RedisObject {
	return &RedisObject{
		value: value,
	}
}

func (obj *RedisObject) GetObjectType() string {
	return obj.objectType
}

func (obj *RedisObject) SetObjectType(objType string) {
	obj.objectType = objType
}

func (obj *RedisObject) GetEncoding() int {
	return obj.encoding
}

func (obj *RedisObject) SetEncoding(encoding int) {
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

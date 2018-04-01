package database

type RedisObject struct {
	objectType int         // 类型
	lru        int         // LRU时间
	refCount   int         // 引用计数
	value      interface{} // 指向的对象
}

func NewRedisObject(value interface{}) *RedisObject {
	return &RedisObject{
		value: value,
	}
}

func (obj *RedisObject) GetValue() string {
	return obj.value.(string)
}

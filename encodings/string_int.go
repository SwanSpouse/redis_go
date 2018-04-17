package encodings

import "time"

type StringInt struct {
	RedisObject
}

func NewRedisStringWithEncodingRawInt(value string, ttl int) *RedisObject {
	obj := &RedisObject{
		objectType: RedisTypeString,
		encoding:   RedisEncodingInt,
		ttl:        ttl,
		value:      value,
	}
	if ttl > 0 {
		obj.expireTime = time.Now().Add(time.Duration(ttl) * time.Second)
	}
	return obj
}

func (si *StringInt) Append(val string) int {
	return 0
}

func (si *StringInt) Incr() (int, error) {
	return 0, nil
}

func (si *StringInt) Decr() (int, error) {
	return 0, nil
}

func (si *StringInt) IncrBy(val int) (int, error) {
	return 0, nil
}

func (si *StringInt) DecrBy(val int) (int, error) {
	return 0, nil
}

func (si *StringInt) Strlen() int {
	return 0
}

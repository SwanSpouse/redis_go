package encodings

import "time"

type StringEmb struct {
	RedisObject
}

func NewRedisStringWithEncodingEmbStr(value string, ttl int) *RedisObject {
	obj := &RedisObject{
		objectType: RedisTypeString,
		encoding:   RedisEncodingEmbStr,
		ttl:        ttl,
		value:      value,
	}
	if ttl > 0 {
		obj.expireTime = time.Now().Add(time.Duration(ttl) * time.Second)
	}
	return obj
}

func (se *StringEmb) Append(val string) int {
	return 0
}

func (se *StringEmb) Incr() (int, error) {
	return 0, nil
}

func (se *StringEmb) Decr() (int, error) {
	return 0, nil
}

func (se *StringEmb) IncrBy(val int) (int, error) {
	return 0, nil
}

func (se *StringEmb) DecrBy(val int) (int, error) {
	return 0, nil
}

func (se *StringEmb) Strlen() int {
	return 0
}

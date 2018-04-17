package encodings

import (
	re "redis_go/error"
	"strconv"
	"time"
)

type StringRaw struct {
	RedisObject
}

func NewRedisStringWithEncodingRawString(value string, ttl int) *RedisObject {
	obj := &RedisObject{
		objectType: RedisTypeString,
		encoding:   RedisEncodingRaw,
		ttl:        ttl,
		value:      value,
	}
	if ttl > 0 {
		obj.expireTime = time.Now().Add(time.Duration(ttl) * time.Second)
	}
	return obj
}

func (sr *StringRaw) Append(val string) int {
	newValue := sr.GetValue().(string) + val
	sr.SetValue(newValue)
	return len(newValue)
}

func (sr *StringRaw) Incr() (int, error) {
	if valueInt, ok := sr.GetValue().(int); !ok {
		return 0, re.ErrNotIntegerOrOutOfRange
	} else {
		sr.SetValue(strconv.Itoa(valueInt + 1))
		return valueInt + 1, nil
	}
}

func (sr *StringRaw) Decr() (int, error) {
	if valueInt, ok := sr.GetValue().(int); !ok {
		return 0, re.ErrNotIntegerOrOutOfRange
	} else {
		sr.SetValue(strconv.Itoa(valueInt - 1))
		return valueInt - 1, nil
	}
}

func (sr *StringRaw) IncrBy(val int) (int, error) {
	return 0, nil
}

func (sr *StringRaw) DecrBy(val int) (int, error) {
	return 0, nil
}

func (sr *StringRaw) Strlen() int {
	return len(sr.GetValue().(string))
}

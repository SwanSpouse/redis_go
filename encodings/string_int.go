package encodings

import (
	re "redis_go/error"
	"redis_go/log"
	"strconv"
	"time"
)

type StringInt struct {
	RedisObject
}

func NewRedisStringWithEncodingRawInt(value int, ttl int) *RedisObject {
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

func (si *StringInt) convertStringIntToStringRaw() (*StringRaw, error) {
	si.SetEncoding(RedisEncodingRaw)
	if valueInt, ok := si.GetValue().(int); !ok {
		return nil, re.ErrConvertEncoding
	} else {
		si.SetValue(strconv.Itoa(valueInt))
		return (*StringRaw)(si), nil
	}
}

func (si *StringInt) Append(val string) int {
	if sr, err := si.convertStringIntToStringRaw(); err != nil {
		log.Errorf("convert string int to string raw error")
		return -1
	} else {
		return sr.Append(val)
	}
}

func (si *StringInt) Incr() (int, error) {
	if valueInt, ok := si.GetValue().(int); !ok {
		return 0, re.ErrNotIntegerOrOutOfRange
	} else {
		si.SetValue(valueInt + 1)
		return valueInt + 1, nil
	}
}

func (si *StringInt) Decr() (int, error) {
	if valueInt, ok := si.GetValue().(int); !ok {
		return 0, re.ErrNotIntegerOrOutOfRange
	} else {
		si.SetValue(valueInt - 1)
		return valueInt - 1, nil
	}
}

func (si *StringInt) IncrBy(val int) (int, error) {
	return 0, nil
}

func (si *StringInt) DecrBy(val int) (int, error) {
	return 0, nil
}

func (si *StringInt) Strlen() int {
	if valueInt, ok := si.GetValue().(int); !ok {
		return -1
	} else {
		return len(strconv.Itoa(valueInt))
	}
}

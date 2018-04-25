package encodings

import (
	"fmt"
	"math"
	re "redis_go/error"
	"strconv"
	"time"
)

type StringInt struct {
	RedisObject
}

func NewStringInt(ttl int, value interface{}) *StringInt {
	var expireTime time.Time
	if ttl > 0 {
		expireTime = time.Now().Add(time.Duration(ttl) * time.Second)
	}
	return &StringInt{
		RedisObject: RedisObject{
			objectType: RedisTypeString,
			encoding:   RedisEncodingInt,
			ttl:        ttl,
			value:      value,
			expireTime: expireTime,
		},
	}
}

func (si *StringInt) String() string {
	return fmt.Sprintf("%d", si.GetValue().(int64))
}

func (si *StringInt) Append(val string) int {
	return 0
}

func (si *StringInt) Incr() (int64, error) {
	if valueInt, ok := si.GetValue().(int64); !ok {
		return 0, re.ErrNotIntegerOrOutOfRange
	} else {
		if valueInt == math.MaxInt64 {
			return 0, re.ErrIncrOrDecrOverflow
		}
		si.SetValue(valueInt + 1)
		return valueInt + 1, nil
	}
}

func (si *StringInt) Decr() (int64, error) {
	if valueInt, ok := si.GetValue().(int64); !ok {
		return 0, re.ErrNotIntegerOrOutOfRange
	} else {
		if valueInt == math.MinInt64 {
			return 0, re.ErrIncrOrDecrOverflow
		}
		si.SetValue(valueInt - 1)
		return valueInt - 1, nil
	}
}

func (si *StringInt) IncrBy(val string) (int64, error) {
	incrValInt, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, re.ErrNotIntegerOrOutOfRange
	}
	if valueInt, ok := si.GetValue().(int64); !ok {
		return 0, re.ErrNotIntegerOrOutOfRange
	} else {
		if (valueInt > 0 && incrValInt > 0 && incrValInt > math.MaxInt64-valueInt) ||
			(valueInt < 0 && incrValInt < 0 && incrValInt < math.MinInt64-valueInt) {
			return 0, re.ErrIncrOrDecrOverflow
		}
		si.SetValue(valueInt + incrValInt)
		return valueInt + incrValInt, nil
	}
}

func (si *StringInt) DecrBy(val string) (int64, error) {
	decrValInt, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, re.ErrNotIntegerOrOutOfRange
	}
	return si.IncrBy(strconv.FormatInt(-1*decrValInt, 10))
}

func (si *StringInt) IncrByFloat(val string) (string, error) {
	return "", re.ErrWrongTypeOrEncoding
}

func (si *StringInt) Strlen() int {
	if valueInt, ok := si.GetValue().(int64); !ok {
		return -1
	} else {
		return len(fmt.Sprintf("%d", valueInt))
	}
}

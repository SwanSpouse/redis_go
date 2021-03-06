package encodings

import (
	"fmt"
	"strconv"
	"time"

	re "github.com/SwanSpouse/redis_go/error"
	"github.com/SwanSpouse/redis_go/util"
)

type StringRaw struct {
	RedisObject
}

func NewStringRaw(ttl int, value interface{}) *StringRaw {
	var expireTime time.Time
	if ttl > 0 {
		expireTime = time.Now().Add(time.Duration(ttl) * time.Second)
	}
	return &StringRaw{
		RedisObject: RedisObject{
			objectType: RedisTypeString,
			encoding:   RedisEncodingRaw,
			ttl:        ttl,
			value:      value,
			expireTime: expireTime,
		},
	}
}

func (sr *StringRaw) String() string {
	return sr.value.(string)
}

func (sr *StringRaw) Append(val string) int {
	newValue := sr.GetValue().(string) + val
	sr.SetValue(newValue)
	return len(newValue)
}

func (sr *StringRaw) Incr() (int64, error) {
	return sr.IncrBy("1")
}

func (sr *StringRaw) Decr() (int64, error) {
	return sr.DecrBy("1")
}

func (sr *StringRaw) IncrBy(val string) (int64, error) {
	incrValInt, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, re.ErrNotIntegerOrOutOfRange
	}
	value := sr.GetValue().(string)
	if valueInt, err := strconv.ParseInt(value, 10, 64); err != nil {
		return 0, re.ErrNotIntegerOrOutOfRange
	} else {
		sr.SetValue(strconv.FormatInt(valueInt+incrValInt, 10))
		return valueInt + incrValInt, nil
	}
}

func (sr *StringRaw) DecrBy(val string) (int64, error) {
	valueInt, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, re.ErrNotIntegerOrOutOfRange
	}
	return sr.IncrBy(strconv.FormatInt(-1*valueInt, 10))
}

func (sr *StringRaw) IncrByFloat(val string) (string, error) {
	valueFloat, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return "", re.ErrValueIsNotFloat
	}
	value := sr.GetValue().(string)
	if originValFloat, err := strconv.ParseFloat(value, 64); err != nil {
		return "", re.ErrValueIsNotFloat
	} else {
		ret, err := util.FormatFloatString(fmt.Sprintf("%f", valueFloat+originValFloat))
		if err != nil {
			return "", err
		}
		sr.SetValue(ret)
		return ret, nil
	}
}

func (sr *StringRaw) Strlen() int {
	return len(sr.GetValue().(string))
}

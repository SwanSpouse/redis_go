package encodings

import (
	"fmt"
	re "redis_go/error"
	"redis_go/loggers"
	"redis_go/raw_type"
	"time"
)

type HashDict struct {
	RedisObject
}

func NewHashDict(ttl int) *HashDict {
	var expireTime time.Time
	if ttl > 0 {
		expireTime = time.Now().Add(time.Duration(ttl) * time.Second)
	}
	return &HashDict{
		RedisObject: RedisObject{
			objectType: RedisTypeHash,
			encoding:   RedisEncodingHT,
			ttl:        ttl,
			value:      raw_type.NewDict(),
			expireTime: expireTime,
		},
	}
}

func (hd *HashDict) HSet(key string, value string) int {
	dict := hd.GetValue().(*raw_type.Dict)
	if dict.Put(key, value) == nil {
		return 1
	}
	return 0
}

func (hd *HashDict) HGet(key string) (string, error) {
	dict := hd.GetValue().(*raw_type.Dict)
	if ret := dict.Get(key); ret == nil {
		return "", re.ErrNilValue
	} else {
		return ret.(string), nil
	}
}

func (hd *HashDict) HExists(key string) int {
	dict := hd.GetValue().(*raw_type.Dict)
	if dict.ContainsKey(key) {
		return 1
	}
	return 0
}

func (hd *HashDict) HDel(keys []string) int {
	dict := hd.GetValue().(*raw_type.Dict)
	succCount := 0
	for _, key := range keys {
		if dict.RemoveKey(key) != nil {
			succCount++
		}
	}
	return succCount
}

func (hd *HashDict) HLen() int {
	dict := hd.GetValue().(*raw_type.Dict)
	return dict.Size()
}

func (hd *HashDict) HGetAll() []string {
	dict := hd.GetValue().(*raw_type.Dict)
	keyValues := dict.KeyValueSet()
	ret := make([]string, 0)
	for key, value := range keyValues {
		ret = append(ret, key.(string))
		ret = append(ret, value.(string))
	}
	return ret
}

func (hd *HashDict) HDebug() {
	loggers.Debug(hd.String())
}

func (hd *HashDict) String() string {
	dict := hd.GetValue().(*raw_type.Dict)
	keyValues := dict.KeyValueSet()
	msg := "current dict is"
	for key, value := range keyValues {
		msg += fmt.Sprintf("[%s=%s]", key, value)
	}
	return msg
}

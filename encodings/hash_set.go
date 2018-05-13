package encodings

import (
	"fmt"
	re "redis_go/error"
	"redis_go/loggers"
	"redis_go/raw_type"
	"time"
)

const (
	RedisHashSetDefaultValueInDict = true
)

type HashSet struct {
	RedisObject
}

func NewHasSet(ttl int) *HashSet {
	var expireTime time.Time
	if ttl > 0 {
		expireTime = time.Now().Add(time.Duration(ttl) * time.Second)
	}
	return &HashSet{
		RedisObject: RedisObject{
			objectType: RedisTypeSet,
			encoding:   RedisEncodingHT,
			ttl:        ttl,
			value:      raw_type.NewDict(),
			expireTime: expireTime,
		},
	}
}

func (hs *HashSet) SAdd(values []string) int {
	set := hs.GetValue().(*raw_type.Dict)
	succCount := 0
	for _, value := range values {
		if set.Put(value, RedisHashSetDefaultValueInDict) == nil {
			succCount += 1
		}
	}
	return succCount
}

func (hs *HashSet) SCard() int {
	set := hs.GetValue().(*raw_type.Dict)
	return set.Size()
}

func (hs *HashSet) SIsMember(value string) int {
	set := hs.GetValue().(*raw_type.Dict)
	if set.ContainsKey(value) {
		return 1
	}
	return 0
}

func (hs *HashSet) SMembers() []string {
	set := hs.GetValue().(*raw_type.Dict)
	ret := make([]string, 0)
	for item := range set.KeySet() {
		ret = append(ret, item.(string))
	}
	return ret
}

func (hs *HashSet) SPop() (string, error) {
	set := hs.GetValue().(*raw_type.Dict)
	if key := set.RandomKey(); key == nil {
		return "", re.ErrNilValue
	} else {
		set.RemoveKey(key)
		return key.(string), nil
	}
}

func (hs *HashSet) SRandMember() (string, error) {
	set := hs.GetValue().(*raw_type.Dict)
	if key := set.RandomKey(); key == nil {
		return "", re.ErrNilValue
	} else {
		return key.(string), nil
	}

}

func (hs *HashSet) SRem(values []string) int {
	set := hs.GetValue().(*raw_type.Dict)
	succCount := 0
	for _, value := range values {
		if set.RemoveKey(value) != nil {
			succCount += 1
		}
	}
	return succCount
}

func (hs *HashSet) SDebug() {
	set := hs.GetValue().(*raw_type.Dict)
	msg := ""
	for key := range set.KeySet() {
		msg += fmt.Sprintf("[%s]=>", key.(string))
	}
	loggers.Info("Set Debug Info: %s", msg)
}

func (hs *HashSet) String() string {
	set := hs.GetValue().(*raw_type.Dict)
	keys := set.KeySet()
	msg := "current set is"
	for key := range keys {
		msg += fmt.Sprintf("[%s]", key)
	}
	return msg
}

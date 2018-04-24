package encodings

import "time"

type StringEmb struct {
	RedisObject
}

func NewStringEmb(ttl int, value interface{}) *StringEmb {
	var expireTime time.Time
	if ttl > 0 {
		expireTime = time.Now().Add(time.Duration(ttl) * time.Second)
	}
	return &StringEmb{
		RedisObject: RedisObject{
			objectType: RedisTypeString,
			encoding:   RedisEncodingEmbStr,
			ttl:        ttl,
			value:      value,
			expireTime: expireTime,
		},
	}
}

func (se *StringEmb) String() string {
	return ""
}

func (se *StringEmb) Append(val string) int {
	return 0
}

func (se *StringEmb) Incr() (int64, error) {
	return 0, nil
}

func (se *StringEmb) Decr() (int64, error) {
	return 0, nil
}

func (se *StringEmb) IncrBy(val string) (int64, error) {
	return 0, nil
}

func (se *StringEmb) DecrBy(val string) (int64, error) {
	return 0, nil
}

func (se *StringEmb) Strlen() int {
	return 0
}

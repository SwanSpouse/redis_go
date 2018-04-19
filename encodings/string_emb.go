package encodings

import "time"

type StringEmb struct {
	RedisObject
}

func NewRedisStringWithEncodingEmbStr(value string, ttl int) *StringEmb {
	var expireTime time.Time
	if ttl > 0 {
		expireTime = time.Now().Add(time.Duration(ttl) * time.Second)
	}
	se := &StringEmb{
		RedisObject: RedisObject{
			objectType: RedisTypeString,
			encoding:   RedisEncodingRaw,
			ttl:        ttl,
			value:      value,
			expireTime: expireTime,
		},
	}
	return se
}

func (se *StringEmb) String() string {
	return ""
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

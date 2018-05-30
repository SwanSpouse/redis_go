package encodings

import (
	"redis_go/raw_type"
	"time"
)

type SortedSet struct {
	RedisObject
}

func NewSortedSet(ttl int) *SortedSet {
	var expireTime time.Time
	if ttl > 0 {
		expireTime = time.Now().Add(time.Duration(ttl) * time.Second)
	}
	return &SortedSet{
		RedisObject: RedisObject{
			objectType: RedisTypeZSet,
			encoding:   RedisEncodingSkipList,
			ttl:        ttl,
			value:      raw_type.NewSkipList(),
			expireTime: expireTime,
		},
	}
}

func (ss *SortedSet) ZAdd(inputs []string) int {
	return 0
}

func (ss *SortedSet) ZCard() int {
	return 0
}

func (ss *SortedSet) ZCount() int {
	return 0
}

func (ss *SortedSet) ZIncrBy(key string, increment float64) float64 {
	return 0.0
}

func (ss *SortedSet) ZRange(start, end int) []string {
	return nil
}

func (ss *SortedSet) ZRangeByScore(lower, upper float64) []string {
	return nil
}

func (ss *SortedSet) ZRevRange(start, end int) []string {
	return nil
}

func (ss *SortedSet) ZRevRangeByScore(lower, upper float64) []string {
	return nil
}

func (ss *SortedSet) ZRank(key string) int {
	return 0
}

func (ss *SortedSet) ZRevRank(key string) int {
	return 0
}

func (ss *SortedSet) ZRem(inputs []string) int {
	return 0
}

func (ss *SortedSet) ZRemRangeByRyRank(start, end int) int {
	return 0
}

func (ss *SortedSet) ZRemRangeByScore(lower, upper float64) int {
	return 0
}

func (ss *SortedSet) ZScore(key string) float64 {
	return 0.0
}

func (ss *SortedSet) String() string {
	return ""
}

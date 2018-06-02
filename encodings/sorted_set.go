package encodings

import (
	"fmt"
	re "redis_go/error"
	"redis_go/raw_type"
	"redis_go/util"
	"strconv"
	"time"
)

type SortedSet struct {
	RedisObject
	dict map[string]*raw_type.SkipNode
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
		dict: make(map[string]*raw_type.SkipNode),
	}
}

func (ss *SortedSet) ZAdd(inputs []string) (int, error) {
	if len(inputs)%2 != 0 {
		return 0, re.ErrSyntaxError
	}
	for i := 0; i < len(inputs); i += 2 {
		_, err := strconv.ParseFloat(inputs[i], 64)
		if err != nil {
			return 0, re.ErrValueIsNotFloat
		}
	}
	skipList := ss.GetValue().(*raw_type.SkipList)
	count := 0
	for i := 0; i < len(inputs); i += 2 {
		score, _ := strconv.ParseFloat(inputs[i], 64)
		key := inputs[i+1]
		if obj, exists := ss.dict[key]; exists {
			if obj.GetScore() != score {
				obj.SetScore(score)
				count += 1
			}
		} else {
			newNode := skipList.Insert(key, score)
			ss.dict[key] = newNode
			count += 1
		}
	}
	return count, nil
}

func (ss *SortedSet) ZCard() int {
	skipList := ss.GetValue().(*raw_type.SkipList)
	return skipList.Length()
}

func (ss *SortedSet) ZCount(lower, upper string) (int, error) {
	fLower, err1 := strconv.ParseFloat(lower, 64)
	fUpper, err2 := strconv.ParseFloat(upper, 64)
	if err1 != nil || err2 != nil {
		return 0, re.ErrValueIsNotFloat
	}
	skipList := ss.GetValue().(*raw_type.SkipList)

	firstNode := skipList.FirstInRange(raw_type.RangeSpec{
		Min: fLower, Max: fUpper, MinEx: false, MaxEx: false,
	})
	count := 0
	for cur := firstNode; cur != nil; cur.GetNextNode() {
		count += 1
	}
	return count, nil
}

func (ss *SortedSet) ZIncrBy(key string, increment string) (float64, error) {
	fIncrement, err := strconv.ParseFloat(increment, 64)
	if err != nil {
		return 0, re.ErrValueIsNotFloat
	}
	if obj, exists := ss.dict[key]; !exists {
		skipList := ss.GetValue().(*raw_type.SkipList)
		newNode := skipList.Insert(key, fIncrement)
		ss.dict[key] = newNode
		return fIncrement, nil
	} else {
		newScore := obj.GetScore() + fIncrement
		obj.SetScore(newScore)
		return newScore, nil
	}
	return 0.0, nil
}

func (ss *SortedSet) ZRange(lower, upper string) ([]string, error) {
	iLower, err1 := strconv.Atoi(lower)
	iUpper, err2 := strconv.Atoi(upper)
	if err1 != nil || err2 != nil {
		return nil, re.ErrNotIntegerOrOutOfRange
	}
	skipList := ss.GetValue().(*raw_type.SkipList)
	ret := make([]string, 0)

	if iLower > skipList.Length() || -iLower > skipList.Length() {
		return ret, re.ErrEmptyListOrSet
	}
	if (iLower > 0 && iUpper > 0) || (iLower < 0 && iUpper < 0) && iLower > iUpper {
		return ret, re.ErrEmptyListOrSet
	}
	if iLower < 0 {
		iLower = skipList.Length() - 1 + iLower
	}
	if iUpper < 0 {
		iUpper = skipList.Length() - 1 + iUpper
	}
	endNode := skipList.GetElementByRank(iUpper)
	for cur := skipList.GetElementByRank(iLower); cur != endNode; cur.GetNextNode() {
		ret = append(ret, cur.GetValue())
		ret = append(ret, util.FloatToSimpleString(cur.GetScore()))
	}
	if endNode != nil {
		ret = append(ret, endNode.GetValue())
		ret = append(ret, util.FloatToSimpleString(endNode.GetScore()))
	}
	return ret, nil
}

func (ss *SortedSet) ZRangeByScore(lower, upper string) ([]string, error) {
	fLower, err1 := strconv.ParseFloat(lower, 64)
	fUpper, err2 := strconv.ParseFloat(upper, 64)
	if err1 != nil || err2 != nil {
		return nil, re.ErrValueIsNotFloat
	}
	skipList := ss.GetValue().(*raw_type.SkipList)
	ret := make([]string, 0)

	spec := raw_type.RangeSpec{
		Min: fLower, Max: fUpper, MinEx: false, MaxEx: false,
	}
	endNode := skipList.LastInRange(spec)
	for cur := skipList.FirstInRange(spec); cur != endNode; cur.GetNextNode() {
		ret = append(ret, cur.GetValue())
		ret = append(ret, util.FloatToSimpleString(cur.GetScore()))
	}
	if endNode != nil {
		ret = append(ret, endNode.GetValue())
		ret = append(ret, util.FloatToSimpleString(endNode.GetScore()))
	}
	return ret, nil
}

func (ss *SortedSet) ZRevRange(lower, upper string) ([]string, error) {
	if ret, err := ss.ZRange(lower, upper); err != nil {
		return ret, err
	} else {
		revRet := make([]string, len(ret))
		i := 0
		for index := len(ret) - 1; index >= 0; index -= 1 {
			revRet[i] = ret[index]
			i += 1
		}
		return revRet, nil
	}
}

func (ss *SortedSet) ZRevRangeByScore(lower, upper string) ([]string, error) {
	if ret, err := ss.ZRange(lower, upper); err != nil {
		return ret, nil
	} else {
		revRet := make([]string, len(ret))
		i := 0
		for index := len(ret) - 1; index >= 0; index -= 1 {
			revRet[i] = ret[index]
			i += 1
		}
		return revRet, nil
	}
}

func (ss *SortedSet) ZRank(key string) (int, error) {
	if obj, exits := ss.dict[key]; !exits {
		return 0, re.ErrNoSuchKey
	} else {
		skipList := ss.GetValue().(*raw_type.SkipList)
		return skipList.GetRank(key, obj.GetScore()), nil
	}
}

func (ss *SortedSet) ZRevRank(key string) (int, error) {
	if obj, exits := ss.dict[key]; !exits {
		return 0, re.ErrNoSuchKey
	} else {
		skipList := ss.GetValue().(*raw_type.SkipList)
		rank := skipList.GetRank(key, obj.GetScore())
		return skipList.Length() - rank, nil
	}
}

func (ss *SortedSet) ZRem(inputs []string) int {
	skipList := ss.GetValue().(*raw_type.SkipList)
	count := 0
	for _, key := range inputs {
		if obj, exists := ss.dict[key]; exists {
			score := obj.GetScore()
			skipList.Delete(key, score)
			delete(ss.dict, key)
			count += 1
		}
	}
	return count
}

func (ss *SortedSet) ZRemRangeByRank(start, end string) (int, error) {
	iStart, err1 := strconv.Atoi(start)
	iEnd, err2 := strconv.Atoi(end)
	if err1 != nil || err2 != nil {
		return 0, re.ErrNotIntegerOrOutOfRange
	}
	skipList := ss.GetValue().(*raw_type.SkipList)
	return skipList.DeleteRangeByRank(iStart, iEnd), nil
}

func (ss *SortedSet) ZRemRangeByScore(lower, upper string) (int, error) {
	fLower, err1 := strconv.ParseFloat(lower, 64)
	fUpper, err2 := strconv.ParseFloat(upper, 64)
	if err1 != nil || err2 != nil {
		return 0, re.ErrValueIsNotFloat
	}
	skipList := ss.GetValue().(*raw_type.SkipList)
	return skipList.DeleteRangeByScore(raw_type.RangeSpec{
		Min: fLower, Max: fUpper, MinEx: false, MaxEx: false,
	}), nil
}

func (ss *SortedSet) ZScore(key string) (float64, error) {
	if obj, exists := ss.dict[key]; !exists {
		return 0.0, re.ErrNoSuchKey
	} else {
		return obj.GetScore(), nil
	}
	return 0.0, nil
}

func (ss *SortedSet) String() string {
	return fmt.Sprintf("SortedSet:%+v", ss.dict)
}

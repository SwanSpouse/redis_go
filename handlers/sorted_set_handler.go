package handlers

import (
	"github.com/SwanSpouse/redis_go/client"
	"github.com/SwanSpouse/redis_go/database"
	"github.com/SwanSpouse/redis_go/encodings"
	re "github.com/SwanSpouse/redis_go/error"
	"github.com/SwanSpouse/redis_go/loggers"
	"github.com/SwanSpouse/redis_go/util"
)

const (
	RedisSortedSetCommandZAdd             = "ZADD"
	RedisSortedSetCommandZCard            = "ZCARD"
	RedisSortedSetCommandZCount           = "ZCOUNT"
	RedisSortedSetCommandZIncrBy          = "ZINCRBY"
	RedisSortedSetCommandZRange           = "ZRANGE"
	RedisSortedSetCommandZRangeByScore    = "ZRANGEBYSCORE"
	RedisSortedSetCommandZRank            = "ZRANK"
	RedisSortedSetCommandZRem             = "ZREM"
	RedisSortedSetCommandZRemRangeByRank  = "ZREMRANGEBYRANK"
	RedisSortedSetCommandZRemRangeByScore = "ZREMRANGEBYSCORE"
	RedisSortedSetCommandZRevRange        = "ZREVRANGE"
	RedisSortedSetCommandZRevRangeByScore = "ZREVRANGEBYSCORE"
	RedisSortedSetCommandZRevRank         = "ZREVRANK"
	RedisSortedSetCommandZScore           = "ZSCORE"
	RedisSortedSetCommandZUnionStore      = "ZUNIONSTORE"
	RedisSortedSetCommandZInterStore      = "ZINTERSTORE"
	RedisSortedSetCommandZScan            = "ZSCAN"
)

// SetHandler可以处理的一种rawType
var setEncodingTypeSortedSet = map[string]bool{
	encodings.RedisEncodingSkipList: true,
}

type SortedSetHandler struct {
}

func getTZSetValueByKey(cli *client.Client, key string) (database.TZSet, error) {
	baseType := cli.SelectedDatabase().SearchKeyInDB(key)
	if baseType == nil {
		return nil, re.ErrNoSuchKey
	}
	if _, ok := setEncodingTypeSortedSet[baseType.GetEncoding()]; !ok || baseType.GetObjectType() != encodings.RedisTypeZSet {
		loggers.Errorf(string(re.ErrWrongType), baseType.GetObjectType(), baseType.GetEncoding())
		return nil, re.ErrWrongType
	}
	if tzs, ok := baseType.(database.TZSet); ok {
		return tzs, nil
	}
	loggers.Errorf("base type:%+v can not convert to TSet", baseType)
	return nil, re.ErrWrongType
}

func createZSetIfNotExists(cli *client.Client, key string) error {
	if _, err := getTZSetValueByKey(cli, key); err != nil && err != re.ErrNoSuchKey {
		return err
	} else if err == re.ErrNoSuchKey {
		cli.SelectedDatabase().SetKeyInDB(key, database.NewRedisSortedSetObject())
	}
	return nil
}

func (handler *SortedSetHandler) ZAdd(cli *client.Client) {
	key := cli.Argv[1]
	if err := createZSetIfNotExists(cli, key); err != nil {
		cli.ResponseReError(err)
		return
	}
	if tss, err := getTZSetValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		count, err := tss.ZAdd(cli.Argv[2:])
		if err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(count)
			if count != 0 {
				cli.Dirty += 1
			}
		}
	}
}

func (handler *SortedSetHandler) ZCard(cli *client.Client) {
	key := cli.Argv[1]
	if tss, err := getTZSetValueByKey(cli, key); err != nil && err == re.ErrNoSuchKey {
		cli.Response(0)
	} else if err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(tss.ZCard())
	}
}

func (handler *SortedSetHandler) ZCount(cli *client.Client) {
	key := cli.Argv[1]
	if tss, err := getTZSetValueByKey(cli, key); err != nil && err == re.ErrNoSuchKey {
		cli.Response(0)
	} else if err != nil {
		cli.ResponseReError(err)
	} else {
		count, err := tss.ZCount(cli.Argv[2], cli.Argv[3])
		if err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(count)
		}
	}
}

func (handler *SortedSetHandler) ZIncrBy(cli *client.Client) {
	key := cli.Argv[1]
	if err := createZSetIfNotExists(cli, key); err != nil {
		cli.ResponseReError(err)
		return
	}
	if tss, err := getTZSetValueByKey(cli, key); err != nil && err == re.ErrNoSuchKey {
		cli.Response(0)
	} else if err != nil {
		cli.ResponseReError(err)
	} else {
		ret, err := tss.ZIncrBy(cli.Argv[2], cli.Argv[3])
		if err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(ret)
			cli.Dirty += 1
		}
	}
}

func (handler *SortedSetHandler) ZRange(cli *client.Client) {
	key := cli.Argv[1]
	if tss, err := getTZSetValueByKey(cli, key); err != nil && err == re.ErrNoSuchKey {
		cli.ResponseReError(re.ErrEmptyListOrSet)
	} else if err != nil {
		cli.ResponseReError(err)
	} else {
		if ret, err := tss.ZRange(cli.Argv[2], cli.Argv[3]); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(ret)
		}
	}
}

func (handler *SortedSetHandler) ZRangeByScore(cli *client.Client) {
	key := cli.Argv[1]
	if tss, err := getTZSetValueByKey(cli, key); err != nil && err == re.ErrNoSuchKey {
		cli.ResponseReError(re.ErrEmptyListOrSet)
	} else if err != nil {
		cli.ResponseReError(err)
	} else {
		if ret, err := tss.ZRangeByScore(cli.Argv[2], cli.Argv[3]); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(ret)
		}
	}
}

func (handler *SortedSetHandler) ZRank(cli *client.Client) {
	key := cli.Argv[1]
	if tss, err := getTZSetValueByKey(cli, key); err != nil && err == re.ErrNoSuchKey {
		cli.Response(nil)
	} else if err != nil {
		cli.ResponseReError(err)
	} else {
		if rank, err := tss.ZRank(cli.Argv[2]); err != nil && err == re.ErrNoSuchKey {
			cli.Response(nil)
		} else {
			cli.Response(rank)
		}
	}
}

func (handler *SortedSetHandler) ZRem(cli *client.Client) {
	key := cli.Argv[1]
	if tss, err := getTZSetValueByKey(cli, key); err != nil && err == re.ErrNoSuchKey {
		cli.Response(0)
	} else if err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(tss.ZRem(cli.Argv[2:]))
	}
}

func (handler *SortedSetHandler) ZRemRangeByRank(cli *client.Client) {
	key := cli.Argv[1]
	if tss, err := getTZSetValueByKey(cli, key); err != nil && err == re.ErrNoSuchKey {
		cli.ResponseReError(re.ErrEmptyListOrSet)
	} else if err != nil {
		cli.ResponseReError(err)
	} else {
		if ret, err := tss.ZRemRangeByRank(cli.Argv[2], cli.Argv[3]); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(ret)
			cli.Dirty += int64(ret)
		}
	}
}

func (handler *SortedSetHandler) ZRemRangeByScore(cli *client.Client) {
	key := cli.Argv[1]
	if tss, err := getTZSetValueByKey(cli, key); err != nil && err == re.ErrNoSuchKey {
		cli.ResponseReError(re.ErrEmptyListOrSet)
	} else if err != nil {
		cli.ResponseReError(err)
	} else {
		if ret, err := tss.ZRemRangeByScore(cli.Argv[2], cli.Argv[3]); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(ret)
			cli.Dirty += int64(ret)
		}
	}
}

func (handler *SortedSetHandler) ZRevRange(cli *client.Client) {
	key := cli.Argv[1]
	if tss, err := getTZSetValueByKey(cli, key); err != nil && err == re.ErrNoSuchKey {
		cli.ResponseReError(re.ErrEmptyListOrSet)
	} else if err != nil {
		cli.ResponseReError(err)
	} else {
		if ret, err := tss.ZRevRange(cli.Argv[2], cli.Argv[3]); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(ret)
		}
	}
}

func (handler *SortedSetHandler) ZRevRangeByScore(cli *client.Client) {
	key := cli.Argv[1]
	if tss, err := getTZSetValueByKey(cli, key); err != nil && err == re.ErrNoSuchKey {
		cli.ResponseReError(re.ErrEmptyListOrSet)
	} else if err != nil {
		cli.ResponseReError(err)
	} else {
		if ret, err := tss.ZRevRangeByScore(cli.Argv[2], cli.Argv[3]); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(ret)
		}
	}
}

func (handler *SortedSetHandler) ZRevRank(cli *client.Client) {
	key := cli.Argv[1]
	if tss, err := getTZSetValueByKey(cli, key); err != nil && err == re.ErrNoSuchKey {
		cli.ResponseReError(re.ErrEmptyListOrSet)
	} else if err != nil {
		cli.ResponseReError(err)
	} else {
		if ret, err := tss.ZRevRank(cli.Argv[2]); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(ret)
		}
	}
}

func (handler *SortedSetHandler) ZScore(cli *client.Client) {
	key := cli.Argv[1]
	if tss, err := getTZSetValueByKey(cli, key); err != nil && err == re.ErrNoSuchKey {
		cli.Response(nil)
	} else if err != nil {
		cli.ResponseReError(err)
	} else {
		if score, err := tss.ZScore(cli.Argv[2]); err == re.ErrNoSuchKey {
			cli.Response(nil)
		} else {
			cli.Response(util.FloatToSimpleString(score))
		}
	}
}

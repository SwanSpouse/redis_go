package handlers

import (
	"redis_go/client"
	"redis_go/database"
	"redis_go/encodings"
	re "redis_go/error"
	"redis_go/loggers"
	"redis_go/util"
)

var (
	_ client.BaseHandler = (*SortedSetHandler)(nil)
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

func (handler *SortedSetHandler) Process(cli *client.Client) {
	switch cli.Cmd.GetName() {
	case RedisSortedSetCommandZAdd:
		handler.ZAdd(cli)
	case RedisSortedSetCommandZCard:
		handler.ZCard(cli)
	case RedisSortedSetCommandZCount:
		handler.ZCount(cli)
	case RedisSortedSetCommandZIncrBy:
		handler.ZIncrBy(cli)
	case RedisSortedSetCommandZRange:
		handler.ZRange(cli)
	case RedisSortedSetCommandZRangeByScore:
		handler.ZRangeByScore(cli)
	case RedisSortedSetCommandZRank:
		handler.ZRank(cli)
	case RedisSortedSetCommandZRem:
		handler.ZRem(cli)
	case RedisSortedSetCommandZRemRangeByRank:
		handler.ZRemRangeByRank(cli)
	case RedisSortedSetCommandZRemRangeByScore:
		handler.ZRemRangeByScore(cli)
	case RedisSortedSetCommandZRevRange:
		handler.ZRevRange(cli)
	case RedisSortedSetCommandZRevRangeByScore:
		handler.ZRevRangeByScore(cli)
	case RedisSortedSetCommandZRevRank:
		handler.ZRevRank(cli)
	case RedisSortedSetCommandZScore:
		handler.ZScore(cli)
	case RedisSortedSetCommandZUnionStore:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisSortedSetCommandZInterStore:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisSortedSetCommandZScan:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	default:
		cli.ResponseReError(re.ErrUnknownCommand, cli.Cmd.GetOriginName())
	}
	cli.Flush()
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
		ret := tss.ZAdd(cli.Argv[2:])
		cli.Response(ret)
		if ret != 0 {
			cli.Dirty += 1
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

}

func (handler *SortedSetHandler) ZIncrBy(cli *client.Client) {

}

func (handler *SortedSetHandler) ZRange(cli *client.Client) {

}

func (handler *SortedSetHandler) ZRangeByScore(cli *client.Client) {

}

func (handler *SortedSetHandler) ZRank(cli *client.Client) {

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

}

func (handler *SortedSetHandler) ZRemRangeByScore(cli *client.Client) {

}

func (handler *SortedSetHandler) ZRevRange(cli *client.Client) {

}

func (handler *SortedSetHandler) ZRevRangeByScore(cli *client.Client) {

}

func (handler *SortedSetHandler) ZRevRank(cli *client.Client) {

}

func (handler *SortedSetHandler) ZScore(cli *client.Client) {
	key := cli.Argv[1]
	if tss, err := getTZSetValueByKey(cli, key); err != nil && err == re.ErrNoSuchKey {
		cli.Response(nil)
	} else if err != nil {
		cli.ResponseReError(err)
	} else {
		score := tss.ZScore(cli.Argv[2])
		cli.Response(util.FloatToSimpleString(score))
	}
}

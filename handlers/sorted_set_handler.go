package handlers

import (
	"redis_go/client"
	"redis_go/encodings"
	re "redis_go/error"
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
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisSortedSetCommandZCard:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisSortedSetCommandZCount:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisSortedSetCommandZIncrBy:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisSortedSetCommandZRange:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisSortedSetCommandZRangeByScore:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisSortedSetCommandZRank:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisSortedSetCommandZRem:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisSortedSetCommandZRemRangeByRank:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisSortedSetCommandZRemRangeByScore:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisSortedSetCommandZRevRange:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisSortedSetCommandZRevRangeByScore:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisSortedSetCommandZRevRank:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisSortedSetCommandZScore:
		cli.ResponseReError(re.ErrFunctionNotImplement)
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

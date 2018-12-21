package handlers

import (
	"strconv"
	"strings"

	"github.com/SwanSpouse/redis_go/client"
	"github.com/SwanSpouse/redis_go/database"
	"github.com/SwanSpouse/redis_go/encodings"
	re "github.com/SwanSpouse/redis_go/error"
	"github.com/SwanSpouse/redis_go/loggers"
)

const (
	RedisListCommandLIndex    = "LINDEX"
	RedisListCommandLInsert   = "LINSERT"
	RedisListCommandLLen      = "LLEN"
	RedisListCommandLPop      = "LPOP"
	RedisListCommandLPush     = "LPUSH"
	RedisListCommandLPushX    = "LPUSHX"
	RedisListCommandLRange    = "LRANGE"
	RedisListCommandLRem      = "LREM"
	RedisListCommandLSet      = "LSET"
	RedisListCommandLTrim     = "LTRIM"
	RedisListCommandRPop      = "RPOP"
	RedisListCommandRPopLPush = "RPOPLPUSH"
	RedisListCommandRPush     = "RPUSH"
	RedisListCommandRpushX    = "RPUSHX"
	RedisListCommandLDebug    = "LDEBUG"
)

// ListHandler可以处理的两种rawType
var listEncodingTypeDict = map[string]bool{
	encodings.RedisEncodingLinkedList: true,
	encodings.RedisEncodingZipList:    true,
}

type ListHandler struct {
}

func getTListValueByKey(cli *client.Client, key string) (database.TList, error) {
	baseType := cli.SelectedDatabase().SearchKeyInDB(key)
	if baseType == nil {
		return nil, re.ErrNoSuchKey
	}
	if _, ok := listEncodingTypeDict[baseType.GetEncoding()]; !ok || baseType.GetObjectType() != encodings.RedisTypeList {
		loggers.Errorf(string(re.ErrWrongType), baseType.GetObjectType(), baseType.GetEncoding())
		return nil, re.ErrWrongType
	}
	if ts, ok := baseType.(database.TList); ok {
		return ts, nil
	}
	loggers.Errorf("base type can not convert to TList")
	return nil, re.ErrWrongType
}

func createListIfNotExists(cli *client.Client, key string) error {
	if _, err := getTListValueByKey(cli, key); err != nil && err != re.ErrNoSuchKey {
		return err
	} else if err == re.ErrNoSuchKey {
		cli.SelectedDatabase().SetKeyInDB(key, database.NewRedisListObject())
	}
	return nil
}

func (handler *ListHandler) LIndex(cli *client.Client) {
	key := cli.Argv[1]
	if ts, err := getTListValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		index, err := strconv.Atoi(cli.Argv[2])
		if err != nil {
			cli.ResponseReError(re.ErrNotIntegerOrOutOfRange)
			return
		}
		if ret, err := ts.LIndex(index); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(ret)
		}
	}
}

func (handler *ListHandler) LInsert(cli *client.Client) {
	key := cli.Argv[1]
	if ts, err := getTListValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
		cli.Response(nil)
	} else {
		var insertFlag int
		switch strings.ToUpper(cli.Argv[2]) {
		case "BEFORE":
			insertFlag = encodings.RedisTypeListInsertBefore
		case "AFTER":
			insertFlag = encodings.RedisTypeListInsertAfter
		default:
			cli.ResponseReError(re.ErrSyntaxError)
			return
		}
		if ret, err := ts.LInsert(insertFlag, cli.Argv[3], cli.Argv[4:]...); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(ret)
			cli.Dirty += 1
		}
	}
}

func (handler *ListHandler) LLen(cli *client.Client) {
	key := cli.Argv[1]
	if ts, err := getTListValueByKey(cli, key); err != nil && err != re.ErrNilValue {
		cli.ResponseReError(err)
	} else if err == re.ErrNilValue {
		cli.Response(0)
	} else {
		cli.Response(ts.LLen())
	}
}

func (handler *ListHandler) LPop(cli *client.Client) {
	key := cli.Argv[1]
	if tl, err := getTListValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(tl.LPop())
		cli.Dirty += 1
	}
}

func (handler *ListHandler) LPush(cli *client.Client) {
	key := cli.Argv[1]
	if err := createListIfNotExists(cli, key); err != nil {
		cli.ResponseReError(err)
		return
	}
	if tl, err := getTListValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(tl.LPush(cli.Argv[2:]))
		cli.Dirty += 1
	}
}

func (handler *ListHandler) LRange(cli *client.Client) {
	key := cli.Argv[1]
	if ts, err := getTListValueByKey(cli, key); err != nil && err != re.ErrNoSuchKey {
		cli.ResponseReError(err)
	} else if err == re.ErrNoSuchKey {
		cli.ResponseReError(re.ErrEmptyListOrSet)
	} else {
		start, startErr := strconv.Atoi(cli.Argv[2])
		stop, stopErr := strconv.Atoi(cli.Argv[3])
		if startErr != nil || stopErr != nil {
			cli.ResponseReError(re.ErrNotIntegerOrOutOfRange)
			return
		}
		cli.Response(ts.LRange(start, stop))
	}
}

func (handler *ListHandler) LRem(cli *client.Client) {
	key := cli.Argv[1]
	if ts, err := getTListValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		index, err := strconv.Atoi(cli.Argv[2])
		if err != nil {
			cli.ResponseReError(re.ErrNotIntegerOrOutOfRange)
			return
		}
		cli.Response(ts.LRem(index, cli.Argv[3]))
		cli.Dirty += 1
	}
}

func (handler *ListHandler) LSet(cli *client.Client) {
	key := cli.Argv[1]
	if ts, err := getTListValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		index, err := strconv.Atoi(cli.Argv[2])
		if err != nil {
			cli.ResponseReError(re.ErrNotIntegerOrOutOfRange)
			return
		}
		if err := ts.LSet(index, cli.Argv[3]); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.ResponseOK()
			cli.Dirty += 1
		}
	}
}

func (handler *ListHandler) LTrim(cli *client.Client) {
	key := cli.Argv[1]
	if ts, err := getTListValueByKey(cli, key); err != nil && err != re.ErrNoSuchKey {
		cli.ResponseReError(err)
	} else if err == re.ErrNoSuchKey {
		cli.ResponseOK()
	} else {
		startPos, startPosErr := strconv.Atoi(cli.Argv[2])
		endPos, endPosErr := strconv.Atoi(cli.Argv[3])
		if startPosErr != nil || endPosErr != nil {
			cli.ResponseReError(re.ErrNotIntegerOrOutOfRange)
			return
		}
		if err := ts.LTrim(startPos, endPos); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.ResponseOK()
			cli.Dirty += 1
		}
	}
}

func (handler *ListHandler) RPop(cli *client.Client) {
	key := cli.Argv[1]
	if tl, err := getTListValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(tl.RPop())
		cli.Dirty += 1
	}
}

func (handler *ListHandler) RPush(cli *client.Client) {
	key := cli.Argv[1]
	if err := createListIfNotExists(cli, key); err != nil {
		cli.ResponseReError(err)
		return
	}
	if tl, err := getTListValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(tl.RPush(cli.Argv[2:]))
		cli.Dirty += 1
	}
}

func (handler *ListHandler) Debug(cli *client.Client) {
	key := cli.Argv[1]
	if tl, err := getTListValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		tl.Debug()
		cli.ResponseOK()
	}
}

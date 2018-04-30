package handlers

import (
	"redis_go/client"
	"redis_go/database"
	"redis_go/encodings"
	re "redis_go/error"
	"redis_go/loggers"
	"strconv"
	"strings"
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

func (handler *ListHandler) Process(cli *client.Client) {
	switch cli.Cmd.GetName() {
	case RedisListCommandLIndex:
		handler.LIndex(cli)
	case RedisListCommandLInsert:
		handler.LInsert(cli)
	case RedisListCommandLLen:
		handler.LLen(cli)
	case RedisListCommandLPop:
		handler.LPop(cli)
	case RedisListCommandLPush:
		handler.LPush(cli)
	case RedisListCommandLPushX:
		handler.LPushX(cli)
	case RedisListCommandLRange:
		handler.LRange(cli)
	case RedisListCommandLRem:
		handler.LRem(cli)
	case RedisListCommandLSet:
		handler.LSet(cli)
	case RedisListCommandLTrim:
		handler.LTrim(cli)
	case RedisListCommandRPop:
		handler.RPop(cli)
	case RedisListCommandRPopLPush:
		handler.RPopLPush(cli)
	case RedisListCommandRPush:
		handler.RPush(cli)
	case RedisListCommandRpushX:
		handler.RPushX(cli)
	case RedisListCommandLDebug:
		handler.Debug(cli)
	default:
		cli.ResponseReError(re.ErrUnknownCommand, cli.Cmd.GetOriginName())
	}
	cli.Flush()
}

func getTListValueByKey(cli *client.Client, key string) (database.TList, error) {
	if cli.Cmd == nil {
		return nil, re.ErrNilCommand
	}
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
	args := cli.Cmd.GetArgs()
	if len(args) < 2 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if ts, err := getTListValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		index, err := strconv.Atoi(args[1])
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
	args := cli.Cmd.GetArgs()
	if len(args) < 4 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if ts, err := getTListValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
		cli.Response(nil)
	} else {
		var insertFlag int
		switch strings.ToUpper(args[1]) {
		case "BEFORE":
			insertFlag = encodings.RedisTypeListInsertBefore
		case "AFTER":
			insertFlag = encodings.RedisTypeListInsertAfter
		default:
			cli.ResponseReError(re.ErrSyntaxError)
			return
		}
		if ret, err := ts.LInsert(insertFlag, args[2], args[3:]...); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(ret)
		}
	}
}

func (handler *ListHandler) LLen(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) != 1 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if ts, err := getTListValueByKey(cli, key); err != nil && err != re.ErrNilValue {
		cli.ResponseReError(err)
	} else if err == re.ErrNilValue {
		cli.Response(0)
	} else {
		cli.Response(ts.LLen())
	}
}

func (handler *ListHandler) LPop(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) < 1 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if tl, err := getTListValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(tl.LPop())
	}
}

func (handler *ListHandler) LPush(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) < 2 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if err := createListIfNotExists(cli, key); err != nil {
		cli.ResponseReError(err)
		return
	}
	if tl, err := getTListValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(tl.LPush(args[1:]))
	}
}

func (handler *ListHandler) LPushX(cli *client.Client) {

}

func (handler *ListHandler) LRange(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) < 3 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if ts, err := getTListValueByKey(cli, key); err != nil && err != re.ErrNoSuchKey {
		cli.ResponseReError(err)
	} else if err == re.ErrNoSuchKey {
		cli.ResponseReError(re.ErrEmptyListOrSet)
	} else {
		start, startErr := strconv.Atoi(args[1])
		stop, stopErr := strconv.Atoi(args[2])
		if startErr != nil || stopErr != nil {
			cli.ResponseReError(re.ErrNotIntegerOrOutOfRange)
			return
		}
		cli.Response(ts.LRange(start, stop))
	}
}

func (handler *ListHandler) LRem(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) < 3 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if ts, err := getTListValueByKey(cli, key); err != nil && err != re.ErrNoSuchKey {
		cli.ResponseReError(err)
	} else {
		index, err := strconv.Atoi(args[1])
		if err != nil {
			cli.ResponseReError(re.ErrNotIntegerOrOutOfRange)
			return
		}
		cli.Response(ts.LRem(index, args[2]))
	}
}

func (handler *ListHandler) LSet(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) < 3 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if ts, err := getTListValueByKey(cli, key); err != nil && err != re.ErrNoSuchKey {
		cli.ResponseReError(err)
	} else {
		index, err := strconv.Atoi(args[1])
		if err != nil {
			cli.ResponseReError(re.ErrNotIntegerOrOutOfRange)
			return
		}
		if err := ts.LSet(index, args[2]); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.ResponseOK()
		}
	}
}

func (handler *ListHandler) LTrim(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) < 3 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if ts, err := getTListValueByKey(cli, key); err != nil && err != re.ErrNoSuchKey {
		cli.ResponseReError(err)
	} else if err == re.ErrNoSuchKey {
		cli.ResponseOK()
	} else {
		startPos, startPosErr := strconv.Atoi(args[1])
		endPos, endPosErr := strconv.Atoi(args[2])
		if startPosErr != nil || endPosErr != nil {
			cli.ResponseReError(re.ErrNotIntegerOrOutOfRange)
			return
		}
		if err := ts.LTrim(startPos, endPos); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.ResponseOK()
		}
	}
}

func (handler *ListHandler) RPop(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) < 1 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if tl, err := getTListValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(tl.RPop())
	}
}

func (handler *ListHandler) RPopLPush(cli *client.Client) {

}

func (handler *ListHandler) RPush(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) < 2 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if err := createListIfNotExists(cli, key); err != nil {
		cli.ResponseReError(err)
		return
	}
	if tl, err := getTListValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(tl.RPush(args[1:]))
	}
}

func (handler *ListHandler) RPushX(cli *client.Client) {

}

func (handler *ListHandler) Debug(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) < 1 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if tl, err := getTListValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		tl.Debug()
		cli.ResponseOK()
	}
}

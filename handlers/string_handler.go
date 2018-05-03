package handlers

import (
	"fmt"
	"redis_go/client"
	"redis_go/database"
	"redis_go/encodings"
	re "redis_go/error"
	"redis_go/loggers"
	"strconv"
)

var (
	_ client.BaseHandler = (*StringHandler)(nil)
)

const (
	RedisStringCommandAppend      = "APPEND"
	RedisStringCommandBitCount    = "BITCOUNT"
	RedisStringCommandBitop       = "BITOP"
	RedisStringCommandGetBit      = "GETBIT"
	RedisStringCommandSetBit      = "SETBIT"
	RedisStringCommandDecr        = "DECR"
	RedisStringCommandDecrBy      = "DECRBY"
	RedisStringCommandGet         = "GET"
	RedisStringCommandGetRange    = "GETRANGE"
	RedisStringCommandGetSet      = "GETSET"
	RedisStringCommandIncr        = "INCR"
	RedisStringCommandIncrBy      = "INCRBY"
	RedisStringCommandIncrByFloat = "INCRBYFLOAT"
	RedisStringCommandMGet        = "MGET"
	RedisStringCommandMSet        = "MSET"
	RedisStringCommandMSetNx      = "MSETNX"
	RedisStringCommandPSetEx      = "PSETEX"
	RedisStringCommandSet         = "SET"
	RedisStringCommandSetNx       = "SETNX"
	RedisStringCommandSetEX       = "SETEX"
	RedisStringCommandSetRange    = "SETRANGE"
	RedisStringCommandStrLen      = "STRLEN"
)

// StringHandler可以处理的三种rawType
var stringEncodingTypeDict = map[string]bool{
	encodings.RedisEncodingInt:    true,
	encodings.RedisEncodingRaw:    true,
	encodings.RedisEncodingEmbStr: true,
}

type StringHandler struct{}

func (handler *StringHandler) Process(cli *client.Client) {
	switch cli.GetCommandName() {
	case RedisStringCommandAppend:
		handler.Append(cli)
	case RedisStringCommandBitCount, RedisStringCommandBitop, RedisStringCommandGetBit, RedisStringCommandSetBit:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisStringCommandDecr:
		handler.Decr(cli)
	case RedisStringCommandDecrBy:
		handler.DecrBy(cli)
	case RedisStringCommandGet:
		handler.Get(cli)
	case RedisStringCommandGetRange:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisStringCommandGetSet:
		handler.GetSet(cli)
	case RedisStringCommandIncr:
		handler.Incr(cli)
	case RedisStringCommandIncrBy:
		handler.IncrBy(cli)
	case RedisStringCommandIncrByFloat:
		handler.IncrByFloat(cli)
	case RedisStringCommandMGet:
		handler.MGet(cli)
	case RedisStringCommandMSet:
		handler.MSet(cli)
	case RedisStringCommandMSetNx:
		handler.MSetNx(cli)
	case RedisStringCommandPSetEx:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisStringCommandSet:
		handler.Set(cli)
	case RedisStringCommandSetNx:
		handler.SetNx(cli)
	case RedisStringCommandSetEX:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisStringCommandSetRange:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisStringCommandStrLen:
		handler.Strlen(cli)
	default:
		cli.ResponseReError(re.ErrUnknownCommand, cli.GetOriginCommandName())
	}
	// 最后统一发送数据
	cli.Flush()
}

func getTStringValueByKey(cli *client.Client, key string) (database.TString, error) {
	// 获取key在数据库中对应的value(TBase:BaseType)
	baseType := cli.SelectedDatabase().SearchKeyInDB(key)
	if baseType == nil {
		return nil, re.ErrNilValue
	}
	/**
	对baseType的类型是否为string进行校验。
	判断baseType的Encoding是否为string的Encoding,baseType的Type是否为RedisTypeString
	*/
	if _, ok := stringEncodingTypeDict[baseType.GetEncoding()]; !ok || baseType.GetObjectType() != encodings.RedisTypeString {
		loggers.Errorf(string(re.ErrWrongType), baseType.GetObjectType(), baseType.GetEncoding())
		return nil, re.ErrWrongType
	}
	if ts, ok := baseType.(database.TString); ok {
		return ts, nil
	}
	loggers.Errorf("base type can not convert to TString")
	return nil, re.ErrWrongType
}

func (handler *StringHandler) Append(cli *client.Client) {
	key := cli.Argv[1]
	ts, err := getTStringValueByKey(cli, key)
	if err != nil {
		cli.ResponseReError(err)
		return
	}
	// 如果TString的编码类型是int,转换成StringRaw再进行处理
	if ts.GetEncoding() == encodings.RedisEncodingInt {
		if valueInt, ok := ts.GetValue().(int64); !ok {
			cli.ResponseReError(re.ErrWrongTypeOrEncoding)
			return
		} else {
			rs := database.NewRedisStringWithEncodingRawString(fmt.Sprintf("%d", valueInt), -1)
			cli.SelectedDatabase().SetKeyInDB(key, rs)
			ts = rs
		}
	}
	cli.Response(ts.Append(cli.Argv[2]))
}

func (handler *StringHandler) Set(cli *client.Client) {
	key := cli.Argv[1]
	cli.SelectedDatabase().SetKeyInDB(key, database.NewRedisStringObject(cli.Argv[2]))
	cli.ResponseOK()
}

func (handler *StringHandler) SetNx(cli *client.Client) {
	key := cli.Argv[1]
	if cli.SelectedDatabase().SearchKeyInDB(key) == nil {
		cli.SelectedDatabase().SetKeyInDB(key, database.NewRedisStringObject(cli.Argv[2]))
		cli.Response(1)
	} else {
		cli.Response(0)
	}
}

func (handler *StringHandler) MSetNx(cli *client.Client) {
	args := cli.Argv
	if len(args)%2 == 0 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.GetOriginCommandName())
		return
	}
	var containsKey bool
	for i := 2; i < len(cli.Argv); i += 2 {
		if cli.SelectedDatabase().SearchKeyInDB(cli.Argv[i]) != nil {
			containsKey = true
			break
		}
	}
	if containsKey {
		cli.Response(0)
	} else {
		for i := 2; i < len(cli.Argv); i += 2 {
			cli.SelectedDatabase().SetKeyInDB(cli.Argv[i], database.NewRedisStringObject(cli.Argv[i+1]))
		}
		cli.Response(1)
	}
}

func (handler *StringHandler) Get(cli *client.Client) {
	key := cli.Argv[1]
	ts, err := getTStringValueByKey(cli, key)
	if err != nil {
		cli.ResponseReError(err)
		return
	}
	cli.Response(ts.String())
}

func (handler *StringHandler) GetSet(cli *client.Client) {
	key := cli.Argv[1]
	ts, err := getTStringValueByKey(cli, key)

	if err != nil && err != re.ErrNilValue {
		cli.ResponseReError(err)
		return
	}
	cli.SelectedDatabase().SetKeyInDB(key, database.NewRedisStringObject(cli.Argv[2]))

	if ts == nil {
		cli.ResponseReError(re.ErrNilValue)
	} else {
		cli.Response(ts.String())
	}
}

func (handler *StringHandler) MGet(cli *client.Client) {
	ret := make([]interface{}, 0)
	for _, key := range cli.Argv[1:] {
		ts, err := getTStringValueByKey(cli, key)
		if err != nil {
			ret = append(ret, nil)
		} else {
			ret = append(ret, ts.String())
		}
	}
	cli.Response(ret)
}

func (handler *StringHandler) MSet(cli *client.Client) {
	if len(cli.Argv)%2 == 0 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.GetOriginCommandName())
		return
	}
	for i := 1; i < len(cli.Argv); i += 2 {
		cli.SelectedDatabase().SetKeyInDB(cli.Argv[i], database.NewRedisStringObject(cli.Argv[i+1]))
	}
	cli.ResponseOK()
}

func (handler *StringHandler) Incr(cli *client.Client) {
	key := cli.Argv[1]
	ts, err := getTStringValueByKey(cli, key)
	if err != nil && err != re.ErrNilValue {
		cli.ResponseReError(err)
	} else if err == re.ErrNilValue {
		cli.SelectedDatabase().SetKeyInDB(key, database.NewRedisStringObject("1"))
		cli.Response(1)
	} else {
		if ret, err := ts.Incr(); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(ret)
		}
	}
}

func (handler *StringHandler) IncrBy(cli *client.Client) {
	key := cli.Argv[1]
	ts, err := getTStringValueByKey(cli, key)
	if err != nil && err != re.ErrNilValue {
		cli.ResponseReError(err)
	} else if err == re.ErrNilValue {
		if _, err := strconv.ParseInt(cli.Argv[1], 10, 64); err != nil {
			cli.ResponseReError(re.ErrNotIntegerOrOutOfRange)
		} else {
			cli.SelectedDatabase().SetKeyInDB(key, database.NewRedisStringObject(cli.Argv[2]))
			cli.Response(cli.Argv[1])
		}
	} else if ret, err := ts.IncrBy(cli.Argv[2]); err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(ret)
	}
}

func (handler *StringHandler) IncrByFloat(cli *client.Client) {
	key := cli.Argv[1]
	ts, err := getTStringValueByKey(cli, key)
	if err != nil && err != re.ErrNilValue {
		cli.ResponseReError(err)
	} else if err == re.ErrNilValue {
		if _, err := strconv.ParseFloat(cli.Argv[2], 64); err != nil {
			cli.ResponseReError(re.ErrValueIsNotFloat)
		} else {
			cli.SelectedDatabase().SetKeyInDB(key, database.NewRedisStringObject(cli.Argv[2]))
			cli.Response(cli.Argv[2])
		}
	} else {
		// 如果TString的编码类型是int,转换成StringRaw再进行处理
		if ts.GetEncoding() == encodings.RedisEncodingInt {
			if valueInt, ok := ts.GetValue().(int64); !ok {
				cli.ResponseReError(re.ErrWrongTypeOrEncoding)
				return
			} else {
				rs := database.NewRedisStringWithEncodingRawString(fmt.Sprintf("%d", valueInt), -1)
				cli.SelectedDatabase().SetKeyInDB(key, rs)
				ts = rs
			}
		}
		if ret, err := ts.IncrByFloat(cli.Argv[2]); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(ret)
		}
	}
}

func (handler *StringHandler) Decr(cli *client.Client) {
	key := cli.Argv[1]
	ts, err := getTStringValueByKey(cli, key)
	if err != nil && err != re.ErrNilValue {
		cli.ResponseReError(err)
	} else if err == re.ErrNilValue {
		cli.SelectedDatabase().SetKeyInDB(key, database.NewRedisStringObject("-1"))
		cli.Response("-1")
	} else {
		if ret, err := ts.Decr(); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(ret)
		}
	}
}

func (handler *StringHandler) DecrBy(cli *client.Client) {
	key := cli.Argv[1]
	ts, err := getTStringValueByKey(cli, key)
	if err != nil && err != re.ErrNilValue {
		cli.ResponseReError(err)
	} else if err == re.ErrNilValue {
		if _, err := strconv.ParseInt(cli.Argv[2], 10, 64); err != nil {
			cli.ResponseReError(re.ErrNotIntegerOrOutOfRange)
		} else {
			cli.SelectedDatabase().SetKeyInDB(key, database.NewRedisStringObject(cli.Argv[2]))
			cli.Response(cli.Argv[2])
		}
	} else {
		if ret, err := ts.DecrBy(cli.Argv[2]); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(ret)
		}
	}
}

func (handler *StringHandler) Strlen(cli *client.Client) {
	key := cli.Argv[1]
	ts, err := getTStringValueByKey(cli, key)
	if err != nil {
		cli.ResponseReError(err)
		return
	}
	cli.Response(ts.Strlen())
}

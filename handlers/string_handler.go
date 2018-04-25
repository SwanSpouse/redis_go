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
	switch cli.Cmd.GetName() {
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
		cli.ResponseReError(re.ErrUnknownCommand, cli.Cmd.GetOriginName())
	}
	// 最后统一发送数据
	cli.Flush()
}

func getTStringValueByKey(cli *client.Client, key string) (database.TString, error) {
	if cli.Cmd == nil {
		return nil, re.ErrNilCommand
	}
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
	args := cli.Cmd.GetArgs()
	if len(args) < 2 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
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
	cli.Response(ts.Append(args[1]))
}

func (handler *StringHandler) Set(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) < 2 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	cli.SelectedDatabase().SetKeyInDB(key, database.NewRedisStringObject(args[1]))
	cli.ResponseOK()
}

func (handler *StringHandler) SetNx(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) < 2 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if cli.SelectedDatabase().SearchKeyInDB(key) == nil {
		cli.SelectedDatabase().SetKeyInDB(key, database.NewRedisStringObject(args[1]))
		cli.Response(1)
	} else {
		cli.Response(0)
	}
}

func (handler *StringHandler) MSetNx(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) < 2 || len(args)%2 == 1 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	var containsKey bool
	for i := 0; i < len(args); i += 2 {
		if cli.SelectedDatabase().SearchKeyInDB(args[i]) != nil {
			containsKey = true
			break
		}
	}
	if containsKey {
		cli.Response(0)
	} else {
		for i := 0; i < len(args); i += 2 {
			cli.SelectedDatabase().SetKeyInDB(args[i], database.NewRedisStringObject(args[i+1]))
		}
		cli.Response(1)
	}
}

func (handler *StringHandler) Get(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	// 判断参数个数是否合理
	if len(args) != 1 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	ts, err := getTStringValueByKey(cli, key)
	if err != nil {
		cli.ResponseReError(err)
		return
	}
	cli.Response(ts.String())
}

func (handler *StringHandler) GetSet(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) < 2 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	ts, err := getTStringValueByKey(cli, key)

	if err != nil && err != re.ErrNilValue {
		cli.ResponseReError(err)
		return
	}
	cli.SelectedDatabase().SetKeyInDB(key, database.NewRedisStringObject(args[1]))

	if ts == nil {
		cli.ResponseReError(re.ErrNilValue)
	} else {
		cli.Response(ts.String())
	}
}

func (handler *StringHandler) MGet(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) < 1 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	ret := make([]interface{}, 0)
	for _, key := range args {
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
	args := cli.Cmd.GetArgs()
	if len(args) < 2 || len(args)%2 == 1 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	for i := 0; i < len(args); i += 2 {
		cli.SelectedDatabase().SetKeyInDB(args[i], database.NewRedisStringObject(args[i+1]))
	}
	cli.ResponseOK()
}

func (handler *StringHandler) Incr(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) != 1 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
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
	args := cli.Cmd.GetArgs()
	if len(args) != 2 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	ts, err := getTStringValueByKey(cli, key)
	if err != nil && err != re.ErrNilValue {
		cli.ResponseReError(err)
	} else if err == re.ErrNilValue {
		if _, err := strconv.ParseInt(args[1], 10, 64); err != nil {
			cli.ResponseReError(re.ErrNotIntegerOrOutOfRange)
		} else {
			cli.SelectedDatabase().SetKeyInDB(key, database.NewRedisStringObject(args[1]))
			cli.Response(args[1])
		}
	} else if ret, err := ts.IncrBy(args[1]); err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(ret)
	}
}

func (handler *StringHandler) IncrByFloat(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) != 2 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	ts, err := getTStringValueByKey(cli, key)
	if err != nil && err != re.ErrNilValue {
		cli.ResponseReError(err)
	} else if err == re.ErrNilValue {
		if _, err := strconv.ParseFloat(args[1], 64); err != nil {
			cli.ResponseReError(re.ErrValueIsNotFloat)
		} else {
			cli.SelectedDatabase().SetKeyInDB(key, database.NewRedisStringObject(args[1]))
			cli.Response(args[1])
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
		if ret, err := ts.IncrByFloat(args[1]); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(ret)
		}
	}
}

func (handler *StringHandler) Decr(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) != 1 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
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
	args := cli.Cmd.GetArgs()
	if len(args) != 2 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	ts, err := getTStringValueByKey(cli, key)
	if err != nil && err != re.ErrNilValue {
		cli.ResponseReError(err)
	} else if err == re.ErrNilValue {
		if _, err := strconv.ParseInt(args[1], 10, 64); err != nil {
			cli.ResponseReError(re.ErrNotIntegerOrOutOfRange)
		} else {
			cli.SelectedDatabase().SetKeyInDB(key, database.NewRedisStringObject(args[1]))
			cli.Response(args[1])
		}
	} else {
		if ret, err := ts.DecrBy(args[1]); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(ret)
		}
	}
}

func (handler *StringHandler) Strlen(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) != 1 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	ts, err := getTStringValueByKey(cli, key)
	if err != nil {
		cli.ResponseReError(err)
		return
	}
	cli.Response(ts.Strlen())
}

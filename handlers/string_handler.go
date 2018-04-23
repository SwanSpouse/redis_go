package handlers

import (
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
	RedisStringCommandSetNX       = "SETNX"
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
	if key, ts, err := handler.getValidKeyAndTypeOrError(cli); err == nil {
		switch cli.Cmd.GetName() {
		case RedisStringCommandAppend:
			handler.Append(cli, ts)
		case RedisStringCommandBitCount, RedisStringCommandBitop, RedisStringCommandGetBit, RedisStringCommandSetBit:
			cli.ResponseReError(re.ErrFunctionNotImplement)
		case RedisStringCommandDecr:
			handler.Decr(cli, ts)
		case RedisStringCommandDecrBy:
		case RedisStringCommandGet:
			handler.Get(cli, ts)
		case RedisStringCommandGetRange:
		case RedisStringCommandGetSet:
		case RedisStringCommandIncr:
			handler.Incr(cli, ts)
		case RedisStringCommandIncrBy:
		case RedisStringCommandIncrByFloat:
		case RedisStringCommandMGet:
		case RedisStringCommandMSet:
		case RedisStringCommandMSetNx:
		case RedisStringCommandPSetEx:
		case RedisStringCommandSet:
			handler.Set(cli, key)
		case RedisStringCommandSetNX:
		case RedisStringCommandSetEX:
		case RedisStringCommandSetRange:
		case RedisStringCommandStrLen:
			handler.Strlen(cli, ts)
		default:
			cli.ResponseReError(re.ErrUnknownCommand, cli.Cmd.GetOriginName())
		}
	} else {
		cli.ResponseReError(err)
	}
	// 最后统一发送数据
	cli.Flush()
}

func (handler *StringHandler) getValidKeyAndTypeOrError(cli *client.Client) (string, database.TString, error) {
	if cli.Cmd == nil {
		return "", nil, re.ErrNilCommand
	}
	args := cli.Cmd.GetArgs()
	// 参数个数的错误都交给每个命令自己来进行处理
	if len(args) == 0 {
		return "", nil, nil
	}
	key := args[0]
	if cli.Cmd.GetName() == RedisStringCommandSet {
		return key, nil, nil
	}
	// 获取key在数据库中对应的value(TBase:BaseType)
	baseType := cli.SelectedDatabase().SearchKeyInDB(key)
	if baseType == nil {
		return "", nil, re.ErrNilValue
	}
	/**
	对baseType的类型是否为string进行校验。
	判断baseType的Encoding是否为string的Encoding,baseType的Type是否为RedisTypeString
	*/
	if _, ok := stringEncodingTypeDict[baseType.GetEncoding()]; !ok || baseType.GetObjectType() != encodings.RedisTypeString {
		loggers.Errorf(string(re.ErrWrongTypeOrEncoding), baseType.GetObjectType(), baseType.GetEncoding())
		return "", nil, re.ErrConvertToTargetType
	}
	if ts, ok := baseType.(database.TString); ok {
		return "", ts, nil
	}
	loggers.Errorf("base type can not convert to TString")
	return "", nil, re.ErrConvertToTargetType
}

func (handler *StringHandler) Append(cli *client.Client, ts database.TString) {
	args := cli.Cmd.GetArgs()
	if len(args) < 2 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	// 如果TString的编码类型是int,转换成StringRaw再进行处理
	if ts.GetEncoding() == encodings.RedisEncodingInt {
		if valueInt, ok := ts.GetValue().(int); !ok {
			cli.ResponseReError(re.ErrWrongTypeOrEncoding)
			return
		} else {
			rs := database.NewRedisStringWithEncodingRawString(strconv.Itoa(valueInt), -1)
			cli.SelectedDatabase().SetKeyInDB(key, rs)
			ts = rs
		}
	}
	cli.Response(ts.Append(args[1]))
}

func (handler *StringHandler) Set(cli *client.Client, key string) {
	args := cli.Cmd.GetArgs()
	if len(args) < 2 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	cli.SelectedDatabase().SetKeyInDB(key, database.NewRedisStringObject(args[1]))
	cli.ResponseOK()
}

func (handler *StringHandler) Get(cli *client.Client, ts database.TString) {
	args := cli.Cmd.GetArgs()
	// 判断参数个数是否合理
	if len(args) != 1 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	cli.Response(ts.String())
}

func (handler *StringHandler) Incr(cli *client.Client, ts database.TString) {
	args := cli.Cmd.GetArgs()
	if len(args) != 1 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	if ret, err := ts.Incr(); err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(ret)
	}
}

func (handler *StringHandler) Decr(cli *client.Client, ts database.TString) {
	args := cli.Cmd.GetArgs()
	if len(args) != 1 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	if ret, err := ts.Decr(); err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(ret)
	}
}

func (handler *StringHandler) Strlen(cli *client.Client, ts database.TString) {
	args := cli.Cmd.GetArgs()
	if len(args) != 1 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	cli.Response(ts.Strlen())
}

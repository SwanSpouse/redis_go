package handlers

import (
	"redis_go/client"
	"redis_go/database"
	"redis_go/encodings"
	re "redis_go/error"
	"redis_go/loggers"
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

func (handler *StringHandler) Process(client *client.Client) {
	if key, ts, err := handler.getValidKeyAndTypeOrError(client); err == nil {
		switch client.Cmd.GetName() {
		case RedisStringCommandAppend:
			handler.Append(client, ts)
		case RedisStringCommandBitCount, RedisStringCommandBitop, RedisStringCommandGetBit, RedisStringCommandSetBit:
			client.ResponseReError(re.ErrFunctionNotImplement)
		case RedisStringCommandDecr:
			handler.Decr(client, ts)
		case RedisStringCommandDecrBy:
		case RedisStringCommandGet:
			handler.Get(client, ts)
		case RedisStringCommandGetRange:
		case RedisStringCommandGetSet:
		case RedisStringCommandIncr:
			handler.Incr(client, ts)
		case RedisStringCommandIncrBy:
		case RedisStringCommandIncrByFloat:
		case RedisStringCommandMGet:
		case RedisStringCommandMSet:
		case RedisStringCommandMSetNx:
		case RedisStringCommandPSetEx:
		case RedisStringCommandSet:
			handler.Set(client, key)
		case RedisStringCommandSetNX:
		case RedisStringCommandSetEX:
		case RedisStringCommandSetRange:
		case RedisStringCommandStrLen:
			handler.Strlen(client, ts)
		default:
			client.ResponseReError(re.ErrUnknownCommand, client.Cmd.GetOriginName())
		}
	} else {
		client.ResponseReError(err)
	}
	// 最后统一发送数据
	client.Flush()
}

func (handler *StringHandler) getValidKeyAndTypeOrError(client *client.Client) (string, database.TString, error) {
	if client.Cmd == nil {
		return "", nil, re.ErrNilCommand
	}
	args := client.Cmd.GetArgs()
	// 参数个数的错误都交给每个命令自己来进行处理
	if len(args) == 0 {
		return "", nil, nil
	}
	key := args[0]
	if client.Cmd.GetName() == RedisStringCommandSet {
		return key, nil, nil
	}
	// 获取key在数据库中对应的value(TBase:BaseType)
	baseType := client.SelectedDatabase().SearchKeyInDB(key)
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

func (handler *StringHandler) Append(client *client.Client, ts database.TString) {
	args := client.Cmd.GetArgs()
	if len(args) < 2 {
		client.ResponseReError(re.ErrWrongNumberOfArgs, client.Cmd.GetOriginName())
		return
	}
	// TODO lmj 应该在这里把类型判断的拦截做了
	client.Response(ts.Append(args[1]))
}

func (handler *StringHandler) Set(client *client.Client, key string) {
	args := client.Cmd.GetArgs()
	if len(args) < 2 {
		client.ResponseReError(re.ErrWrongNumberOfArgs, client.Cmd.GetOriginName())
		return
	}
	client.SelectedDatabase().SetKeyInDB(key, database.NewRedisStringObject(args[1]))
	client.ResponseOK()
}

func (handler *StringHandler) Get(client *client.Client, ts database.TString) {
	args := client.Cmd.GetArgs()
	// 判断参数个数是否合理
	if len(args) != 1 {
		client.ResponseReError(re.ErrWrongNumberOfArgs, client.Cmd.GetOriginName())
		return
	}
	client.Response(ts.String())
}

func (handler *StringHandler) Incr(client *client.Client, ts database.TString) {
	args := client.Cmd.GetArgs()
	if len(args) != 1 {
		client.ResponseReError(re.ErrWrongNumberOfArgs, client.Cmd.GetOriginName())
		return
	}
	if ret, err := ts.Incr(); err != nil {
		client.ResponseReError(err)
	} else {
		client.Response(ret)
	}
}

func (handler *StringHandler) Decr(client *client.Client, ts database.TString) {
	args := client.Cmd.GetArgs()
	if len(args) != 1 {
		client.ResponseReError(re.ErrWrongNumberOfArgs, client.Cmd.GetOriginName())
		return
	}
	if ret, err := ts.Decr(); err != nil {
		client.ResponseReError(err)
	} else {
		client.Response(ret)
	}
}

func (handler *StringHandler) Strlen(client *client.Client, ts database.TString) {
	args := client.Cmd.GetArgs()
	if len(args) != 1 {
		client.ResponseReError(re.ErrWrongNumberOfArgs, client.Cmd.GetOriginName())
		return
	}
	client.Response(ts.Strlen())
}

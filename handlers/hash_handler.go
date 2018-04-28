package handlers

import (
	"redis_go/client"
	"redis_go/database"
	"redis_go/encodings"
	re "redis_go/error"
	"redis_go/loggers"
)

const (
	RedisHashCommandHDel         = "HDEL"
	RedisHashCommandHExists      = "HEXISTS"
	RedisHashCommandHGet         = "HGET"
	RedisHashCommandHGetAll      = "HGETALL"
	RedisHashCommandHIncrBy      = "HINCRBY"
	RedisHashCommandHIncrByFloat = "HINCRBYFLOAT"
	RedisHashCommandHKeys        = "HKEYS"
	RedisHashCommandHLen         = "HLEN"
	RedisHashCommandHMGet        = "HMGET"
	RedisHashCommandHMSet        = "HMSET"
	RedisHashCommandHSet         = "HSET"
	RedisHashCommandHSetNX       = "HSETNX"
	RedisHashCommandHVals        = "HVALS"
	RedisHashCommandHScan        = "HSCAN"
	RedisHashCommandHStrLen      = "HSTRLEN"
	RedisHashCommandHDebug       = "HDEBUG"
)

// HashHandler可以处理的一种rawType
var hashEncodingTypeDict = map[string]bool{
	encodings.RedisEncodingHT: true,
}

type HashHandler struct {
}

func (handler *HashHandler) Process(cli *client.Client) {
	switch cli.Cmd.GetName() {
	case RedisHashCommandHDel:
		handler.HDel(cli)
	case RedisHashCommandHExists:
		handler.HExists(cli)
	case RedisHashCommandHGet:
		handler.HGet(cli)
	case RedisHashCommandHGetAll:
		handler.HGetAll(cli)
	case RedisHashCommandHIncrBy:
	case RedisHashCommandHIncrByFloat:
	case RedisHashCommandHKeys:
		handler.HKeys(cli)
	case RedisHashCommandHLen:
		handler.HLen(cli)
	case RedisHashCommandHMGet:
		handler.HMGet(cli)
	case RedisHashCommandHMSet:
		handler.HMSet(cli)
	case RedisHashCommandHSet:
		handler.HSet(cli)
	case RedisHashCommandHSetNX:
		handler.HSetNX(cli)
	case RedisHashCommandHVals:
		handler.HVals(cli)
	case RedisHashCommandHScan:
		handler.HScan(cli)
	case RedisHashCommandHStrLen:
		handler.HStrLen(cli)
	case RedisHashCommandHDebug:
		handler.HDebug(cli)
	default:
		cli.ResponseReError(re.ErrUnknownCommand, cli.Cmd.GetOriginName())
	}
	// 最后统一发送数据
	cli.Flush()
}

func getTHashValueByKey(cli *client.Client, key string) (database.THash, error) {
	if cli.Cmd == nil {
		return nil, re.ErrNilCommand
	}
	baseType := cli.SelectedDatabase().SearchKeyInDB(key)
	if baseType == nil {
		return nil, re.ErrNilValue
	}
	if _, ok := hashEncodingTypeDict[baseType.GetEncoding()]; !ok || baseType.GetObjectType() != encodings.RedisTypeHash {
		loggers.Errorf(string(re.ErrWrongType), baseType.GetObjectType(), baseType.GetEncoding())
		return nil, re.ErrWrongType
	}
	if th, ok := baseType.(database.THash); ok {
		return th, nil
	}
	loggers.Errorf("base type can not convert to THash")
	return nil, re.ErrWrongType
}

func createHashIfNotExists(cli *client.Client, key string) error {
	if _, err := getTHashValueByKey(cli, key); err != nil && err != re.ErrNoSuchKey {
		return err
	} else if err == re.ErrNoSuchKey {
		cli.SelectedDatabase().SetKeyInDB(key, database.NewRedisHashObject())
	}
	return nil
}

func (handler *HashHandler) HDel(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) < 2 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(th.HDel(args[1:]))
	}
}

func (handler *HashHandler) HExists(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) != 2 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(th.HExists(args[1]))
	}
}

func (handler *HashHandler) HGet(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) != 2 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		if ret, err := th.HGet(args[1]); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(ret)
		}
	}
}

func (handler *HashHandler) HGetAll(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) != 1 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(th.HGetAll())
	}
}

func (handler *HashHandler) HKeys(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) != 1 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		ret := make([]string, 0)
		keyValues := th.HGetAll()
		for i := 0; i < len(keyValues); i += 2 {
			ret = append(ret, keyValues[i])
		}
		cli.Response(ret)
	}
}

func (handler *HashHandler) HLen(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) != 1 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(th.HLen())
	}
}

func (handler *HashHandler) HMGet(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) < 2 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		ret := make([]string, 0)
		for _, item := range args[1:] {
			value, _ := th.HGet(item)
			ret = append(ret, value)
		}
		cli.Response(ret)
	}
}

func (handler *HashHandler) HMSet(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) < 2 || len(args)%2 == 0 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if err := createHashIfNotExists(cli, key); err != nil {
		cli.ResponseReError(err)
		return
	}
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		for i := 1; i < len(args); i += 2 {
			th.HSet(args[i], args[i+1])
		}
		cli.ResponseOK()
	}
}

func (handler *HashHandler) HSet(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) < 3 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if err := createHashIfNotExists(cli, key); err != nil {
		cli.ResponseReError(err)
		return
	}
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(th.HSet(args[1], args[2]))
	}
}

func (handler *HashHandler) HSetNX(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) < 3 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if err := createHashIfNotExists(cli, key); err != nil {
		cli.ResponseReError(err)
		return
	}
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		if val, _ := th.HGet(args[1]); val == "" {
			cli.Response(th.HSet(args[1], args[2]))
		}
	}
}

func (handler *HashHandler) HVals(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) != 1 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		ret := make([]string, 0)
		keyValues := th.HGetAll()
		for i := 0; i < len(keyValues); i += 2 {
			ret = append(ret, keyValues[i+1])
		}
		cli.Response(ret)
	}
}

func (handler *HashHandler) HScan(cli *client.Client) {

}

func (handler *HashHandler) HStrLen(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) != 2 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		if ret, err := th.HGet(args[1]); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(len(ret))
		}
	}
}

func (handler *HashHandler) HDebug(cli *client.Client) {
	args := cli.Cmd.GetArgs()
	if len(args) != 1 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.Cmd.GetOriginName())
		return
	}
	key := args[0]
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		th.HDebug()
		cli.ResponseOK()
	}
}

package handlers

import (
	"redis_go/client"
	"redis_go/database"
	"redis_go/encodings"
	re "redis_go/error"
	"redis_go/loggers"
)

var (
	_ client.BaseHandler = (*HashHandler)(nil)
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
	switch cli.GetCommandName() {
	case RedisHashCommandHDel:
		handler.HDel(cli)
	case RedisHashCommandHExists:
		handler.HExists(cli)
	case RedisHashCommandHGet:
		handler.HGet(cli)
	case RedisHashCommandHGetAll:
		handler.HGetAll(cli)
	case RedisHashCommandHIncrBy:
		handler.HIncrBy(cli)
	case RedisHashCommandHIncrByFloat:
		handler.HIncrByFloat(cli)
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
		//handler.HScan(cli)
	case RedisHashCommandHStrLen:
		handler.HStrLen(cli)
	case RedisHashCommandHDebug:
		handler.HDebug(cli)
	default:
		cli.ResponseReError(re.ErrUnknownCommand, cli.GetOriginCommandName())
	}
	// 最后统一发送数据
	cli.Flush()
}

func getTHashValueByKey(cli *client.Client, key string) (database.THash, error) {
	baseType := cli.SelectedDatabase().SearchKeyInDB(key)
	if baseType == nil {
		return nil, re.ErrNoSuchKey
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
	key := cli.Argv[1]
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(th.HDel(cli.Argv[2:]))
	}
}

func (handler *HashHandler) HExists(cli *client.Client) {
	key := cli.Argv[1]
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(th.HExists(cli.Argv[2]))
	}
}

func (handler *HashHandler) HGet(cli *client.Client) {
	key := cli.Argv[1]
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		if ret, err := th.HGet(cli.Argv[2]); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(ret)
		}
	}
}

func (handler *HashHandler) HGetAll(cli *client.Client) {
	key := cli.Argv[1]
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(th.HGetAll())
	}
}

func (handler *HashHandler) HKeys(cli *client.Client) {
	key := cli.Argv[1]
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
	key := cli.Argv[1]
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(th.HLen())
	}
}

func (handler *HashHandler) HMGet(cli *client.Client) {
	key := cli.Argv[1]
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		ret := make([]string, 0)
		for _, item := range cli.Argv[2:] {
			value, _ := th.HGet(item)
			ret = append(ret, value)
		}
		cli.Response(ret)
	}
}

func (handler *HashHandler) HMSet(cli *client.Client) {
	key := cli.Argv[1]
	if len(cli.Argv)%2 == 1 {
		cli.ResponseReError(re.ErrWrongNumberOfArgs, cli.GetOriginCommandName())
		return
	}
	if err := createHashIfNotExists(cli, key); err != nil {
		cli.ResponseReError(err)
		return
	}
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		for i := 2; i < len(cli.Argv); i += 2 {
			th.HSet(cli.Argv[i], cli.Argv[i+1])
		}
		cli.ResponseOK()
	}
}

func (handler *HashHandler) HSet(cli *client.Client) {
	key := cli.Argv[1]
	if err := createHashIfNotExists(cli, key); err != nil {
		cli.ResponseReError(err)
		return
	}
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(th.HSet(cli.Argv[2], cli.Argv[3]))
	}
}

func (handler *HashHandler) HSetNX(cli *client.Client) {
	key := cli.Argv[1]
	if err := createHashIfNotExists(cli, key); err != nil {
		cli.ResponseReError(err)
		return
	}
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		if val, _ := th.HGet(cli.Argv[2]); val == "" {
			cli.Response(th.HSet(cli.Argv[2], cli.Argv[3]))
		} else {
			cli.Response(0)
		}
	}
}

func (handler *HashHandler) HVals(cli *client.Client) {
	key := cli.Argv[1]
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

func (handler *HashHandler) HStrLen(cli *client.Client) {
	key := cli.Argv[1]
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		if ret, err := th.HGet(cli.Argv[2]); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(len(ret))
		}
	}
}

func (handler *HashHandler) HIncrBy(cli *client.Client) {
	key := cli.Argv[1]
	if err := createHashIfNotExists(cli, key); err != nil {
		cli.ResponseReError(err)
		return
	}
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		if ret, err := th.HIncrBy(cli.Argv[2], cli.Argv[3]); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(ret)
		}
	}
}

func (handler *HashHandler) HIncrByFloat(cli *client.Client) {
	key := cli.Argv[1]
	if err := createHashIfNotExists(cli, key); err != nil {
		cli.ResponseReError(err)
		return
	}
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		if ret, err := th.HIncrByFloat(cli.Argv[2], cli.Argv[3]); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(ret)
		}
	}
}

func (handler *HashHandler) HDebug(cli *client.Client) {
	key := cli.Argv[1]
	if th, err := getTHashValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		th.HDebug()
		cli.ResponseOK()
	}
}

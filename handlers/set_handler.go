package handlers

import (
	"github.com/SwanSpouse/redis_go/client"
	"github.com/SwanSpouse/redis_go/database"
	"github.com/SwanSpouse/redis_go/encodings"
	re "github.com/SwanSpouse/redis_go/error"
	"github.com/SwanSpouse/redis_go/loggers"
)

const (
	RedisSetCommandSADD        = "SADD"
	RedisSetCommandSCARD       = "SCARD"
	RedisSetCommandSDIFF       = "SDIFF"
	RedisSetCommandSDIFFSTORE  = "SDIFFSTORE"
	RedisSetCommandSINTER      = "SINTER"
	RedisSetCommandSINTERSTORE = "SINTERSTORE"
	RedisSetCommandSISMEMBER   = "SISMEMBER"
	RedisSetCommandSMEMBERS    = "SMEMBERS"
	RedisSetCommandSMOVE       = "SMOVE"
	RedisSetCommandSPOP        = "SPOP"
	RedisSetCommandSRANDMEMBER = "SRANDMEMBER"
	RedisSetCommandSREM        = "SREM"
	RedisSetCommandSUNION      = "SUNION"
	RedisSetCommandSUNIONSTORE = "SUNIONSTORE"
	RedisSetCommandSSCAN       = "SSCAN"
)

// SetHandler可以处理的一种rawType
var setEncodingTypeDict = map[string]bool{
	encodings.RedisEncodingHT: true,
}

type SetHandler struct{}

func getTSetValueByKey(cli *client.Client, key string) (database.TSet, error) {
	baseType := cli.SelectedDatabase().SearchKeyInDB(key)
	if baseType == nil {
		return nil, re.ErrNoSuchKey
	}
	if _, ok := setEncodingTypeDict[baseType.GetEncoding()]; !ok || baseType.GetObjectType() != encodings.RedisTypeSet {
		loggers.Errorf(string(re.ErrWrongType), baseType.GetObjectType(), baseType.GetEncoding())
		return nil, re.ErrWrongType
	}
	if ts, ok := baseType.(database.TSet); ok {
		return ts, nil
	}
	loggers.Errorf("base type:%+v can not convert to TSet", baseType)
	return nil, re.ErrWrongType
}

func createSetIfNotExists(cli *client.Client, key string) error {
	if _, err := getTSetValueByKey(cli, key); err != nil && err != re.ErrNoSuchKey {
		return err
	} else if err == re.ErrNoSuchKey {
		cli.SelectedDatabase().SetKeyInDB(key, database.NewRedisSetObject())
	}
	return nil
}

func (handler *SetHandler) SAdd(cli *client.Client) {
	key := cli.Argv[1]
	if err := createSetIfNotExists(cli, key); err != nil {
		cli.ResponseReError(err)
		return
	}
	if ts, err := getTSetValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		ret := ts.SAdd(cli.Argv[2:])
		cli.Response(ret)
		if ret != 0 {
			cli.Dirty += 1
		}
	}
}

func (handler *SetHandler) SCard(cli *client.Client) {
	key := cli.Argv[1]
	if ts, err := getTSetValueByKey(cli, key); err != nil && err == re.ErrNoSuchKey {
		cli.Response(0)
	} else if err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(ts.SCard())
	}
}

func (handler *SetHandler) SDiff(cli *client.Client) {

}

func (handler *SetHandler) SDiffStore(cli *client.Client) {

}

func (handler *SetHandler) SInter(cli *client.Client) {

}

func (handler *SetHandler) SInterStore(cli *client.Client) {

}

func (handler *SetHandler) SIsMember(cli *client.Client) {
	key := cli.Argv[1]
	if ts, err := getTSetValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(ts.SIsMember(cli.Argv[2]))
	}
}

func (handler *SetHandler) SMembers(cli *client.Client) {
	key := cli.Argv[1]
	if ts, err := getTSetValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		cli.Response(ts.SMembers())
	}
}

func (handler *SetHandler) SMove(cli *client.Client) {

}

func (handler *SetHandler) SPop(cli *client.Client) {
	key := cli.Argv[1]
	if ts, err := getTSetValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		if key, err := ts.SPop(); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(key)
			cli.Dirty += 1
		}
	}
}

func (handler *SetHandler) SRem(cli *client.Client) {
	key := cli.Argv[1]
	if ts, err := getTSetValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		ret := ts.SRem(cli.Argv[2:])
		cli.Response(ret)
		if ret != 0 {
			cli.Dirty += 1
		}
	}
}

func (handler *SetHandler) SRandMember(cli *client.Client) {
	key := cli.Argv[1]
	if ts, err := getTSetValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		if key, err := ts.SRandMember(); err != nil {
			cli.ResponseReError(err)
		} else {
			cli.Response(key)
		}
	}
}

func (handler *SetHandler) SDebug(cli *client.Client) {
	key := cli.Argv[1]
	if ts, err := getTSetValueByKey(cli, key); err != nil {
		cli.ResponseReError(err)
	} else {
		ts.SDebug()
		cli.ResponseOK()
	}
}

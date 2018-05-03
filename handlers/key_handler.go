package handlers

import (
	"redis_go/client"
	re "redis_go/error"
	"strings"
)

var (
	_ client.BaseHandler = (*KeyHandler)(nil)
)

const (
	RedisKeyCommandDel    = "DEL"
	RedisKeyCommandExists = "EXISTS"
	RedisKeyCommandType   = "TYPE"
	RedisKeyCommandObject = "OBJECT"
)

const (
	CommandObjectSubTypeRefCount  = "REFCOUNT"
	CommandObjectSubTypeEncodings = "ENCODING"
	CommandObjectSubTypeIdleTime  = "IDLETIME"
)

type KeyHandler struct{}

func (handler *KeyHandler) Process(cli *client.Client) {
	if cli.Cmd == nil {
		cli.ResponseReError(re.ErrNilCommand)
		return
	}
	switch cli.GetCommandName() {
	case RedisKeyCommandDel:
		handler.Del(cli)
	case "DUMP":
	case RedisKeyCommandExists:
		handler.Exists(cli)
	case "EXPIRE":
	case "EXPIREAT":
	case "KEYS":
	case "MIGRATE":
	case "MOVE":
	case RedisKeyCommandObject:
		handler.Object(cli)
	case "PERSIST":
	case "PEXPIRE":
	case "PEXPIREAT":
	case "PTTL":
	case "RANDOMKEY":
	case "RENAME":
	case "RENAMENX":
	case "RESTORE":
	case "SORT":
	case "TTL":
	case RedisKeyCommandType:
		handler.Type(cli)
	case "SCAN":
	default:
		cli.ResponseReError(re.ErrUnknownCommand, cli.GetOriginCommandName())
	}
	cli.Flush()
}

func (handler *KeyHandler) Del(cli *client.Client) {
	successCount := cli.SelectedDatabase().RemoveKeyInDB(cli.Argv)
	cli.Response(successCount)
}

func (handler *KeyHandler) Exists(cli *client.Client) {
	successCount, _ := cli.SelectedDatabase().SearchKeysInDB(cli.Argv)
	cli.Response(int64(len(successCount)))
}

func (handler *KeyHandler) Type(cli *client.Client) {
	if tb := cli.SelectedDatabase().SearchKeyInDB(cli.Argv[1]); tb == nil {
		cli.Response(nil)
	} else {
		cli.Response(tb.GetObjectType())
	}
}

func (handler *KeyHandler) Object(cli *client.Client) {
	encodings := cli.Argv[1]
	key := cli.Argv[2]
	switch strings.ToUpper(encodings) {
	case CommandObjectSubTypeEncodings:
		if tb := cli.SelectedDatabase().SearchKeyInDB(key); tb == nil {
			cli.Response(nil)
		} else {
			cli.Response(tb.GetEncoding())
		}
	case CommandObjectSubTypeRefCount:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case CommandObjectSubTypeIdleTime:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	default:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	}
}

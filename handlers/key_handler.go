package handlers

import (
	"math/rand"
	"redis_go/client"
	re "redis_go/error"
	"strings"
	"time"
)

var (
	_ client.BaseHandler = (*KeyHandler)(nil)
)

const (
	RedisKeyCommandDel       = "DEL"
	RedisKeyCommandExists    = "EXISTS"
	RedisKeyCommandType      = "TYPE"
	RedisKeyCommandObject    = "OBJECT"
	RedisKeyCommandDump      = "DUMP"
	RedisKeyCommandExpire    = "EXPIRE"
	RedisKeyCommandExpireAt  = "EXPIREAT"
	RedisKeyCommandKeys      = "KEYS"
	RedisKeyCommandMigrate   = "MIGRATE"
	RedisKeyCommandMove      = "MOVE"
	RedisKeyCommandPersist   = "PERSIST"
	RedisKeyCommandPExpire   = "PEXPIRE"
	RedisKeyCommandPExpireAt = "PEXPIREAT"
	RedisKeyCommandPTTL      = "PTTL"
	RedisKeyCommandRandomKey = "RANDOMKEY"
	RedisKeyCommandRename    = "RENAME"
	RedisKeyCommandRenameNx  = "RENAMENX"
	RedisKeyCommandRestore   = "RESTORE"
	RedisKeyCommandSort      = "SORT"
	RedisKeyCommandTTL       = "TTL"
	RedisKeyCommandScan      = "SCAN"
)

const (
	CommandObjectSubTypeRefCount  = "REFCOUNT"
	CommandObjectSubTypeEncodings = "ENCODING"
	CommandObjectSubTypeIdleTime  = "IDLETIME"
)

type KeyHandler struct{}

func (handler *KeyHandler) Process(cli *client.Client) {
	switch cli.GetCommandName() {
	case RedisKeyCommandDel:
		handler.Del(cli)
	case RedisKeyCommandDump:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisKeyCommandExists:
		handler.Exists(cli)
	case RedisKeyCommandExpire:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisKeyCommandExpireAt:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisKeyCommandKeys:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisKeyCommandMigrate:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisKeyCommandMove:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisKeyCommandObject:
		handler.Object(cli)
	case RedisKeyCommandPersist:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisKeyCommandPExpire:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisKeyCommandPExpireAt:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisKeyCommandPTTL:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisKeyCommandRandomKey:
		handler.RandomKey(cli)
	case RedisKeyCommandRename:
		handler.Rename(cli)
	case RedisKeyCommandRenameNx:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisKeyCommandRestore:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisKeyCommandSort:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisKeyCommandTTL:
		cli.ResponseReError(re.ErrFunctionNotImplement)
	case RedisKeyCommandType:
		handler.Type(cli)
	case RedisKeyCommandScan:
		cli.ResponseReError(re.ErrFunctionNotImplement)
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
	cli.Response(len(successCount))
}

func (handler *KeyHandler) Type(cli *client.Client) {
	if tb := cli.SelectedDatabase().SearchKeyInDB(cli.Argv[1]); tb == nil {
		cli.Response(nil)
	} else {
		cli.Response(tb.GetObjectType())
	}
}

func (handler *KeyHandler) RandomKey(cli *client.Client) {
	keys := cli.SelectedDatabase().GetAllKeys()
	if len(keys) == 0 {
		cli.Response(nil)
	} else {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		cli.Response(keys[r.Intn(len(keys))])
	}
}

func (handler *KeyHandler) Rename(cli *client.Client) {
	if tb := cli.SelectedDatabase().SearchKeyInDB(cli.Argv[1]); tb == nil {
		cli.ResponseReError(re.ErrNoSuchKey)
	} else {
		cli.SelectedDatabase().RemoveKeyInDB([]string{cli.Argv[1]})
		cli.SelectedDatabase().SetKeyInDB(cli.Argv[2], tb)
		cli.ResponseOK()
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

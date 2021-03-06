package handlers

import (
	"math/rand"
	"strings"
	"time"

	"github.com/SwanSpouse/redis_go/client"
	re "github.com/SwanSpouse/redis_go/error"
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

func (handler *KeyHandler) Del(cli *client.Client) {
	successCount := cli.SelectedDatabase().RemoveKeyInDB(cli.Argv)
	cli.Dirty += successCount
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
		cli.Dirty += 1
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

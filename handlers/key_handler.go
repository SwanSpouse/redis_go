package handlers

import (
	"redis_go/client"
	re "redis_go/error"
	"redis_go/protocol"
)

type KeyHandler struct{}

func (handler *KeyHandler) Process(client *client.Client, command *protocol.Command) {
	switch command.GetName() {
	case "DEL":
		handler.Del(client, command)
	case "DUMP":
	case "EXISTS":
		handler.Exists(client, command)
	case "EXPIRE":
	case "EXPIREAT":
	case "KEYS":
	case "MIGRATE":
	case "MOVE":
	case "OBJECT":
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
	case "TYPE":
	case "SCAN":
	default:
		client.ResponseWriter.AppendErrorf("ERR unknown command %s", command.GetOriginName())
	}
	client.ResponseWriter.Flush()
}

func (handler *KeyHandler) Del(client *client.Client, command *protocol.Command) {
	args := command.GetArgs()
	if len(args) < 2 {
		client.ResponseWriter.AppendErrorf(re.ErrWrongNumberOfArgs, command.GetOriginName())
		return
	}
	successCount := client.SelectedDatabase().RemoveKeyInDB(args[1:])
	client.ResponseWriter.AppendInt(successCount)
}

func (handler *KeyHandler) Exists(client *client.Client, command *protocol.Command) {
	args := command.GetArgs()
	if len(args) < 2 {
		client.ResponseWriter.AppendErrorf(re.ErrWrongNumberOfArgs, command.GetOriginName())
		return
	}
	successCount, _ := client.SelectedDatabase().SearchKeysInDB(args[1:])
	client.ResponseWriter.AppendInt(int64(len(successCount)))
}

package handlers

import (
	"redis_go/client"
	re "redis_go/error"
)

type KeyHandler struct{}

func (handler *KeyHandler) Process(client *client.Client) {
	if client.Cmd == nil {
		client.ResponseError("ERR nil command")
		return
	}
	switch client.Cmd.GetName() {
	case "DEL":
		handler.Del(client)
	case "DUMP":
	case "EXISTS":
		handler.Exists(client)
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
		client.ResponseError("ERR unknown command %s", client.Cmd.GetOriginName())
	}
	client.Flush()
}

func (handler *KeyHandler) Del(client *client.Client) {
	args := client.Cmd.GetArgs()
	if len(args) < 2 {
		client.ResponseError(string(re.ErrWrongNumberOfArgs), client.Cmd.GetOriginName())
		return
	}
	successCount := client.SelectedDatabase().RemoveKeyInDB(args[1:])
	client.Response(successCount)
}

func (handler *KeyHandler) Exists(client *client.Client) {
	args := client.Cmd.GetArgs()
	if len(args) < 2 {
		client.ResponseError(string(re.ErrWrongNumberOfArgs), client.Cmd.GetOriginName())
		return
	}
	successCount, _ := client.SelectedDatabase().SearchKeysInDB(args[1:])
	client.Response(int64(len(successCount)))
}

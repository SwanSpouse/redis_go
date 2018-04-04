package handlers

import (
	"redis_go/client"
	"redis_go/database"
	"redis_go/protocol"
	"redis_go/tcp"
	"strings"
)

type StringHandler struct {
}

func (handler *StringHandler) Process(client *client.Client, command *protocol.Command) {
	switch strings.ToUpper(command.GetName()) {
	case "APPEND":
	case "INCR":
	case "SET":
		handler.Set(client, command)
	case "GET":
		handler.Get(client, command)
	default:
		client.ResponseWriter.AppendErrorf("ERR unknown command %s", command.GetOriginName())
		return
	}
}

func (handler *StringHandler) Set(client *client.Client, command *protocol.Command) {
	args := command.GetArgs()
	if len(args) != 2 {
		client.ResponseWriter.AppendErrorf(tcp.ErrWrongNumberOfArgs, command.GetOriginName())
		client.ResponseWriter.Flush()
		return
	}
	key := args[0]
	value := database.NewRedisObject(args[1])
	client.GetChosenDB().SetKeyInDB(key, value)
	client.ResponseWriter.AppendOK()
}

func (handler *StringHandler) Get(client *client.Client, command *protocol.Command) {
	args := command.GetArgs()
	if len(args) != 1 {
		client.ResponseWriter.AppendErrorf(tcp.ErrWrongNumberOfArgs, command.GetOriginName())
		client.ResponseWriter.Flush()
		return
	}
	key := args[0]
	if obj := client.GetChosenDB().SearchKeyInDB(key); obj != nil {
		client.ResponseWriter.AppendBulkString("")
	} else {
		client.ResponseWriter.AppendNil()
	}
}

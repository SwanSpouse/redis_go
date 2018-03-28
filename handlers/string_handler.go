package handlers

import (
	"redis_go/client"
	"redis_go/protocol"
	"redis_go/redis_database"
	"redis_go/tcp"
	"strings"
)

type StringHandler struct {
}

func (sh *StringHandler) Process(databases []*redis_database.Database, client *client.Client, command *protocol.Command) {
	switch strings.ToUpper(command.GetName()) {
	case "SET":
		sh.Set(databases, client, command)
	case "GET":
		sh.Get(databases, client, command)
	default:
		client.ResponseWriter.AppendErrorf("ERR unknown command %s", command.GetOriginName())
		return
	}
}

func (sh *StringHandler) Set(databases []*redis_database.Database, client *client.Client, command *protocol.Command) {
	args := command.GetArgs()
	if len(args) != 2 {
		client.ResponseWriter.AppendErrorf(tcp.ErrWrongNumberOfArgs, command.GetOriginName())
	}
	key := args[0]
	value := redis_database.NewRedisObject(args[1])
	databases[0].SetKeyInDB(key, value)
	client.ResponseWriter.AppendOK()
}

func (sh *StringHandler) Get(databases []*redis_database.Database, client *client.Client, command *protocol.Command) {
	args := command.GetArgs()
	if len(args) != 1 {
		client.ResponseWriter.AppendErrorf(tcp.ErrWrongNumberOfArgs, command.GetOriginName())
	}
	key := args[0]
	obj := databases[0].SearchKeyInDB(key)
	client.ResponseWriter.AppendBulkString(obj.GetValue())
}

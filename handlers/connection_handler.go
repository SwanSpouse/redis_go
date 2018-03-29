package handlers

import (
	"redis_go/client"
	"redis_go/log"
	"redis_go/protocol"
	"strings"
)

type ConnectionHandler struct {
}

func (handler *ConnectionHandler) Process(client *client.Client, command *protocol.Command) {
	switch strings.ToUpper(command.GetName()) {
	case "PING":
		handler.ping(client, command)
	case "AUTH":
		handler.auth(client, command)
	default:
		client.ResponseWriter.AppendErrorf("ERR unknown command %s", command.GetOriginName())
		return
	}
}

func (handler *ConnectionHandler) ping(client *client.Client, command *protocol.Command) {
	msg := "PONG"
	log.Info("message we send to client %+v", msg)
	client.ResponseWriter.AppendArrayLen(1)
	client.ResponseWriter.AppendBulkString("PONG")
}

func (handler *ConnectionHandler) auth(client *client.Client, command *protocol.Command) {

}

package handlers

import (
	"redis_go/client"
	"redis_go/protocol"
	"strings"
	"redis_go/log"
)

type BaseHandler interface {
	Process(client *client.Client, command *protocol.Command)
}

type ConnectionHandler struct {
}

func (sh *ConnectionHandler) Process(client *client.Client, command *protocol.Command) {
	switch strings.ToUpper(command.GetName()) {
	case "PING":
		sh.Ping(client, command)
	case "AUTH":
		sh.Auth(client, command)
	default:
		client.ResponseWriter.AppendErrorf("ERR unknown command %s", command.GetOriginName())
		return
	}
}

func (sh *ConnectionHandler) Ping(client *client.Client, command *protocol.Command) {
	msg := "PONG"
	log.Info("message we send to client %+v", msg)
	client.ResponseWriter.AppendInlineString("PONG")
}

func (sh *ConnectionHandler) Auth(client *client.Client, command *protocol.Command) {

}

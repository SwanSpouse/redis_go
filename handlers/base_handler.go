package handlers

import (
	"redis_go/client"
	"redis_go/protocol"
	"strings"
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
		return
	}
}

func (sh *ConnectionHandler) Ping(client *client.Client, command *protocol.Command) {

}

func (sh *ConnectionHandler) Auth(client *client.Client, command *protocol.Command) {

}

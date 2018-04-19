package handlers

import (
	"redis_go/client"
	re "redis_go/error"
	"redis_go/log"
	"strings"
)

type ConnectionHandler struct {
}

func (handler *ConnectionHandler) Process(client *client.Client) {
	if client.Cmd == nil {
		client.ResponseReError(re.ErrNilCommand)
		return
	}
	switch strings.ToUpper(client.Cmd.GetName()) {
	case "PING":
		handler.ping(client)
	case "AUTH":
		handler.auth(client)
	default:
		client.ResponseReError(re.ErrUnknownCommand, client.Cmd.GetOriginName())
		return
	}
	client.Flush()
}

func (handler *ConnectionHandler) ping(client *client.Client) {
	msg := "PONG"
	log.Info("message we send to client %+v", msg)
	client.Response(msg)
}

func (handler *ConnectionHandler) auth(client *client.Client) {

}

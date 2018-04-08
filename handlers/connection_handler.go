package handlers

import (
	"redis_go/client"
	"redis_go/log"
	"strings"
)

type ConnectionHandler struct {
}

func (handler *ConnectionHandler) Process(client *client.Client) {
	if client.Cmd == nil {
		client.AppendErrorf("ERR nil command")
		return
	}
	switch strings.ToUpper(client.Cmd.GetName()) {
	case "PING":
		handler.ping(client)
	case "AUTH":
		handler.auth(client)
	default:
		client.AppendErrorf("ERR unknown command %s", client.Cmd.GetOriginName())
		return
	}
}

func (handler *ConnectionHandler) ping(client *client.Client) {
	msg := "PONG"
	log.Info("message we send to client %+v", msg)
	client.AppendInlineString("PONG")
}

func (handler *ConnectionHandler) auth(client *client.Client) {

}

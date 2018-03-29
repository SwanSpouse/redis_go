package handlers

import (
	"redis_go/client"
	"redis_go/protocol"
)

type PubSubHandler struct{}

func (handler *PubSubHandler) Process(client *client.Client, command *protocol.Command) {

}

package handlers

import (
	"redis_go/client"
	"redis_go/protocol"
)

type HashHandler struct {
}

func (handler *HashHandler) Process(client *client.Client, command *protocol.Command) {}

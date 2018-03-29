package handlers

import (
	"redis_go/client"
	"redis_go/protocol"
)

type ListHandler struct {
}

func (handler *ListHandler) Process(client *client.Client, command *protocol.Command) {}

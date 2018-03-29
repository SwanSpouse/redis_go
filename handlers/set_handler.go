package handlers

import (
	"redis_go/protocol"
	"redis_go/client"
)

type SetHandler struct{}

func (handler *SetHandler) Process(client *client.Client, command *protocol.Command) {}

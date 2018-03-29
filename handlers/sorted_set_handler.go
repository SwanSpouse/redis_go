package handlers

import (
	"redis_go/client"
	"redis_go/protocol"
)

type SortedSetHandler struct {
}

func (handler *SortedSetHandler) Process(client *client.Client, command *protocol.Command) {}

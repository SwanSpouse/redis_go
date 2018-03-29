package handlers

import (
	"redis_go/client"
	"redis_go/protocol"
)

type TransactionHandler struct{}

func (handler *TransactionHandler) Process(client *client.Client, command *protocol.Command) {}

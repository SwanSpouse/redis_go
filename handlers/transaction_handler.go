package handlers

import (
	"redis_go/client"
)

type TransactionHandler struct{}

func (handler *TransactionHandler) Process(client *client.Client) {}

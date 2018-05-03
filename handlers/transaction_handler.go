package handlers

import (
	"redis_go/client"
)

var (
	_ client.BaseHandler = (*TransactionHandler)(nil)
)

type TransactionHandler struct{}

func (handler *TransactionHandler) Process(client *client.Client) {}

package handlers

import (
	"redis_go/client"
)

var (
	_ client.BaseHandler = (*SetHandler)(nil)
)

type SetHandler struct{}

func (handler *SetHandler) Process(client *client.Client) {}

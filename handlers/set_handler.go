package handlers

import (
	"redis_go/client"
)

type SetHandler struct{}

func (handler *SetHandler) Process(client *client.Client) {}

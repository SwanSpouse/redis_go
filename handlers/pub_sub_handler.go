package handlers

import (
	"redis_go/client"
)

var (
	_ client.BaseHandler = (*PubSubHandler)(nil)
)

type PubSubHandler struct{}

func (handler *PubSubHandler) Process(client *client.Client) {

}

package handlers

import (
	"github.com/SwanSpouse/redis_go/client"
)

var (
	_ client.BaseHandler = (*PubSubHandler)(nil)
)

type PubSubHandler struct{}

func (handler *PubSubHandler) Process(client *client.Client) {

}

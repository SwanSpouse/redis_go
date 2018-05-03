package handlers

import (
	"redis_go/client"
)

var (
	_ client.BaseHandler = (*SortedSetHandler)(nil)
)

type SortedSetHandler struct {
}

func (handler *SortedSetHandler) Process(client *client.Client) {}

package handlers

import (
	"redis_go/client"
)

type SortedSetHandler struct {
}

func (handler *SortedSetHandler) Process(client *client.Client) {}

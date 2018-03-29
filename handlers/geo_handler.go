package handlers

import (
	"redis_go/protocol"
	"redis_go/client"
)

type GeoHandler struct {
}

func (handler *GeoHandler) Process(client *client.Client, command *protocol.Command) {

}

package handlers

import (
	"redis_go/client"
	"redis_go/protocol"
)

type KeyHandler struct {
}

func (handler *KeyHandler) Process(client *client.Client, command *protocol.Command) {

}

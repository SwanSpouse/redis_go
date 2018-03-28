package handlers

import (
	"redis_go/client"
	"redis_go/protocol"
)

var (
	_ BaseHandler = (*ConnectionHandler)(nil)
	_ BaseHandler = (*StringHandler)(nil)
)

type BaseHandler interface {
	Process(*client.Client, *protocol.Command)
}

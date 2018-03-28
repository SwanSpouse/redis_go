package handlers

import (
	"redis_go/client"
	"redis_go/protocol"
	"redis_go/redis_database"
)

var (
	_ BaseHandler = (*ConnectionHandler)(nil)
	_ BaseHandler = (*StringHandler)(nil)
)

type BaseHandler interface {
	Process([]*redis_database.Database, *client.Client, *protocol.Command)
}

package handlers

import (
	"redis_go/client"
)

var (
	_ BaseHandler = (*KeyHandler)(nil)
	_ BaseHandler = (*StringHandler)(nil)
	_ BaseHandler = (*HashHandler)(nil)
	_ BaseHandler = (*ListHandler)(nil)
	_ BaseHandler = (*SetHandler)(nil)
	_ BaseHandler = (*SortedSetHandler)(nil)
	_ BaseHandler = (*GeoHandler)(nil)
	_ BaseHandler = (*PubSubHandler)(nil)
	_ BaseHandler = (*TransactionHandler)(nil)
	_ BaseHandler = (*ConnectionHandler)(nil)
)

type BaseHandler interface {
	Process(*client.Client)
}

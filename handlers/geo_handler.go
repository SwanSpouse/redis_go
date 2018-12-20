package handlers

import (
	"github.com/SwanSpouse/redis_go/client"
)

var (
	_ client.BaseHandler = (*GeoHandler)(nil)
)

type GeoHandler struct {
}

func (handler *GeoHandler) Process(client *client.Client) {

}

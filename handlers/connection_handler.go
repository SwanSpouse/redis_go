package handlers

import (
	"redis_go/client"
	re "redis_go/error"
	"redis_go/loggers"
)

var (
	_ client.BaseHandler = (*ConnectionHandler)(nil)
)

const (
	RedisConnectionCommandPing = "PING"
	RedisConnectionCommandAuth = "AUTH"
)

type ConnectionHandler struct {
}

func (handler *ConnectionHandler) Process(cli *client.Client) {
	switch cli.GetCommandName() {
	case RedisConnectionCommandPing:
		handler.ping(cli)
	case RedisConnectionCommandAuth:
		handler.auth(cli)
	default:
		cli.ResponseReError(re.ErrUnknownCommand, cli.GetOriginCommandName())
		return
	}
	cli.Flush()
}

func (handler *ConnectionHandler) ping(client *client.Client) {
	msg := "PONG"
	loggers.Info("message we send to client %+v", msg)
	client.Response(msg)
}

func (handler *ConnectionHandler) auth(client *client.Client) {

}

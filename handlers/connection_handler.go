package handlers

import (
	"github.com/SwanSpouse/redis_go/client"
		"github.com/SwanSpouse/redis_go/loggers"
)

const (
	RedisConnectionCommandPing   = "PING"
	RedisConnectionCommandAuth   = "AUTH"
	RedisConnectionCommandSelect = "SELECT"
	RedisConnectionCommandEcho   = "ECHO"
	RedisConnectionCommandQuit   = "QUIT"
)

type ConnectionHandler struct {
}

func (handler *ConnectionHandler) Ping(cli *client.Client) {
	msg := "PONG"
	loggers.Info("message we send to cli %+v", msg)
	cli.Response(msg)
}

func (handler *ConnectionHandler) Auth(cli *client.Client) {

}

func (handler *ConnectionHandler) Echo(cli *client.Client) {
	cli.Response(cli.Argv[1])
}

func (handler *ConnectionHandler) Quit(cli *client.Client) {

}

func (handler *ConnectionHandler) CmdSelect(cli *client.Client) {
	cli.ResponseOK()
}

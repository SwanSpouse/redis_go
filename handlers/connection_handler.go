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
	RedisConnectionCommandPing   = "PING"
	RedisConnectionCommandAuth   = "AUTH"
	RedisConnectionCommandSelect = "SELECT"
	RedisConnectionCommandEcho   = "ECHO"
	RedisConnectionCommandQuit   = "QUIT"
)

type ConnectionHandler struct {
}

func (handler *ConnectionHandler) Process(cli *client.Client) {
	switch cli.Cmd.GetName() {
	case RedisConnectionCommandPing:
		handler.ping(cli)
	case RedisConnectionCommandAuth:
		handler.auth(cli)
	case RedisConnectionCommandEcho:
		handler.echo(cli)
	case RedisConnectionCommandQuit:
		handler.quit(cli)
	case RedisConnectionCommandSelect:
		handler.cmdSelect(cli)
	default:
		cli.ResponseReError(re.ErrUnknownCommand, cli.Cmd.GetOriginName())
		return
	}
	cli.Flush()
}

func (handler *ConnectionHandler) ping(cli *client.Client) {
	msg := "PONG"
	loggers.Info("message we send to cli %+v", msg)
	cli.Response(msg)
}

func (handler *ConnectionHandler) auth(cli *client.Client) {

}

func (handler *ConnectionHandler) echo(cli *client.Client) {
	cli.Response(cli.Argv[1])
}

func (handler *ConnectionHandler) quit(cli *client.Client) {

}

func (handler *ConnectionHandler) cmdSelect(cli *client.Client) {
	cli.ResponseOK()
}

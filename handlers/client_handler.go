package handlers

import (
	"strings"

	"github.com/SwanSpouse/redis_go/client"
	re "github.com/SwanSpouse/redis_go/error"
)

const (
	RedisClientCommand = "CLIENT"

	RedisClientSubCommandGetName = "GETNAME"
	RedisClientSubCommandSetName = "SETNAME"
	RedisClientSubCommandList    = "LIST"
	RedisClientSubCommandKill    = "KILL"
	RedisClientSubCommandReply   = "REPLY"
	RedisClientSubCommandPause   = "PAUSE"
)

type ClientHandler struct {
}

func (handler *ClientHandler) Client(cli *client.Client) {
	switch strings.ToUpper(cli.Argv[1]) {
	case RedisClientSubCommandGetName:
		handler.getName(cli)
	case RedisClientSubCommandSetName:
		handler.setName(cli)
	case RedisClientSubCommandList:
		handler.list(cli)
	case RedisClientSubCommandKill:
		handler.kill(cli)
	case RedisClientSubCommandReply:
		handler.reply(cli)
	case RedisClientSubCommandPause:
		handler.pause(cli)
	default:
		cli.ResponseReError(re.ErrClientCommand)
	}
	cli.Flush()
}

func (handler *ClientHandler) getName(cli *client.Client) {
	cli.Response(cli.Name)
}

func (handler *ClientHandler) setName(cli *client.Client) {
	if len(cli.Argv) != 2 || len(cli.Argv[1]) == 0 {
		cli.ResponseReError(re.ErrClientCommand)
	} else {
		cli.Name = cli.Argv[1]
		cli.ResponseOK()
	}
}

func (handler *ClientHandler) list(cli *client.Client) {
	panic("not implement")
}

func (handler *ClientHandler) kill(cli *client.Client) {
	panic("not implement")
}

func (handler *ClientHandler) reply(cli *client.Client) {
	panic("not implement")
}

func (handler *ClientHandler) pause(cli *client.Client) {
	panic("not implement")
}

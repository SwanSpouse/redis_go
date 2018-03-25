package redis_command

import "redis_go/resp"

type ConnectionCommand interface {
	AUTH(c *networking.Client)
	ECHO(c *networking.Client)
	PING(c *networking.Client)

	QUIT(c *networking.Client)
	SELECT(c *networking.Client)
}

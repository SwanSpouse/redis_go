package redis_command

import "redis_go/resp"

type ConnectionCommand interface {
	AUTH(c *resp.Client)
	ECHO(c *resp.Client)
	PING(c *resp.Client)

	QUIT(c *resp.Client)
	SELECT(c *resp.Client)
}

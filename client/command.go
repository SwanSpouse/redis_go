package client

import (
	"fmt"
	"strings"
)

const (
	/* Command Flags. Please check the command table defined in the server_command_table.go file
	* for more information about the meaning of every flag. */
	RedisCmdWrite            = 1    /* "w" flag */
	RedisCmdReadOnly         = 2    /* "r" flag */
	RedisCmdDenyOom          = 4    /* "m" flag */
	RedisCmdForceReplication = 8    /* "f" flag */
	RedisCmdAdmin            = 16   /* "a" flag */
	RedisCmdPubSub           = 32   /* "p" flag */
	RedisCmdNoScript         = 64   /* "s" flag */
	RedisCMdRandom           = 128  /* "R" flag */
	RedisCmdSortForScript    = 256  /* "S" flag */
	RedisCmdLoading          = 512  /* "l" flag */
	RedisCmdStable           = 1024 /* "t" flag */
	RedisCmdSkipMonitor      = 2048 /* "M" flag */
)

type BaseHandler interface {
	Process(*Client)
}

type Command struct {
	name        string      // command name
	Arity       int         // command args
	SFlags      string      //
	Flags       int         //
	microsecond int64       // execute time in microsecond
	calls       int64       // call times
	Handler     BaseHandler // redis命令的对应handler
}

func NewCommand(name string, arity int, sflags string, proc BaseHandler) *Command {
	return &Command{
		name:    name,
		Arity:   arity,
		SFlags:  sflags,
		Handler: proc,
	}
}

func (c *Command) GetMicrosecond() int64 {
	return c.microsecond
}

func (c *Command) GetCalls() int64 {
	return c.calls
}

func (c *Command) String() string {
	return fmt.Sprintf("current command args info: command name: %s. args:", c.name)
}

func (c *Command) GetName() string {
	return strings.ToUpper(c.name)
}

func (c *Command) GetOriginName() string {
	return strings.ToUpper(c.name)
}

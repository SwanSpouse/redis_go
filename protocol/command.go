package protocol

import (
	"redis_go/tcp"
	"strings"
)

type Command struct {
	name        string                // command name
	arity       int                   // command args
	sflags      rune                  //
	flags       int                   //
	args        []tcp.CommandArgument // command args
	microsecond int64                 // execute time in microsecond
	calls       int64                 // call times
}

func NewCommand() *Command {
	return &Command{args: make([]tcp.CommandArgument, 0)}
}

func (c *Command) GetArgs() []tcp.CommandArgument {
	return c.args
}

func (c *Command) AddArgs(arg tcp.CommandArgument) {
	if c.args == nil {
		c.args = make([]tcp.CommandArgument, 0)
	}
	c.args = append(c.args, arg)
}

func (c *Command) GetName() string {
	return strings.ToUpper(c.name)
}

func (c *Command) GetOriginName() string {
	return c.name
}

func (c *Command) SetName(name string) {
	c.name = name
}

func (c *Command) GetMicrosecond() int64 {
	return c.microsecond
}

func (c *Command) GetCalls() int64 {
	return c.calls
}

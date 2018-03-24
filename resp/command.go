package resp

import (
	"strings"
)

// CommandArgument is an argument of a command
type CommandArgument []byte

type Command struct {
	name             string               // command name
	redisCommandProc func(client *Client) // function deal with the command
	arity            int                  // command args
	sflags           rune                 //
	flags            int                  //
	args             []CommandArgument    // command args
	microsecond      int64                // execute time in microsecond
	calls            int64                // call times
}

func NewCommand() *Command {
	return &Command{args: make([]CommandArgument, 0)}
}

func (c *Command) GetArgs() []CommandArgument {
	return c.args
}

func (c *Command) AddArgs(arg CommandArgument) {
	if c.args == nil {
		c.args = make([]CommandArgument, 0)
	}
	c.args = append(c.args, arg)
}

func (c *Command) GetName() string {
	return c.name
}

func (c *Command) SetName(name string) {
	c.name = strings.ToLower(name)
}

func (c *Command) GetMicrosecond() int64 {
	return c.microsecond
}

func (c *Command) GetCalls() int64 {
	return c.calls
}
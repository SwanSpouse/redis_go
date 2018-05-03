package client

import (
	"fmt"
)

type BaseHandler interface {
	Process(*Client)
}

type Command struct {
	name        string      // command name
	arity       int         // command args
	Argc        int         // argument counter
	Argv        []string    // argument vector
	sflags      string      //
	flags       int         //
	microsecond int64       // execute time in microsecond
	calls       int64       // call times
	Handler     BaseHandler // redis命令的对应handler
}

func NewCommand(name string, arity int, sflags string, proc BaseHandler) *Command {
	return &Command{
		name:    name,
		arity:   arity,
		sflags:  sflags,
		Argv:    make([]string, 0),
		Argc:    0,
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
	ret := fmt.Sprintf("current command args info: command name: %s. args:", c.name)
	for _, item := range c.Argv {
		ret += fmt.Sprintf("%s ", string(item))
	}
	return ret
}

package resp

import "context"

// CommandArgument is an argument of a command
type CommandArgument []byte

// --------------------------------------------------------------------

// Command instances are parsed by a RequestReader
type Command struct {
	// Name refers to the command name
	Name string

	// Args returns arguments
	Args []CommandArgument

	ctx context.Context
}

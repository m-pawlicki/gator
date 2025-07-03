package commands

import (
	"fmt"

	"github.com/m-pawlicki/gator/internal/state"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	commands map[string]func(*state.State, Command) error
}

func NewCommands() *Commands {
	return &Commands{
		commands: make(map[string]func(*state.State, Command) error),
	}
}

func (c *Commands) Run(s *state.State, cmd Command) error {
	val, ok := c.commands[cmd.Name]
	if ok {
		val(s, cmd)
	}
	return fmt.Errorf("command doesn't exist")
}

func (c *Commands) Register(name string, f func(*state.State, Command) error) {
	c.commands[name] = f
}

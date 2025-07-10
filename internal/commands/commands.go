package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/m-pawlicki/gator/internal/database"
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

func MiddlewareLoggedIn(handler func(s *state.State, cmd Command, user database.User) error) func(*state.State, Command) error {
	login := func(s *state.State, cmd Command) error {
		user, err := s.DB.GetUser(context.Background(), s.Config.User)
		if err != nil {
			fmt.Println("User not logged in.")
			os.Exit(1)
		}
		return handler(s, cmd, user)
	}
	return login
}

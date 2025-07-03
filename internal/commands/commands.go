package commands

import (
	"fmt"

	"github.com/m-pawlicki/gator/internal/state"
)

type Command struct {
	name string
	args []string
}

func handlerLogin(s *state.State, cmd Command) error {
	if len(cmd.args) < 1 {

	}
	s.Cfg.SetUser(cmd.args[0])
	fmt.Println("User has been set.")
	return nil
}

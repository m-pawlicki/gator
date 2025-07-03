package handlers

import (
	"fmt"
	"os"

	"github.com/m-pawlicki/gator/internal/commands"
	"github.com/m-pawlicki/gator/internal/state"
)

func HandlerLogin(s *state.State, cmd commands.Command) error {
	if len(cmd.Args) < 1 {
		fmt.Println("Error: Username required.")
		os.Exit(1)
	}
	s.Config.SetUser(cmd.Args[0])
	fmt.Printf("User %s has been set.\n", cmd.Args[0])
	return nil
}

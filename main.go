package main

import _ "github.com/lib/pq"

import (
	"fmt"
	"os"

	"github.com/m-pawlicki/gator/internal/commands"
	"github.com/m-pawlicki/gator/internal/config"
	"github.com/m-pawlicki/gator/internal/handlers"
	"github.com/m-pawlicki/gator/internal/state"
)

func main() {
	cfg := config.Read()
	st := state.NewState(&cfg)
	cmds := commands.NewCommands()
	cmds.Register("login", handlers.HandlerLogin)
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Error: Not enough arguments provided.")
		os.Exit(1)
	}
	cmd := commands.Command{Name: args[1], Args: args[2:]}
	cmds.Run(st, cmd)
}

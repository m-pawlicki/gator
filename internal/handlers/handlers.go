package handlers

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/m-pawlicki/gator/internal/commands"
	"github.com/m-pawlicki/gator/internal/database"
	"github.com/m-pawlicki/gator/internal/state"
)

func HandlerLogin(s *state.State, cmd commands.Command) error {
	if len(cmd.Args) < 1 {
		fmt.Println("Error: Username required.")
		os.Exit(1)
	}
	ctx := context.Background()
	_, err := s.DB.GetUser(ctx, cmd.Args[0])
	if err != nil {
		fmt.Println("Error: User doesn't exist.")
		os.Exit(1)
	} else {
		s.Config.SetUser(cmd.Args[0])
		fmt.Printf("User %s has been set.\n", cmd.Args[0])
	}
	return nil
}

func HandlerRegister(s *state.State, cmd commands.Command) error {
	if len(cmd.Args) < 1 {
		fmt.Println("Error: Username required.")
		os.Exit(1)
	}
	ctx := context.Background()
	username := cmd.Args[0]
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}
	dbParams := database.CreateUserParams(params)
	usr, _ := s.DB.GetUser(ctx, username)
	if usr.Name == username {
		fmt.Println("Error: User already exists.")
		os.Exit(1)
	} else {
		s.DB.CreateUser(ctx, dbParams)
		s.Config.SetUser(username)
		fmt.Printf("User created: %v\n", params)
	}
	return nil
}

func HandlerReset(s *state.State, cmd commands.Command) error {
	ctx := context.Background()
	err := s.DB.ResetUsers(ctx)
	if err != nil {
		fmt.Println("Error: Failed to reset users.")
		os.Exit(1)
	}
	fmt.Println("Users reset.")
	return nil
}

func HandlerUsers(s *state.State, cmd commands.Command) error {
	ctx := context.Background()
	users, err := s.DB.GetUsers(ctx)
	if err != nil {
		fmt.Println("Error: Couldn't retrieve users.")
		os.Exit(1)
	}
	for _, usr := range users {
		if usr == s.Config.User {
			fmt.Printf("* %s (current)\n", usr)
		} else {
			fmt.Printf("* %s\n", usr)
		}
	}
	return nil
}

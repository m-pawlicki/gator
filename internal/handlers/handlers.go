package handlers

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/m-pawlicki/gator/internal/commands"
	"github.com/m-pawlicki/gator/internal/database"
	"github.com/m-pawlicki/gator/internal/rss"
	"github.com/m-pawlicki/gator/internal/state"
)

func HandlerLogin(s *state.State, cmd commands.Command) error {
	if len(cmd.Args) < 1 {
		fmt.Println("Username required.")
		fmt.Println("Usage: login <username>")
		os.Exit(1)
	}
	ctx := context.Background()
	_, err := s.DB.GetUser(ctx, cmd.Args[0])
	if err != nil {
		fmt.Println("User doesn't exist.")
		os.Exit(1)
	} else {
		s.Config.SetUser(cmd.Args[0])
		fmt.Printf("User %s has been set.\n", cmd.Args[0])
	}
	return nil
}

func HandlerRegister(s *state.State, cmd commands.Command) error {
	if len(cmd.Args) < 1 {
		fmt.Println("Username required.")
		fmt.Println("Usage: register <username>")
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
		fmt.Println("User already exists.")
		os.Exit(1)
	} else {
		s.DB.CreateUser(ctx, dbParams)
		s.Config.SetUser(username)
		fmt.Printf("User created: %s\n", params)
	}
	return nil
}

func HandlerReset(s *state.State, cmd commands.Command) error {
	ctx := context.Background()
	err := s.DB.ResetUsers(ctx)
	if err != nil {
		fmt.Println("Failed to reset users.")
		os.Exit(1)
	}
	fmt.Println("Users reset.")
	return nil
}

func HandlerUsers(s *state.State, cmd commands.Command) error {
	ctx := context.Background()
	users, err := s.DB.GetUsers(ctx)
	if err != nil {
		fmt.Println("Couldn't retrieve users.")
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

func HandlerAgg(s *state.State, cmd commands.Command) error {
	ctx := context.Background()
	feedURL := "https://www.wagslane.dev/index.xml"
	feed, err := rss.FetchFeed(ctx, feedURL)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Printf("%s\n", feed)
	return nil
}

func HandlerAddFeed(s *state.State, cmd commands.Command) error {
	if len(cmd.Args) < 1 {
		fmt.Println("Title and URL missing.")
		fmt.Println("Usage: addfeed <title> <url>")
		os.Exit(1)
	}
	if len(cmd.Args) < 2 {
		fmt.Println("URL missing.")
		fmt.Println("Usage: addfeed <title> <url>")
		os.Exit(1)
	}
	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]
	ctx := context.Background()
	currUser := s.Config.User
	userID, err := s.DB.GetUser(ctx, currUser)
	if err != nil {
		fmt.Println("Error retrieving user.")
		os.Exit(1)
	}
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    userID.ID,
	}
	_, err = s.DB.CreateFeed(ctx, params)
	if err != nil {
		fmt.Println("Couldn't add feed.")
		os.Exit(1)
	} else {
		fmt.Printf("Feed added: %s\n", params)
	}
	return nil
}

func HandlerFeeds(s *state.State, cmd commands.Command) error {
	ctx := context.Background()
	feeds, err := s.DB.GetFeeds(ctx)
	if err != nil {
		fmt.Println("Couldn't get feeds.")
		os.Exit(1)
	}
	if len(feeds) < 1 {
		fmt.Println("No feeds available.")
		return nil
	}
	for _, feed := range feeds {
		username, err := s.DB.GetUserFromID(ctx, feed.UserID)
		if err != nil {
			fmt.Println("Couldn't get user.")
			os.Exit(1)
		}
		fmt.Printf("Feed name: %s | Feed URL: %s | Added by: %s\n", feed.Name, feed.Url, username)
	}
	return nil
}

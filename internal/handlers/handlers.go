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
		newUsr, err := s.DB.CreateUser(ctx, dbParams)
		if err != nil {
			fmt.Println("Error creating user.")
			os.Exit(1)
		}
		s.Config.SetUser(username)
		fmt.Printf("User %s created at %s\n", newUsr.Name, newUsr.CreatedAt)
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
	fmt.Println("List of all users:")
	for _, usr := range users {
		if usr.Name == s.Config.User {
			fmt.Printf("* %s (current)\n", usr.Name)
		} else {
			fmt.Printf("* %s\n", usr.Name)
		}
	}
	return nil
}

func HandlerAgg(s *state.State, cmd commands.Command) error {
	if len(cmd.Args) < 1 {
		fmt.Println("Please enter a duration.")
		fmt.Println("Usage example: agg <1s/1m/1hr>")
		os.Exit(1)
	}
	duration, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		fmt.Println("Couldn't parse duration.")
		os.Exit(1)
	}
	fmt.Println("Collecting feeds every", duration)
	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		rss.ScrapeFeeds(s)
	}
}

func HandlerAddFeed(s *state.State, cmd commands.Command, user database.User) error {
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
	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	}
	newFeed, err := s.DB.CreateFeed(ctx, params)
	if err != nil {
		fmt.Println("Couldn't add feed.")
		os.Exit(1)
	} else {
		fmt.Printf("Feed added: %s - %s\n", newFeed.Name, newFeed.Url)
	}
	fefoParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: newFeed.CreatedAt,
		UpdatedAt: newFeed.UpdatedAt,
		UserID:    newFeed.UserID,
		FeedID:    newFeed.ID,
	}
	fefo, err := s.DB.CreateFeedFollow(ctx, fefoParams)
	if err != nil {
		fmt.Println("Failed to follow feed after creation.")
		os.Exit(1)
	} else {
		fmt.Printf("%s is now following the feed '%s'", fefo.UserName, fefo.FeedName)
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

func HandlerFollow(s *state.State, cmd commands.Command, user database.User) error {
	if len(cmd.Args) < 1 {
		fmt.Println("URL missing.")
		fmt.Println("Usage: follow <url>")
		os.Exit(1)
	}
	ctx := context.Background()
	url := cmd.Args[0]
	feed, err := s.DB.GetFeedByURL(ctx, url)
	if err != nil {
		fmt.Println("Couldn't find feed by URL.")
		os.Exit(1)
	}
	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	}
	fefo, err := s.DB.CreateFeedFollow(ctx, params)
	if err != nil {
		fmt.Println("Couldn't follow feed.")
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Printf("%s is now following the feed: %s", fefo.UserName, fefo.FeedName)
	}
	return nil
}

func HandlerFollowing(s *state.State, cmd commands.Command, user database.User) error {
	ctx := context.Background()
	feeds, err := s.DB.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		fmt.Println("Failed to get feeds for user.")
		os.Exit(1)
	}
	fmt.Printf("Feeds %s is currently following:\n", user.Name)
	for _, feed := range feeds {
		fmt.Printf("- %s\n", feed.FeedName)
	}
	return nil
}

func HandlerUnfollow(s *state.State, cmd commands.Command, user database.User) error {
	if len(cmd.Args) < 1 {
		fmt.Println("URL missing.")
		fmt.Println("Usage: unfollow <url>")
		os.Exit(1)
	}
	ctx := context.Background()
	params := database.DeleteFeedFollowParams{
		UserID: user.ID,
		Url:    cmd.Args[0],
	}
	err := s.DB.DeleteFeedFollow(ctx, params)
	if err != nil {
		fmt.Println("Error unfollowing feed.")
		os.Exit(1)
	}
	fmt.Println("Unfollowed feed successfully.")
	return nil
}

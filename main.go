package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/m-pawlicki/gator/internal/commands"
	"github.com/m-pawlicki/gator/internal/config"
	"github.com/m-pawlicki/gator/internal/database"
	"github.com/m-pawlicki/gator/internal/handlers"
	"github.com/m-pawlicki/gator/internal/state"
)

func main() {
	cfg := config.Read()
	st := state.NewState(&cfg)
	dbURL := cfg.DB
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Error: Couldn't open database.")
		os.Exit(1)
	}
	dbQueries := database.New(db)
	st.DB = dbQueries
	cmds := commands.NewCommands()
	cmds.Register("login", handlers.HandlerLogin)
	cmds.Register("register", handlers.HandlerRegister)
	cmds.Register("reset", handlers.HandlerReset)
	cmds.Register("users", handlers.HandlerUsers)
	cmds.Register("agg", handlers.HandlerAgg)
	cmds.Register("feeds", handlers.HandlerFeeds)
	cmds.Register("addfeed", commands.MiddlewareLoggedIn(handlers.HandlerAddFeed))
	cmds.Register("follow", commands.MiddlewareLoggedIn(handlers.HandlerFollow))
	cmds.Register("following", commands.MiddlewareLoggedIn(handlers.HandlerFollowing))
	cmds.Register("unfollow", commands.MiddlewareLoggedIn(handlers.HandlerUnfollow))
	cmds.Register("browse", commands.MiddlewareLoggedIn(handlers.HandlerBrowse))
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Error: Not enough arguments provided.")
		os.Exit(1)
	}
	cmd := commands.Command{Name: args[1], Args: args[2:]}
	cmds.Run(st, cmd)
}

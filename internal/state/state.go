package state

import (
	"github.com/m-pawlicki/gator/internal/config"
	"github.com/m-pawlicki/gator/internal/database"
)

type State struct {
	DB     *database.Queries
	Config *config.Config
}

func NewState(cfg *config.Config) *State {
	return &State{
		Config: cfg,
	}
}

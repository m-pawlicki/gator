package state

import (
	"github.com/m-pawlicki/gator/internal/config"
)

type State struct {
	Config *config.Config
}

func NewState(cfg *config.Config) *State {
	return &State{
		Config: cfg,
	}
}

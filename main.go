package main

import (
	"fmt"

	"github.com/m-pawlicki/gator/internal/config"
)

func main() {
	cfg := config.Read()
	cfg.SetUser("micah")
	cfg = config.Read()
	fmt.Println(cfg)
}

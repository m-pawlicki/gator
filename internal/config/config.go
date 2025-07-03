package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DB   string `json:"db_url"`
	User string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() Config {
	var c Config
	path, err := getConfigFilePath()
	if err != nil {
		fmt.Println(err)
	}
	jsonBytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(jsonBytes, &c)
	if err != nil {
		fmt.Println(err)
	}
	return c
}

func (cfg Config) SetUser(user string) {
	cfg.User = user
	write(cfg)
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(home, configFileName)
	return path, nil
}

func write(cfg Config) error {
	jsonBytes, err := json.Marshal(cfg)
	if err != nil {
		fmt.Println(err)
	}
	path, err := getConfigFilePath()
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile(path, jsonBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

package config

import (
	"fmt"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type LaunchMode string

const (
	LocalEnv LaunchMode = "local"
	ProdEnv  LaunchMode = "prod"
)

type Config struct {
	Port string `envconfig:"PORT" default:"8080"`
}

func Load(launchMode LaunchMode, path string) (*Config, error) {
	switch launchMode {
	case LocalEnv:
		cfgPath := filepath.Join(path, fmt.Sprintf("%s.env", launchMode))
		err := godotenv.Load(cfgPath)
		if err != nil {
			return nil, fmt.Errorf("load .env config file: %w", err)
		}
	case ProdEnv:
	// all settings should be provided as env variables
	default:
		return nil, fmt.Errorf("unexpected LAUNCH_MODE: [%s]", launchMode)
	}
	config := new(Config)
	err := envconfig.Process("", config)
	if err != nil {
		return nil, fmt.Errorf("get config from env: %w", err)
	}
	return config, nil
}

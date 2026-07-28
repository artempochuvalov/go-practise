package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port        string
	DatabaseURL string
}

func Load() (Config, error) {
	var config Config

	envDbUrl := os.Getenv("DATABASE_URL")
	if envDbUrl == "" {
		return config, errors.New("no db url provided")
	}

	envPort := os.Getenv("PORT")
	if envPort == "" {
		envPort = "8080"
	}

	_, err := strconv.Atoi(envPort)
	if err != nil {
		return config, err
	}

	config.DatabaseURL = envDbUrl
	config.Port = fmt.Sprintf(":%s", envPort)

	return config, nil
}

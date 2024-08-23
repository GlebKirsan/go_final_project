package config

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port          string `envconfig:"TODO_PORT" default:"7540"`
	DBFile        string `envconfig:"TODO_DBFILE" default:"scheduler.db"`
	LogLevel      string `envconfig:"TODO_LOG_LEVEL"`
	MigrationPath string `envconfig:"TODO_MIGRATION_PATH" default:"internal/database/migration"`
	Pass          string `envconfig:"TODO_PASSWORD" default:"dummy"`
	Secret        string `envconfig:"TODO_TOKEN_SECRET" default:"secret"`
}

var (
	config Config
	once   sync.Once
)

func Get() *Config {
	once.Do(func() {
		err := envconfig.Process("", &config)
		if err != nil {
			log.Fatal(err)
		}
		configBytes, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Configuration:", string(configBytes))
	})

	return &config
}

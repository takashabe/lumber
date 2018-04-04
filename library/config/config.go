package config

import (
	"log"
	"os"

	"github.com/jinzhu/configor"
)

// Config represent configuration
var Config = struct {
	DB struct {
		Name     string `env:"DB_NAME"`
		Host     string `env:"DB_HOST"`
		User     string `env:"DB_USER"`
		Password string `env:"DB_PASSWORD"`
		Port     int    `env:"DB_PORT"`
	}

	Server struct {
		Port int `default:"8080"`
	}
}{}

func init() {
	conf := os.Getenv("APP_CONF")
	err := configor.Load(&Config, conf)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
}

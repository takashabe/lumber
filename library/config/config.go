package config

import (
	"log"
	"os"

	"github.com/jinzhu/configor"
)

// Config represent configuration
var Config = struct {
	DB struct {
		Name     string `env:"LUMBER_DB_NAME"`
		Host     string `env:"LUMBER_DB_HOST"`
		User     string `env:"LUMBER_DB_USER"`
		Password string `env:"LUMBER_DB_PASSWORD"`
		Port     int    `env:"LUMBER_DB_PORT"`
	}

	Server struct {
		Port int `default:"8080" env:"LUMBER_SERVER_PORT"`
	}
}{}

func init() {
	conf := os.Getenv("APP_CONF")
	err := configor.Load(&Config, conf)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
}

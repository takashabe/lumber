package utils

import (
	"database/sql"
	"fmt"
	"os"
)

func getEnvWithDefault(name, def string) string {
	if env := os.Getenv(name); len(env) != 0 {
		return env
	}
	return def
}

// ConnectMySQL returns sql.DB for mysql
func ConnectMySQL() (*sql.DB, error) {
	var (
		user     = getEnvWithDefault("DB_USER", "root")
		password = getEnvWithDefault("DB_PASSWORD", "")
		host     = getEnvWithDefault("DB_HOST", "localhost")
		port     = getEnvWithDefault("DB_PORT", "3306")
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/lumber?parseTime=true", user, password, host, port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(32)
	db.SetMaxOpenConns(32)

	return db, nil
}

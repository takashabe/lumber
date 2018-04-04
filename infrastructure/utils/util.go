package utils

import (
	"database/sql"
	"fmt"

	"github.com/takashabe/lumber/library/config"
)

// ConnectMySQL returns sql.DB for mysql
func ConnectMySQL() (*sql.DB, error) {
	conf := config.Config.DB

	var (
		name     = conf.Name
		user     = conf.User
		password = conf.Password
		host     = conf.Host
		port     = conf.Port
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", user, password, host, port, name)
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

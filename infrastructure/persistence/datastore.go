package persistence

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	"github.com/takashabe/lumber/domain/repository"
)

// SQLRepositoryAdapter provides accessors to the RDB
type SQLRepositoryAdapter struct {
	Conn *sql.DB
}

// NewEntryRepository returns initialized Datastore
func NewEntryRepository() (repository.EntryRepository, error) {
	getEnvWithDefault := func(name, def string) string {
		if env := os.Getenv(name); len(env) != 0 {
			return env
		}
		return def
	}

	user := getEnvWithDefault("DB_USER", "root")
	password := getEnvWithDefault("DB_PASSWORD", "")
	host := getEnvWithDefault("DB_HOST", "localhost")
	port := getEnvWithDefault("DB_PORT", "3306")
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

	return &EntryRepositoryImpl{
		&SQLRepositoryAdapter{Conn: db},
	}, nil
}

func (a *SQLRepositoryAdapter) query(q string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := a.Conn.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Query(args...)
}

func (a *SQLRepositoryAdapter) queryRow(q string, args ...interface{}) (*sql.Row, error) {
	stmt, err := a.Conn.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryRow(args...), nil
}

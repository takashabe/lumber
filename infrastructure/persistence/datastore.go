package persistence

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql" // mysql driver
)

// Datastore represent MySQL adapter
type Datastore struct {
	Conn *sql.DB
}

// NewDatastore returns initialized Datastore
func NewDatastore() (*Datastore, error) {
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

	return &Datastore{Conn: db}, nil
}

// Close calls DB.Close
func (d *Datastore) Close() error {
	return d.Conn.Close()
}

func (d *Datastore) query(q string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := d.Conn.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Query(args...)
}

func (d *Datastore) queryRow(q string, args ...interface{}) (*sql.Row, error) {
	stmt, err := d.Conn.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryRow(args...), nil
}

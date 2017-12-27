package persistence

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // mysql driver
)

// SQLRepositoryAdapter provides accessors to the RDB
type SQLRepositoryAdapter struct {
	Conn *sql.DB
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

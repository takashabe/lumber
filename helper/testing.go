package helper

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/takashabe/go-fixture"
	_ "github.com/takashabe/go-fixture/mysql" // driver
)

// LoadFixture load fixture files
func LoadFixture(t *testing.T, file string) {
	db, err := newDatastore()
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	defer db.Close()

	f, err := fixture.NewFixture(db, "mysql")
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	err = f.Load(file)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
}

// LoadFixtureSQL load sql fixture files
func LoadFixtureSQL(t *testing.T, file string) {
	db, err := newDatastore()
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	defer db.Close()

	f, err := fixture.NewFixture(db, "mysql")
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	err = f.LoadSQL(file)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
}

// SetupTables initialize the database by fixture of the schema
func SetupTables() {
	db, err := newDatastore()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	f, err := fixture.NewFixture(db, "mysql")
	if err != nil {
		panic(err)
	}
	err = f.LoadSQL("../_sql/schema.sql")
	if err != nil {
		panic(err)
	}
}

// newDatastore returns sql.DB
// Porting from datastore package
func newDatastore() (*sql.DB, error) {
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

	return db, nil
}

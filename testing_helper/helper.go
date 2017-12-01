package helper

import (
	"testing"

	"github.com/takashabe/go-fixture"
	_ "github.com/takashabe/go-fixture/mysql" // driver
	"github.com/takashabe/lumber/datastore"
)

// LoadFixture load fixture files
func LoadFixture(t *testing.T, file string) {
	db, err := datastore.NewDatastore()
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	f := fixture.NewFixture(db.Conn, "mysql")
	err = f.Load(file)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
}

// SetupTables initialize the database by fixture of the schema
func SetupTables() {
	db, err := datastore.NewDatastore()
	if err != nil {
		panic(err)
	}

	f := fixture.NewFixture(db.Conn, "mysql")
	err = f.LoadSQL("../_sql/schema.sql")
	if err != nil {
		panic(err)
	}
}

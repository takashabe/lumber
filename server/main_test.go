package server

import (
	"os"
	"testing"

	"github.com/takashabe/go-fixture"
	_ "github.com/takashabe/go-fixture/mysql" // driver
	"github.com/takashabe/lumber/datastore"
)

func TestMain(m *testing.M) {
	setupTables()
	os.Exit(m.Run())
}

func setupTables() {
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

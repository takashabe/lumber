package persistence

import (
	"os"
	"testing"

	fixture "github.com/takashabe/go-fixture"
	_ "github.com/takashabe/go-fixture/mysql"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

func setup() {
	db, err := NewDatastore()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	f := fixture.NewFixture(db.Conn, "mysql")
	err = f.LoadSQL("testdata/schema.sql")
	if err != nil {
		panic(err)
	}
}

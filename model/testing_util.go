package model

import (
	"testing"

	"github.com/takashabe/go-fixture"
	_ "github.com/takashabe/go-fixture/mysql" // driver
	"github.com/takashabe/lumber/datastore"
)

func loadFixture(t *testing.T, file string) {
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

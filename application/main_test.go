package application

import (
	"os"
	"testing"

	"github.com/takashabe/lumber/domain/repository"
	"github.com/takashabe/lumber/helper"
	"github.com/takashabe/lumber/infrastructure/persistence"
)

func TestMain(m *testing.M) {
	helper.SetupTables()
	os.Exit(m.Run())
}

func getEntryRepository(t *testing.T) repository.EntryRepository {
	r, err := persistence.NewEntryRepository()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	return r
}

func getTokenRepository(t *testing.T) repository.TokenRepository {
	r, err := persistence.NewTokenRepository()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	return r
}

package interfaces

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/takashabe/lumber/helper"
	"github.com/takashabe/lumber/infrastructure/persistence"
)

func TestMain(m *testing.M) {
	helper.SetupTables()
	os.Exit(m.Run())
}

func setupServer(t *testing.T) *httptest.Server {
	entryRepo, err := persistence.NewEntryRepository()
	tokenRepo, err := persistence.NewTokenRepository()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	server := &Server{
		Entry: NewEntryHandler(entryRepo, tokenRepo),
		Token: NewTokenHandler(tokenRepo),
	}
	return httptest.NewServer(server.Routes())
}

func sendRequest(t *testing.T, method, url string, body io.Reader) *http.Response {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	return res
}

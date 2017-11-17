package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/takashabe/go-fixture"
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

func setupServer(t *testing.T) *httptest.Server {
	s, err := NewServer()
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	return httptest.NewServer(s.Routes())
}

func TestGetEntry(t *testing.T) {
	ts := setupServer(t)
	defer ts.Close()

	cases := []struct {
		input  int
		expect int
	}{
		{1, http.StatusOK},
		{0, http.StatusNotFound},
	}
	for i, c := range cases {
		loadFixture(t, "fixture/entries.yml")
		res, err := http.Get(ts.URL + fmt.Sprintf("/api/entry/%d", c.input))
		if err != nil {
			t.Fatalf("#%d: want non error, got %v", i, err)
		}
		defer res.Body.Close()

		if res.StatusCode != c.expect {
			t.Errorf("#%d: want %d, got %d", i, c.expect, res.StatusCode)
		}
	}
}

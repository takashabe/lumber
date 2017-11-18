package server

import (
	"fmt"
	"io"
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

func sendRequest(t *testing.T, method, url string, body io.Reader) *http.Response {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	return res
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
		res := sendRequest(t, "GET", fmt.Sprintf("%s/api/entry/%d", ts.URL, c.input), nil)
		defer res.Body.Close()

		if res.StatusCode != c.expect {
			t.Errorf("#%d: want %d, got %d", i, c.expect, res.StatusCode)
		}
	}
}
		if err != nil {
			t.Fatalf("#%d: want non error, got %v", i, err)
		}
		defer res.Body.Close()

		if res.StatusCode != c.expect {
			t.Errorf("#%d: want %d, got %d", i, c.expect, res.StatusCode)
		}
	}
}

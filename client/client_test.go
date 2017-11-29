package client

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/takashabe/lumber/server"
)

func setupServer(t *testing.T) *httptest.Server {
	s, err := server.NewServer()
	if err != nil {
		t.Fatalf("want non error, got %v", err)
	}
	ts := httptest.NewServer(s.Routes())
	os.Setenv(LumberServerAddress, ts.URL)
	return ts
}

func TestCreateAndGetEntry(t *testing.T) {
	ts := setupServer(t)
	defer ts.Close()

	cases := []struct {
		input         string
		expectTitle   string
		expectContent string
	}{
		{
			"testdata/minimum.md",
			"title",
			"<p>content</p>",
		},
	}
	for i, c := range cases {
		ctx := context.Background()
		client, err := New()
		if err != nil {
			t.Fatalf("#%d: want non error, got %v", i, err)
		}
		id, err := client.CreateEntry(ctx, c.input)
		if err != nil {
			t.Fatalf("#%d: want non error, got %v", i, err)
		}
		entry, err := client.Entry(id).Get(ctx)
		if err != nil {
			t.Fatalf("#%d: want non error, got %v", i, err)
		}
		if entry.ID != id {
			t.Errorf("#%d: want id %d, got %d", i, id, entry.ID)
		}
		if entry.Title != c.expectTitle || entry.Content != c.expectContent {
			t.Errorf("#%d: want title %s and content %s, got %s and %s",
				i, c.expectTitle, c.expectContent, entry.Title, entry.Content)
		}
	}
}

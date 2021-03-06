package client

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/takashabe/lumber/helper"
	"github.com/takashabe/lumber/infrastructure/persistence"
	"github.com/takashabe/lumber/interfaces"
)

func setupServer(t *testing.T) *httptest.Server {
	er, err := persistence.NewEntryRepository()
	tr, err := persistence.NewTokenRepository()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	server := &interfaces.Server{
		Entry: interfaces.NewEntryHandler(er, tr),
		Token: interfaces.NewTokenHandler(tr),
	}
	ts := httptest.NewServer(server.Routes())
	os.Setenv(LumberServerAddress, ts.URL)
	return ts
}

func TestCreateAndGetEntry(t *testing.T) {
	ts := setupServer(t)
	defer ts.Close()
	helper.InitializeTable()

	cases := []struct {
		input         string
		expectTitle   string
		expectContent string
	}{
		{
			"testdata/minimum.md",
			"min_title",
			"<p>content</p>",
		},
	}
	for i, c := range cases {
		ctx := context.Background()
		client, err := New()
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		id, err := client.CreateEntry(ctx, c.input)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		entry, err := client.Entry(id).Get(ctx)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
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

func TestCreateDuplicateEntry(t *testing.T) {
	ts := setupServer(t)
	defer ts.Close()
	helper.InitializeTable()

	cases := []struct {
		input string
	}{
		{"testdata/minimum.md"},
		{"testdata/minimum.md"},
	}
	for i, c := range cases {
		ctx := context.Background()
		client, err := New()
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		_, err = client.CreateEntry(ctx, c.input)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
	}
}

func TestEditEntry(t *testing.T) {
	ts := setupServer(t)
	defer ts.Close()

	cases := []struct {
		inputID       int
		inputFile     string
		expectTitle   string
		expectContent string
	}{
		{
			1,
			"testdata/minimum.md",
			"min_title",
			"<p>content</p>",
		},
		{
			1,
			"testdata/minimum2.md",
			"min_title_2",
			"<h2>content</h2>\n\n<ul>\n<li>list</li>\n<li>list2</li>\n</ul>",
		},
	}
	for i, c := range cases {
		helper.LoadFixture(t, "fixture/entries.yml")
		ctx := context.Background()
		client, err := New()
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		entryClient := client.Entry(c.inputID)
		err = entryClient.Edit(ctx, c.inputFile)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		entry, err := entryClient.Get(ctx)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		if entry.Title != c.expectTitle || entry.Content != c.expectContent {
			t.Errorf("#%d: want title %s and content %s, got %s and %s",
				i, c.expectTitle, c.expectContent, entry.Title, entry.Content)
		}
	}
}

func TestDeleteEntry(t *testing.T) {
	ts := setupServer(t)
	defer ts.Close()

	cases := []struct {
		input int
	}{
		{
			1,
		},
	}
	for i, c := range cases {
		helper.LoadFixture(t, "fixture/entries.yml")
		ctx := context.Background()
		client, err := New()
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
		err = client.Entry(c.input).Delete(ctx)
		if err != nil {
			t.Fatalf("#%d: want non error, got %#v", i, err)
		}
	}
}

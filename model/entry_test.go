package model

import (
	"database/sql"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/takashabe/lumber/config"
)

func TestNewEntry(t *testing.T) {
	cases := []struct {
		inputFilePath     string
		expectTitle       string
		expectContentHead string
		expectContentTail string
		expectErr         error
	}{
		{
			"testdata/go-pubsub_readme.md",
			"go-pubsub",
			"<p><a href=",
			"</ul>",
			nil,
		},
		{
			"testdata/minimum.md",
			"title",
			"<p>content</p>",
			"<p>content</p>",
			nil,
		},
		{
			"testdata/empty.md",
			"",
			"",
			"",
			config.ErrEmptyEntry,
		},
	}
	for i, c := range cases {
		data, err := ioutil.ReadFile(c.inputFilePath)
		if err != nil {
			t.Fatalf("#%d: want non error, got %v", i, err)
		}
		entry, err := NewEntry(data, EntryStatusPublic)
		if err != c.expectErr {
			t.Fatalf("#%d: want error %v, got %v", i, c.expectErr, err)
		}
		if err != nil {
			continue
		}

		if entry.Title != c.expectTitle {
			t.Fatalf("#%d: want title %s, got %s", i, c.expectTitle, entry.Title)
		}
		if !strings.HasPrefix(entry.Content, c.expectContentHead) {
			t.Fatalf("#%d: want content prefix %s", i, c.expectContentHead)
		}
		if !strings.HasSuffix(entry.Content, c.expectContentTail) {
			t.Fatalf("#%d: want content suffix %s", i, c.expectContentTail)
		}
	}
}

func TestGetEntry(t *testing.T) {
	cases := []struct {
		input     int
		expectID  int
		expectErr error
	}{
		{1, 1, nil},
		{0, 0, sql.ErrNoRows},
	}
	for i, c := range cases {
		loadFixture(t, "fixture/entries.yml")

		act, err := GetEntry(c.input)
		if err != c.expectErr {
			t.Fatalf("#%d: want error %v, got %v", i, c.expectErr, err)
		}
		if err != nil {
			continue
		}

		if act.ID != c.expectID {
			t.Fatalf("#%d: want %d, got %d", i, c.expectID, act.ID)
		}
	}
}

func TestPost(t *testing.T) {
	cases := []struct {
		inputFilePath string
	}{
		{"testdata/go-pubsub_readme.md"},
		{"testdata/minimum.md"},
	}
	for i, c := range cases {
		data, err := ioutil.ReadFile(c.inputFilePath)
		if err != nil {
			t.Fatalf("#%d: want non error, got %v", i, err)
		}
		entry, err := NewEntry(data, EntryStatusPublic)
		if err != nil {
			t.Fatalf("#%d: want non error, got %v", i, err)
		}
		err = entry.Post()
		if err != nil {
			t.Fatalf("#%d: want non error, got %v", i, err)
		}
	}
}
